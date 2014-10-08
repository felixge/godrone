package main

type Sensors struct {
	// Acceleration in m/s^2
	Acc PRY
	// Rotation in deg/s
	Gyro PRY
	// Altitude in m
	Sonar float64
}

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
			Pitch: (float64(data.GyroPitch) - c.Gyroeros.Pitch) / c.GyroScale.Pitch,
			Roll:  (float64(data.GyroRoll) - c.GyroZeros.Roll) / c.GyroScale.Roll,
			Yaw:   (float64(data.GyroYaw) - c.GyroZeros.Yaw) / c.GyroScale.Yaw,
		},
		Sonar: (float64(data.Ultrasound&0x7FFF) - c.SonarZero) / c.SonarScale,
	}
}
