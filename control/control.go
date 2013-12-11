package control

import (
	"github.com/felixge/godrone/attitude"
	"github.com/felixge/pidctrl"
	"math"
	"sync"
)

const (
	rotationMax = 0.15
	throttleMin = 0.5
	throttleFly = 0.66
	throttleMax = 1 - rotationMax
)

func NewControl(roll, pitch, yaw, altitude [3]float64) *Control {
	return &Control{
		roll:     pidctrl.NewPIDController(roll[0], roll[1], roll[2]),
		pitch:    pidctrl.NewPIDController(pitch[0], pitch[1], pitch[2]),
		yaw:      pidctrl.NewPIDController(yaw[0], yaw[1], yaw[2]),
		altitude: pidctrl.NewPIDController(altitude[0], altitude[1], altitude[2]),
	}
}

type Control struct {
	l        sync.Mutex
	roll     *pidctrl.PIDController
	pitch    *pidctrl.PIDController
	yaw      *pidctrl.PIDController
	altitude *pidctrl.PIDController
}

func (c *Control) Set(s attitude.Data) {
	c.l.Lock()
	defer c.l.Unlock()

	c.roll.Set(s.Roll)
	c.pitch.Set(s.Pitch)
	c.yaw.Set(s.Yaw)
	c.altitude.Set(s.Altitude)
}

func (c *Control) Update(a attitude.Data) (speeds [4]float64) {
	c.l.Lock()
	defer c.l.Unlock()

	var (
		adjRoll  = c.roll.Update(a.Roll)
		adjPitch = c.pitch.Update(a.Pitch)
		adjYaw   = c.yaw.Update(a.Yaw)
		throttle = throttleFly + c.altitude.Update(a.Altitude)
	)

	// kill motors if we want to land and are close enough to the ground 0.30 is
	// the min value the ultrasound can sense.
	if c.altitude.Get() <= 0 && a.Altitude <= 0.30 {
		return
	}

	throttle = clipToRange(throttle, throttleMin, throttleMax)

	speeds = [4]float64{
		throttle + clip(+adjRoll+adjPitch-adjYaw, rotationMax),
		throttle + clip(-adjRoll+adjPitch+adjYaw, rotationMax),
		throttle + clip(-adjRoll-adjPitch-adjYaw, rotationMax),
		throttle + clip(+adjRoll-adjPitch+adjYaw, rotationMax),
	}
	return
}

func clip(val, max float64) float64 {
	return math.Max(math.Min(val, max), -max)
}

func clipToRange(val, min, max float64) float64 {
	return math.Min(math.Max(val, min), max)
}
