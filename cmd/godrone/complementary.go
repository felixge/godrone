package main

import (
	"math"
	"time"
)

const (
	rad2Deg = 180 / math.Pi
)

// Implements a simple complementation filter.
// see http://www.pieter-jan.com/node/11
type Complementary struct {
	AccGain  float64
	GyroGain float64
}

func (c Complementary) Update(state *PRY, acc, gyro PRY, dt time.Duration) {
	if c.AccGain+c.GyroGain != 1 {
		panic("Gains must add up to 1")
	}
	var (
		dts    = dt.Seconds()
		accDeg = PRY{
			Pitch: degAngle(acc.Roll, acc.Yaw),
			Roll:  degAngle(acc.Pitch, acc.Yaw),
		}
		gyroDeg = PRY{
			Pitch: state.Pitch + (gyro.Pitch * dts),
			Roll:  state.Roll + (gyro.Roll * dts),
			Yaw:   state.Yaw + (gyro.Yaw * dts),
		}
	)
	state.Pitch = gyroDeg.Pitch*c.GyroGain + accDeg.Pitch*c.AccGain
	state.Roll = gyroDeg.Roll*c.GyroGain + accDeg.Roll*c.AccGain
	state.Yaw = gyroDeg.Yaw
}

func degAngle(a, b float64) float64 {
	return math.Atan2(a, b) * rad2Deg
}
