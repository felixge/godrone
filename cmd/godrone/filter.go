package main

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
	// SonarMin in m
	SonarMin float64
	// SonarMax in m
	SonarMax float64
}

func (f Filter) Update(state *State, sensors Sensors, dt time.Duration) {
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
			Pitch: state.Orientation.Pitch + (sensors.Gyro.Pitch * dts),
			Roll:  state.Orientation.Roll + (sensors.Gyro.Roll * dts),
			Yaw:   state.Orientation.Yaw + (sensors.Gyro.Yaw * dts),
		}
	)
	// Implements a simple complementation filter.
	// see http://www.pieter-jan.com/node/11
	state.Orientation.Pitch = gyroDeg.Pitch*f.GyroGain + accDeg.Pitch*f.AccGain
	state.Orientation.Roll = gyroDeg.Roll*f.GyroGain + accDeg.Roll*f.AccGain
	// @TODO Integrate gyro yaw with magotometer yaw
	state.Orientation.Yaw = gyroDeg.Yaw
	// The sonar seems to be unable to detect the ground when its very close, so
	// we treat all values below this min as 0.
	var sonar = sensors.Sonar
	if sensors.Sonar < f.SonarMin {
		sonar = 0
	}
	// The sonar sometimes reads very high values when on the ground. Ignoring
	// the sonar above a certain altitude solves the problem.
	// @TODO Use barometer above SonarMax
	if sensors.Sonar < f.SonarMax {
		state.Altitude += (sonar - state.Altitude) * f.SonarGain
	}
	if state.Altitude < 0 {
		state.Altitude = 0
	}
}

func degAngle(a, b float64) float64 {
	const rad2Deg = 180 / math.Pi
	return math.Atan2(a, b) * rad2Deg
}
