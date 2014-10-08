package main

import (
	"fmt"
	"log"
	"math"
	"time"
)

type Calibration struct {
	AccZeros   PRY
	AccScale   PRY
	GyroZeros  PRY
	GyroScale  PRY
	SonarScale float64
	SonarZero  float64
}

func (c Calibration) Convert(data Navdata) Sensors {
	return Sensors{
		Acc: PRY{
			Pitch: (float64(data.AccPitch) - c.AccZeros.Pitch) / c.AccScale.Pitch,
			Roll:  (float64(data.AccRoll) - c.AccZeros.Roll) / c.AccScale.Roll,
			Yaw:   (float64(data.AccYaw) - c.AccZeros.Yaw) / c.AccScale.Yaw,
		},
		Gyro: PRY{
			Pitch: (float64(data.GyroPitch) - c.GyroZeros.Pitch) / c.GyroScale.Pitch,
			Roll:  (float64(data.GyroRoll) - c.GyroZeros.Roll) / c.GyroScale.Roll,
			Yaw:   (float64(data.GyroYaw) - c.GyroZeros.Yaw) / c.GyroScale.Yaw,
		},
		Sonar: (float64(data.Ultrasound&0x7FFF) - c.SonarZero) / c.SonarScale,
	}
}

type Calibrator struct {
	Samples    int
	MaxStdDev  float64
	Navboard   *Navboard
	Motorboard *Motorboard
}

func (c *Calibrator) Calibrate(r *Calibration) error {
	log.Printf("Calibrating")
	var (
		err     error
		samples [][6]float64
		sums    [6]float64
		means   [6]float64
		sqrSums [6]float64
		stdDevs [6]float64
	)
	for i := 0; i < c.Samples; i++ {
		var data Navdata
		data, err = c.Navboard.Read()
		if err != nil {
			break
		}
		sample := [6]float64{
			float64(data.AccPitch),
			float64(data.AccRoll),
			float64(data.AccYaw),
			float64(data.GyroPitch),
			float64(data.GyroRoll),
			float64(data.GyroYaw),
		}
		samples = append(samples, sample)
		for i, val := range sample {
			sums[i] += val
		}
		// Make the LEDs blink green while reading calibration samples
		if i%20 > 10 {
			c.Motorboard.WriteLeds(Leds(LedGreen))
		} else {
			c.Motorboard.WriteLeds(Leds(LedOff))
		}
		_ = data
	}
	if err == nil {
		for i, sum := range sums {
			means[i] = sum / float64(len(samples))
		}
		for _, sample := range samples {
			for i, val := range sample {
				sqrSums[i] += (val - means[i]) * (val - means[i])
			}
		}
		for i, sqrSum := range sqrSums {
			stdDevs[i] = math.Sqrt(sqrSum / float64(len(samples)))
		}
	}
	for i, stdDev := range stdDevs {
		// Detect if there is too much sensor noise (due to the drone moving, or
		// sensor issues).
		if stdDev > c.MaxStdDev {
			err = fmt.Errorf("Standard deviation too high")
		} else {
			var val = means[i]
			switch i {
			case 0:
				r.AccZeros.Pitch = val
			case 1:
				r.AccZeros.Roll = val
			case 2:
				// Yaw should measure 1G during calibration, so we'll assume its zero
				// value to be of similar value as pitch/yaw
				r.AccZeros.Yaw = (r.AccZeros.Pitch + r.AccZeros.Roll) / 2
				// Now we can estimate the scale of 1G for all sensors.
				r.AccScale.Pitch = val - r.AccZeros.Yaw
				r.AccScale.Roll = val - r.AccZeros.Yaw
				r.AccScale.Yaw = val - r.AccZeros.Yaw
			case 3:
				r.GyroZeros.Pitch = val
			case 4:
				r.GyroZeros.Roll = val
			case 5:
				r.GyroZeros.Yaw = val
			}
		}
	}
	if err != nil {
		log.Printf("Failed to calibrate: %s", err)
		c.Motorboard.WriteLeds(Leds(LedRed))
		time.Sleep(time.Second)
	} else {
		log.Printf("Calibration succeeded")
	}
	c.Motorboard.WriteLeds(Leds(LedGreen))
	return nil
}
