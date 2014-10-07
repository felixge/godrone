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
		// @TODO allow user to trim these values when drone is flat on ground
		calibration = Calibration{
			AccZeros:        PRY{Pitch: 2032, Roll: 2078, Yaw: 2020},
			AccSensitivity:  PRY{Pitch: 1, Roll: 1, Yaw: 1},
			GyroZeros:       PRY{Pitch: 7.8, Roll: -19, Yaw: 10.3},
			GyroSensitivity: PRY{Pitch: -16, Roll: 16, Yaw: 16},
		}
		// seems to work well in practice
		filter = Complementary{GyroGain: 0.98, AccGain: 0.02}
		// We know the navboard runs at 200 Hz. Assuming that the navboards
		// internal clock is more accurate than the system clock, we should get
		// better precision from hardcoding this assumption
		dt = time.Second / 200
		// Orientation and altitude of the drone as determined by sensor filtering.
		attitude Attitude
	)
	// blink the LEDs to signal that godrone has started
	motorboard.WriteLeds(Leds(LedOff))
	time.Sleep(200 * time.Millisecond)
	for {
		navdata, err := navboard.Read()
		if err != nil {
			log.Printf("Failed to read navdata: %s", err)
			continue
		}
		sensors := calibration.Convert(navdata)
		filter.Update(&attitude.PRY, sensors.Acc, sensors.Gyro, dt)
		log.Printf("%s", attitude)
		motorboard.WriteSpeeds([4]float64{0.1, 0.1, 0.1, 0.1})
		motorboard.WriteLeds(Leds(LedGreen))
	}
}
