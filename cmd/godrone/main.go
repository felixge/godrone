package main

import (
	"log"
	"time"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
	log.Printf("Godrone started")
	navboard, err := OpenNavboard("/dev/ttyO1")
	if err != nil {
		log.Fatalf("Failed to open navboard: %s", err)
	}
	defer navboard.Close()
	motorboard, err := OpenMotorboard("/dev/ttyO0")
	if err != nil {
		log.Fatalf("Failed to open navboard: %s", err)
	}
	defer motorboard.Close()
	var (
		// found by experimenting
		// overwritten by calibration (except for gyro/sonar scale)
		calibration = Calibration{
			AccZeros:   PRY{Pitch: 2038, Roll: 2070, Yaw: 2048},
			AccScale:   PRY{Pitch: 1, Roll: 1, Yaw: 1},
			GyroZeros:  PRY{Pitch: 16, Roll: -27, Yaw: 1.5},
			GyroScale:  PRY{Pitch: -16, Roll: 16, Yaw: 16},
			SonarScale: 3500,
		}
		calibrator = &Calibrator{
			Motorboard: motorboard,
			Navboard:   navboard,
			Samples:    200,
			MaxStdDev:  5,
		}
		// We know the navboard runs at 200 Hz. Assuming that the navboards
		// internal clock is more accurate than the system clock, we should get
		// better precision from hardcoding this assumption
		dt = time.Second / 200
		// also found by experimenting :)
		filter = Filter{
			GyroGain:  0.98,
			AccGain:   0.02,
			SonarGain: 0.1,
			SonarMax:  4,
		}
		controller = &Controller{
			Pitch:        PID{P: 0.02, I: 0.0001, D: 0},
			Roll:         PID{P: 0.02, I: 0.0001, D: 0},
			Yaw:          PID{P: 0.02, I: 0, D: 0},
			Altitude:     PID{P: 0.2, I: 0.05, D: 0.05},
			LandAltitude: 0.3,
		}
		actual      State
		desired     State
		calibrateCh = make(chan struct{}, 1)
		stateCh     = make(chan State, 1000)
	)
	calibrateCh <- struct{}{}
	motorboard.WriteLeds(Leds(LedGreen))
	go mission(stateCh)
	for {
		motorboard.WriteLeds(Leds(LedGreen))
	updateLoop:
		for {
			select {
			case <-calibrateCh:
				if actual.Fly {
					log.Printf("Refusing to calibrate while flying")
					break
				}
				for {
					if err := calibrator.Calibrate(&calibration); err == nil {
						break
					}
				}
			case desired = <-stateCh:
			default:
				break updateLoop
			}
		}
		navdata, err := navboard.Read()
		if err != nil {
			log.Printf("Failed to read navdata: %s", err)
			continue
		}
		sensors := calibration.Convert(navdata)
		filter.Update(&actual, sensors, dt)
		if err := controller.Apply(desired, &actual, dt, motorboard); err != nil {
			log.Printf("Failed to apply control: %s", err)
			continue
		}
	}
}

func mission(stateCh chan State) {
	var state State
	state.Fly = true
	state.Altitude = 0.6
	state.Orientation.Pitch = -1
	state.Orientation.Roll = -0.5
	stateCh <- state
}
