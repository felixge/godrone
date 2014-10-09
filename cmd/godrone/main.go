package main

import (
	"fmt"
	"log"
	"time"

	"github.com/felixge/godrone"
)

type State int

const (
	None State = iota
	Landed
	LandStart
	Land
	Calibrate
	TakeoffStart
	Takeoff
	Fly
)

const (
	// TakeoffAltitude is the altitude to aim for on takeoff
	TakeoffAltitude = 0.5
	// LandAltitude is the altitude at which to cutoff the engines when landing.
	// Must be > 0 due to the sonar being unable to measure small distances.
	LandAltitude = 0.3
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
	log.Printf("Godrone started")
	firmware, err := godrone.NewFirmware()
	if err != nil {
		log.Fatalf("%s", err)
	}
	defer firmware.Close()
	go serveHttp()
	firmware.Motorboard.WriteLeds(godrone.Leds(godrone.LedGreen))
	var state = Calibrate
	for {
		var err error
		switch state {
		case Calibrate:
			// @TODO The LEDs don't seem to turn off when this is called again after
			// a calibration errors, instead they just blink. Not sure why.
			firmware.Motorboard.WriteLeds(godrone.Leds(godrone.LedOff))
			log.Printf("Calibrating sensors")
			err = firmware.Calibrate()
			if err != nil {
				firmware.Motorboard.WriteLeds(godrone.Leds(godrone.LedRed))
				time.Sleep(time.Second)
			} else {
				log.Printf("Finished calibration")
				state = Landed
				firmware.Motorboard.WriteLeds(godrone.Leds(godrone.LedGreen))
			}
		case Landed:
			err = firmware.Observe()
		case TakeoffStart:
			firmware.Desired.Altitude = TakeoffAltitude
			firmware.Desired.PRY = godrone.PRY{}
			state = Takeoff
		case Takeoff:
			err = firmware.Fly()
			if firmware.Actual.Altitude >= firmware.Desired.Altitude {
				state = Fly
			}
		case Fly:
			err = firmware.Fly()
		case Land:
			firmware.Desired.Altitude = 0
			firmware.Desired.PRY = godrone.PRY{}
			err = firmware.Fly()
			if firmware.Actual.Altitude <= LandAltitude {
				state = Landed
			}
		default:
			panic(fmt.Errorf("Unhandled state: %s", state))
		}
		if err != nil {
			log.Printf("%s", err)
		}
	}
}

func serveHttp() {
}
