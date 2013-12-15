// Package control implements flight stablization.
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

// NewControl returns a new Control with the given roll, pitch and yaw PID
// values.
func NewControl(roll, pitch, yaw []float64) *Control {
	return &Control{
		roll:  pidctrl.NewPIDController(roll[0], roll[1], roll[2]),
		pitch: pidctrl.NewPIDController(pitch[0], pitch[1], pitch[2]),
		yaw:   pidctrl.NewPIDController(yaw[0], yaw[1], yaw[2]),
	}
}

// Control determines how to translate a given setpoint and measured attitude
// into motor output.
type Control struct {
	l        sync.Mutex
	roll     *pidctrl.PIDController
	pitch    *pidctrl.PIDController
	yaw      *pidctrl.PIDController
	throttle float64
}

// Set sets the desired attitude and throttle setpoint. For now it does not
// support altitude stabilization.
func (c *Control) Set(s attitude.Attitude, throttle float64) {
	c.l.Lock()
	defer c.l.Unlock()

	c.roll.Set(s.Roll)
	c.pitch.Set(s.Pitch)
	c.yaw.Set(s.Yaw)
	c.throttle = throttle
}

// Update takes the last measured attitude and converts it into motor outputs
// to correct it towards the setpoint.
func (c *Control) Update(a attitude.Attitude) (speeds [4]float64) {
	c.l.Lock()
	defer c.l.Unlock()

	var (
		adjRoll  = c.roll.Update(a.Roll)
		adjPitch = c.pitch.Update(a.Pitch)
		adjYaw   = c.yaw.Update(a.Yaw)
		throttle = clipToRange(c.throttle*throttleMax, 0, throttleMax)
	)

	if throttle == 0 {
		return
	}

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
