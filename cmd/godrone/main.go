package main

import (
	"flag"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/felixge/godrone"
	"github.com/gorilla/websocket"
)

var verbose = flag.Int("verbose", 0, "verbosity: 1=some 2=lots")
var addr = flag.String("addr", ":80", "Address to listen on (default is \":80\")")
var dummy = flag.Bool("dummy", false, "Dummy drone: do not open navboard/motorboard.")

const pitchLimit = 30 // degrees
const rollLimit = 30  // degrees

func main() {
	cutoutReason := "not calibrated"
	flag.Parse()

	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
	log.Printf("Godrone started")

	var firmware *godrone.Firmware
	if *dummy {
		firmware, _ = godrone.NewCustomFirmware(&mockNavboard{}, &mockMotorboard{})
	} else {
		var err error
		firmware, err = godrone.NewFirmware()
		if err != nil {
			log.Fatalf("%s", err)
		}
		defer firmware.Close()
	}

	// This is the channel to send requests to the main control loop.
	// It is the only way to send and receive info the the
	// godrone.Firmware object (which is non-concurrent).
	reqCh := make(chan Request)

	// Autonomy/guard goroutines.
	go monitorAngles(reqCh)
	go lander(reqCh)

	// The websocket input goroutine.
	go serveHttp(reqCh)

	calibrate := func() {
		for {
			// @TODO The LEDs don't seem to turn off when this is called again after
			// a calibration errors, instead they just blink. Not sure why.
			firmware.Motorboard.WriteLeds(godrone.Leds(godrone.LedOff))
			log.Printf("Calibrating sensors")
			err := firmware.Calibrate()
			if err != nil {
				firmware.Motorboard.WriteLeds(godrone.Leds(godrone.LedRed))
				time.Sleep(time.Second)
			} else {
				log.Printf("Finished calibration")
				firmware.Motorboard.WriteLeds(godrone.Leds(godrone.LedGreen))
				cutoutReason = ""
				break
			}
		}
	}
	calibrate()
	log.Print("Up, up and away!")

	// This is the main control loop.
	flying := false
	for {
		select {
		case req := <-reqCh:
			var res Response
			if req.Include != nil {
				if req.Include.Calibration {
					res.Calibration = &firmware.Calibration
				}
			}
			res.Cutout = cutoutReason
			res.Actual = firmware.Actual
			res.Desired = firmware.Desired
			res.Time = time.Now()

			if req.SetDesired != nil {
				// Log changes to desired
				if Verbose() && firmware.Desired != *req.SetDesired {
					log.Print("New desired attitude:", firmware.Desired)
				}
				firmware.Desired = *req.SetDesired
			}
			if req.Cutout != "" {
				cutoutReason = req.Cutout
				log.Print("Cutout: ", cutoutReason)
				firmware.Desired.Altitude = 0
			}
			if req.Calibrate {
				if req.CustomCalibrationData != nil {
					firmware.Calibration = *req.CustomCalibrationData
				} else {
					calibrate()
				}
			}
			req.response <- res
			if reallyVerbose() {
				log.Print("Request: ", req, "Response: ", res)
			}
		default:
		}

		var err error
		if firmware.Desired.Altitude > 0 {
			err = firmware.Control()
			flying = true
		} else {
			// Something subtle to note here: When the motors are running,
			// but then desired altitude goes to zero (for instance due to
			// the emergency command in the UI) we end up here.
			//
			// We never actually send a motors=0 command. Instead we count on
			// the failsafe behavior of the motorboard, which puts the motors to
			// zero if it does not receive a new motor command soon enough.
			err = firmware.Observe()
			if flying {
				log.Print("Motor cutoff.")
				flying = false
			}
		}
		if err != nil {
			log.Printf("%s", err)
		}
	}
}

type IncludeData struct {
	Calibration bool
}

type Request struct {
	SetDesired            *godrone.Placement
	Include               *IncludeData
	Calibrate             bool
	CustomCalibrationData *godrone.Calibration
	Land                  bool
	Cutout                string
	response              chan Response
}

func newRequest() Request { return Request{response: make(chan Response)} }

