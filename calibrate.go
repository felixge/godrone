package godrone

import (
	"fmt"
	"math"
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
	Samples   int
	MaxStdDev float64
}

func (c *Calibrator) Calibrate(navboard NavdataReader, r *Calibration) error {
	var (
		samples [][6]float64
		sums    [6]float64
		means   [6]float64
		sqrSums [6]float64
		stdDevs [6]float64
	)
	for i := 0; i < c.Samples; i++ {
		data, err := navboard.Read()
		if err != nil {
			return err
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
	}
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
	// Detect if there is too much sensor noise (due to the drone moving, or
	// sensor issues).
	for _, stdDev := range stdDevs {
		if stdDev > c.MaxStdDev {
			return fmt.Errorf("Standard deviation too high")
		}
	}
	for i, mean := range means {
		switch i {
		case 0:
			r.AccZeros.Pitch = mean
		case 1:
			r.AccZeros.Roll = mean
		case 2:
			if mean < r.AccZeros.Pitch || mean < r.AccZeros.Roll {
				return fmt.Errorf("Bad yaw, is drone on its back?")
			}
			// Yaw should measure 1G during calibration, so we'll assume its zero
			// value to be between pitch and roll.
			r.AccZeros.Yaw = (r.AccZeros.Pitch + r.AccZeros.Roll) / 2
			// Now we can estimate the scale of 1G for all sensors.
			r.AccScale.Pitch = mean - r.AccZeros.Yaw
			r.AccScale.Roll = mean - r.AccZeros.Yaw
			r.AccScale.Yaw = mean - r.AccZeros.Yaw
		case 3:
			r.GyroZeros.Pitch = mean
		case 4:
			r.GyroZeros.Roll = mean
		case 5:
			r.GyroZeros.Yaw = mean
		}
	}
	return nil
}
