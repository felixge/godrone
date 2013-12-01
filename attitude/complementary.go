package attitude

import (
	"github.com/felixge/godrone/imu"
	"math"
	"time"
)

func NewComplementary() *Complementary {
	return &Complementary{gGain: 0.98, aGain: 0.02}

}

type Complementary struct {
	data    Data
	updated time.Time
	aGain   float64
	gGain   float64
}

const rad2Deg = 180 / math.Pi

func (c *Complementary) Update(d imu.Data) Data {
	var (
		a   Data
		now = time.Now()
	)

	a.Pitch = angle(d.Ax, d.Az)
	a.Roll = angle(d.Ay, d.Az)

	if !c.updated.IsZero() {
		dt := now.Sub(c.updated).Seconds()
		c.data.Roll += d.Gx * dt
		c.data.Pitch += d.Gy * dt
		c.data.Yaw += d.Gz * dt
	}

	c.data.Pitch = c.data.Pitch*c.gGain + a.Pitch*c.aGain
	c.data.Roll = c.data.Roll*c.gGain + a.Roll*c.aGain

	c.updated = now
	return c.data
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
