package main

import (
	"fmt"
	"github.com/felixge/godrone"
	"github.com/felixge/godrone/attitude"
	"github.com/felixge/godrone/drivers/motorboard"
	"github.com/felixge/pidctrl"
	"os"
	"os/signal"
)

var (
	navboard = godrone.DefaultNavboard
	log      = godrone.DefaultLogger
	speed    = 0.01
)

func main() {
	m, err := motorboard.NewMotorboard(motorboard.DefaultTTY, log)
	if err != nil {
		panic(err)
	}

	if err := navboard.Calibrate(); err != nil {
		log.Fatal(err)
	}
	defer navboard.Close()

	go loop(m)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	<-sigCh

	log.Debug("Shutting down")
}

func loop(m *motorboard.Motorboard) {
	var (
		p         = 0.9
		i         = 0.4
		d         = 0.2
		basespeed = 0.5
		roll      = pidctrl.NewPIDController(p, i, d)
		pitch     = pidctrl.NewPIDController(p, i, d)
		yaw       = pidctrl.NewPIDController(p, i, d)
		filter    = attitude.NewComplementary()
	)

	for {
		data, err := navboard.NextData()
		if err != nil {
			continue
		}

		a := filter.Update(data.Data)
		ro := roll.Update(a.Roll / 90)
		po := pitch.Update(a.Pitch / 90)
		yo := yaw.Update(a.Yaw / 90) / 5

		fmt.Printf("roll: %.2f pitch: %.2f yaw: %.2f - %s\r", ro, po, a)

		speeds := [4]float64{basespeed, basespeed, basespeed, basespeed}

		if ro > 0 {
			speeds[0], speeds[3] = speeds[0]+ro, speeds[3]+ro
		} else if ro < 0 {
			speeds[1], speeds[2] = speeds[1]-ro, speeds[2]-ro
		}

		if po > 0 {
			speeds[0], speeds[1] = speeds[0]+po, speeds[1]+po
		} else if po < 0 {
			speeds[2], speeds[3] = speeds[2]-po, speeds[3]-po
		}

		if yo > 0 {
			speeds[1], speeds[3] = speeds[1]+yo, speeds[3]+yo
		} else if yo < 0 {
			speeds[0], speeds[2] = speeds[0]-yo, speeds[2]-yo
		}

		for i, speed := range speeds {
			m.SetSpeed(i, speed)
		}
	}
}
