package attitude

import (
	"github.com/felixge/godrone/imu"
	"math"
	"time"
)

const rad2Deg = 180 / math.Pi

// NewComplementary returns a new Complementary filter.
func NewComplementary() *Complementary {
	return &Complementary{gGain: 0.98, aGain: 0.02}

}

// Implements a simple complementation filter.
// see http://www.pieter-jan.com/node/11
type Complementary struct {
	state   Attitude
	updated time.Time
	aGain   float64
	gGain   float64
}

// Update processes the given imu.Data and returns the new Attitude.
func (c *Complementary) Update(d imu.Data) Attitude {
	var (
		accel Attitude
		now   = time.Now()
	)

	accel.Pitch = angle(d.Ax, d.Az)
	accel.Roll = angle(d.Ay, d.Az)

	if !c.updated.IsZero() {
		dt := now.Sub(c.updated).Seconds()
		c.state.Roll += d.Gx * dt
		c.state.Pitch += d.Gy * dt
		c.state.Yaw += d.Gz * dt
	}

	c.state.Pitch = c.state.Pitch*c.gGain + accel.Pitch*c.aGain
	c.state.Roll = c.state.Roll*c.gGain + accel.Roll*c.aGain
	c.state.Altitude = d.UsAltitude

	c.updated = now
	return c.state
}

func angle(a, b float64) float64 {
	r := math.Atan2(a, b) * rad2Deg
	if r > 0 {
		r = 180 - r
	} else if r < 0 {
		r = -180 - r
	}
	return r
}
