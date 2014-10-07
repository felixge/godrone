package main

import (
	"math"
	"time"
)

const (
	rad2Deg = 180 / math.Pi
)

var (
	// found by experimenting
	aZeros       = PRY{Pitch: 2032, Roll: 2080, Yaw: 2020}
	gZeros       = PRY{Pitch: -21, Roll: 12, Yaw: 10}
	gSensitivity = PRY{Pitch: -16, Roll: 16, Yaw: 16}
)

// NewComplementary returns a new Complementary filter.
func NewComplementary() Complementary {
	return Complementary{gGain: 0.98, aGain: 0.02}

}

// Implements a simple complementation filter.
// see http://www.pieter-jan.com/node/11
type Complementary struct {
	aGain float64
	gGain float64
}

func (c *Complementary) Update(state *PRY, data Navdata, dt time.Duration) {
	var (
		dts = dt.Seconds()
		// assuming same sensitivity for all accelerometers
		aNormal = PRY{
			Pitch: float64(data.APitch) - aZeros.Pitch,
			Roll:  float64(data.ARoll) - aZeros.Roll,
			Yaw:   float64(data.AYaw) - aZeros.Yaw,
		}
		aDeg = PRY{
			Pitch: degAngle(aNormal.Roll, aNormal.Yaw),
			Roll:  degAngle(aNormal.Pitch, aNormal.Yaw),
		}
		gNormal = PRY{
			Pitch: (float64(data.GPitch) - gZeros.Pitch) / gSensitivity.Pitch,
			Roll:  (float64(data.GRoll) - gZeros.Roll) / gSensitivity.Roll,
			Yaw:   (float64(data.GYaw) - gZeros.Yaw) / gSensitivity.Yaw,
		}
		gDeg = PRY{
			Pitch: state.Pitch + (gNormal.Pitch * dts),
			Roll:  state.Roll + (gNormal.Roll * dts),
			Yaw:   state.Yaw + (gNormal.Yaw * dts),
		}
	)
	state.Pitch = gDeg.Pitch*c.gGain + aDeg.Pitch*c.aGain
	state.Roll = gDeg.Roll*c.gGain + aDeg.Roll*c.aGain
	state.Yaw = gDeg.Yaw
}

func degAngle(a, b float64) float64 {
	return math.Atan2(a, b) * rad2Deg
}
