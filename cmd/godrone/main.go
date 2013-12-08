package main

import (
	"github.com/felixge/godrone/attitude"
	"github.com/felixge/godrone/control"
	"github.com/felixge/godrone/drivers/motorboard"
	"github.com/felixge/godrone/drivers/navboard"
	"github.com/felixge/godrone/http"
	"github.com/felixge/log"
	gohttp "net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Config struct {
	NavboardTTY   string
	MotorboardTTY string
	RollPID       [3]float64
	PitchPID      [3]float64
	YawPID        [3]float64
	HttpAddr      string
}

var (
	defaultRollPitchPID = [3]float64{0.02, 0, 0}

	green  = motorboard.Leds(motorboard.LedGreen)
	orange = motorboard.Leds(motorboard.LedOrange)
	red    = motorboard.Leds(motorboard.LedRed)
)

var DefaultConfig = Config{
	NavboardTTY:   "/dev/ttyO1",
	MotorboardTTY: "/dev/ttyO0",
	RollPID:       defaultRollPitchPID,
	PitchPID:      defaultRollPitchPID,
	YawPID:        [3]float64{1, 0, 0},
	HttpAddr:      ":80",
}

type Instances struct {
	log        *log.Logger
	navboard   *navboard.Navboard
	motorboard *motorboard.Motorboard
	attitude   *attitude.Complementary
	control    *control.Control
	http       gohttp.Handler
}

func main() {
	config := DefaultConfig
	i, err := NewInstances(config)
	if err != nil {
		panic(err)
	}
	i.log.Info("Starting godrone")
	defer i.motorboard.Close()

	i.motorboard.SetLeds(green)
	time.Sleep(500 * time.Millisecond)
	i.motorboard.SetLeds(red)

	i.log.Info("Calibrating sensors")
	for {
		if err := i.navboard.Calibrate(); err == nil {
			break
		}
	}
	i.motorboard.SetLeds(green)

	navDataCh := make(chan navboard.Data)
	go readNavData(i.navboard, navDataCh)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT)

	go gohttp.ListenAndServe(config.HttpAddr, i.http)

	i.log.Info("Entering main loop")
mainLoop:
	for {
		select {
		case navData := <-navDataCh:
			attitudeData := i.attitude.Update(navData.Data)
			motorSpeeds := i.control.Update(attitudeData)
			if err := i.motorboard.SetSpeeds(motorSpeeds); err != nil {
				i.log.Error("Could not set motor speeds. err=%s", err)
			}
		case sig := <-sigCh:
			i.log.Info("Received signal=%s, shutting down", sig)
			break mainLoop
		}
	}
}

func readNavData(board *navboard.Navboard, ch chan<- navboard.Data) {
	for {
		navData, err := board.NextData()
		if err != nil {
			continue
		}
		select {
		case ch <- navData:
		default:
		}
	}
}

func NewInstances(c Config) (i Instances, err error) {
	i.log = log.DefaultLogger
	i.navboard = navboard.NewNavboard(c.NavboardTTY, i.log)
	i.motorboard, err = motorboard.NewMotorboard(c.MotorboardTTY)
	if err != nil {
		return
	}
	i.attitude = attitude.NewComplementary()
	i.control = control.NewControl(c.RollPID, c.PitchPID, c.YawPID)
	i.http = http.NewHandler(i.control, i.log)
	return
}
