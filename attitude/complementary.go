package attitude

import (
	"github.com/felixge/godrone/imu"
	"math"
)

type Complementary struct{}

const rad2Deg = 180 / math.Pi

func (c *Complementary) Update(d imu.Data) Data {
	a := Data{}
	a.Pitch = angle(d.Ax, d.Az)
	a.Roll = angle(d.Ay, d.Az)
	a.Yaw = angle(d.Ax, d.Ay)

	return a
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
