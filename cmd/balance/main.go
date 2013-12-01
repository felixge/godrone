package main

import (
	"fmt"
	"github.com/felixge/godrone"
	"github.com/felixge/godrone/attitude"
	"github.com/felixge/godrone/drivers/motorboard"
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
		c          = attitude.NewComplementary()
		baseSpeeds = [4]float64{0.005, 0.005, 0.005, 0.005}
	)

	for {
		data, err := navboard.NextData()
		if err != nil {
			continue
		}
		var (
			a      = c.Update(data.Data)
			speeds = baseSpeeds
			roll   = a.Roll / 90
		)

		fmt.Printf("\r%s", a)

		roll *= 2
		if roll > 0 {
			speeds[1], speeds[2] = speeds[1]+roll, speeds[2]+roll
		} else if roll < 0 {
			speeds[0], speeds[3] = speeds[0]-roll, speeds[3]-roll
		}

		for i, speed := range speeds {
			m.SetSpeed(i, speed)
		}
	}
}
