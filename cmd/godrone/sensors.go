package main

type Sensors struct {
	// Acceleration in m/s^2
	Acc PRY
	// Rotation in deg/s
	Gyro PRY
}

type Calibration struct {
	AccZeros        PRY
	AccSensitivity  PRY
	GyroZeros       PRY
	GyroSensitivity PRY
}

func (t Calibration) Convert(data Navdata) Sensors {
	return Sensors{
		Acc: PRY{
			Pitch: (float64(data.AccPitch) - t.AccZeros.Pitch) / t.AccSensitivity.Pitch,
			Roll:  (float64(data.AccRoll) - t.AccZeros.Roll) / t.AccSensitivity.Roll,
			Yaw:   (float64(data.AccYaw) - t.AccZeros.Yaw) / t.AccSensitivity.Yaw,
		},
		Gyro: PRY{
			Pitch: (float64(data.GyroPitch) - t.GyroZeros.Pitch) / t.GyroSensitivity.Pitch,
			Roll:  (float64(data.GyroRoll) - t.GyroZeros.Roll) / t.GyroSensitivity.Roll,
			Yaw:   (float64(data.GyroYaw) - t.GyroZeros.Yaw) / t.GyroSensitivity.Yaw,
		},
	}
}