type Response struct {
	Time        time.Time            `json:"time"`
	Actual      godrone.Placement    `json:"actual,omitempty"`
	Desired     godrone.Placement    `json:"desired,omitempty"`
	Calibration *godrone.Calibration `json:"calibration,omitempty"`
	Cutout      string               `json:"cutout,omitempty"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func serveHttp(reqCh chan<- Request) {
	err := http.ListenAndServe(*addr, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		first := true
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("Failed to upgrade ws: %s", err)
			return
		}
		defer conn.Close()
		for {
			req := newRequest()
			if err := conn.ReadJSON(&req); err != nil {
				log.Printf("Failed to read from ws: %s", err)
				return
			}

			// On the first request from the websocket, we ignore what they
			// asked for. This way, they find out our current altitude, and
			// can sync up to us, thereby not causing a crash on reconnect.
			if first {
				req.SetDesired = nil
				first = false
			}

			if req.Land {
				log.Print("Landing requested.")
				landCh <- landStart
			}

			// Cancel landing when there is flight input.
			if req.SetDesired != nil {
				landCh <- landCancel
			}

			// Send the request into the control loop.
			reqCh <- req
			res := <-req.response

			if err := conn.WriteJSON(res); err != nil {
				log.Printf("Failed to write to ws: %s", err)
				return
			}
		}
	}))
	if err != nil {
		log.Fatalf("Failed to serve http: %s", err)
	}
}

func reallyVerbose() bool { return *verbose > 1 }
func Verbose() bool       { return *verbose > 0 }

// A MockNavboard is a NavBoardReader that does not talk to real hardware.
type mockNavboard struct {
	seq uint16
}

func (b *mockNavboard) Read() (data godrone.Navdata, err error) {
	b.seq++
	return godrone.Navdata{
		Seq:      b.seq,
		AccRoll:  2000,
		AccPitch: 2000,
		AccYaw:   8000,
	}, nil
}

// A MockMotorboard is a MotorLedWriter that does not talk to real hardware.
type mockMotorboard struct {
	speeds [4]float64
}

func (m *mockMotorboard) WriteLeds(leds [4]godrone.LedColor) error {
	if reallyVerbose() {
		log.Print("Mock WriteLeds: ", leds)
	}
	return nil
}

func (m *mockMotorboard) WriteSpeeds(speeds [4]float64) error {
	if reallyVerbose() && m.speeds != speeds {
		log.Print("Mock WriteSpeeds: ", speeds)
	}
	m.speeds = speeds
	return nil
}

// This runs in its own goroutine. It cuts the motors if pitch or roll are over the limit.
func monitorAngles(reqCh chan<- Request) {
	for {
		// Ask for the current status (by sending an empty request).
		req := newRequest()
		reqCh <- req
		resp := <-req.response

		// If we are still flying...
		if resp.Cutout == "" {
			if math.Abs(resp.Actual.PRY.Pitch) > pitchLimit {
				req.Cutout = "Pitch angle is too high"
			}
			if math.Abs(resp.Actual.PRY.Roll) > rollLimit {
				req.Cutout = "Roll angle is too high"
			}

			// if we decided to cutout, send the command in
			if req.Cutout != "" {
				reqCh <- req
				// ignore the response (but we have to pick it up
				// to unblock the control loop)
				<-req.response
			}
		}

		time.Sleep(100 * time.Millisecond)
	}
}

// A landCmd is a command sent into the lander goroutine via the landCh.
type landCmd int

const (
	landStart landCmd = iota
	landCancel
)

// Buffered because nothing good would come from blocking on sending
// and maybe some bad.
var landCh = make(chan landCmd, 10)

// Descent reate in meters/sec.
const descentRate = 0.5 // m/s
// Descend until here, then cut motors.
const descendUntil = 0.3 // m

// How often landing adjustments are issued (number per second).
const landingHz = 5

var landingSleep = time.Second / landingHz

// A goroutine that commands a controlled landing, once requested by a send on landStart.
// It can be canceled by sending a signal to landCancel.
func lander(reqCh chan<- Request) {
	for {
		// waiting for start
		cmd := landCancel
		for cmd != landStart {
			cmd = <-landCh
		}

		// Bring desired altitude down to .3 at .5 m/s. At .3 cut motors,
		// to drop to the ground.
		for cmd == landStart {
			// Ask for the current status (by sending an empty request).
			req := newRequest()
			reqCh <- req
			resp := <-req.response

			// Descend a bit
			newAlt := resp.Desired.Altitude
			if newAlt < 2*descendUntil {
				// slow down near the cutoff
				newAlt -= descentRate / 3 / landingHz
			} else {
				newAlt -= descentRate / landingHz
			}

			// Cutoff motors?
			if newAlt < descendUntil {
				newAlt = 0
				cmd = landCancel
			}

			// Update the req and send it in to the control loop.
			req.SetDesired = &godrone.Placement{}
			*req.SetDesired = resp.Desired
			req.SetDesired.Altitude = newAlt
			reqCh <- req
			resp = <-req.response
			time.Sleep(landingSleep)

			// check for cancel (or a dup start)
			select {
			case cmd = <-landCh:
				// Got a cmd.
				// It is either a duplicate start, so we keep going
				// or a cancel, so the loop ends, and we got to the top.
			default:
				// Do not block on chan read.
			}
		}
	}
}
