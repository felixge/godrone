package godrone

import (
	"math"
	"time"
)

type Filter struct {
	// AccGain as fraction of 1
	AccGain float64
	// GyroGain as fraction of 1
	GyroGain float64
	// SonarGain as fraction of 1
	SonarGain float64
	// SonarMax in m
	SonarMax float64
}

func (f Filter) Update(placement *Placement, sensors Sensors, dt time.Duration) {
	if f.AccGain+f.GyroGain != 1 {
		panic("Gains must add up to 1")
	}
	var (
		dts    = dt.Seconds()
		accDeg = PRY{
			Pitch: degAngle(sensors.Acc.Roll, sensors.Acc.Yaw),
			Roll:  degAngle(sensors.Acc.Pitch, sensors.Acc.Yaw),
		}
		gyroDeg = PRY{
			Pitch: placement.Pitch + (sensors.Gyro.Pitch * dts),
			Roll:  placement.Roll + (sensors.Gyro.Roll * dts),
			Yaw:   placement.Yaw + (sensors.Gyro.Yaw * dts),
		}
	)
	// Implements a simple complementation filter.
	// see http://www.pieter-jan.com/node/11
	placement.Pitch = gyroDeg.Pitch*f.GyroGain + accDeg.Pitch*f.AccGain
	placement.Roll = gyroDeg.Roll*f.GyroGain + accDeg.Roll*f.AccGain
	// @TODO Integrate gyro yaw with magotometer yaw
	placement.Yaw = gyroDeg.Yaw
	// The sonar sometimes reads very high values when on the ground. Ignoring
	// the sonar above a certain altitude solves the problem.
	// @TODO Use barometer above SonarMax
	if sensors.Sonar < f.SonarMax {
		placement.Altitude += (sensors.Sonar - placement.Altitude) * f.SonarGain
	}
}

func degAngle(a, b float64) float64 {
	const rad2Deg = 180 / math.Pi
	return math.Atan2(a, b) * rad2Deg
}
