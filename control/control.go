package control

import (
	"github.com/felixge/godrone/attitude"
	"github.com/felixge/pidctrl"
	"math"
	"sync"
)

const (
	throttleMax = 0.85
	adjMax      = 1 - throttleMax
)

func NewControl(rollPID, pitchPID, yawPID [3]float64) *Control {
	return &Control{
		roll:  pidctrl.NewPIDController(rollPID[0], rollPID[1], rollPID[2]),
		pitch: pidctrl.NewPIDController(pitchPID[0], pitchPID[1], pitchPID[2]),
		yaw:   pidctrl.NewPIDController(yawPID[0], yawPID[1], yawPID[2]),
	}
}

type Control struct {
	l                sync.Mutex
	throttle         float64
	roll, pitch, yaw *pidctrl.PIDController
}

func (c *Control) Set(s attitude.Data, throttle float64) {
	c.l.Lock()
	defer c.l.Unlock()

	c.roll.Set(s.Roll)
	c.pitch.Set(s.Pitch)
	c.yaw.Set(s.Yaw)
	c.throttle = throttle
}

func (c *Control) Update(a attitude.Data) (speeds [4]float64) {
	c.l.Lock()
	defer c.l.Unlock()

	var (
		adjRoll  = c.roll.Update(a.Roll)
		adjPitch = c.pitch.Update(a.Pitch)
		//adjYaw   = c.yaw.Update(a.Yaw)
		adjYaw   = 0.0
		throttle = math.Min(c.throttle, 1) * throttleMax
	)

	if throttle <= 0 {
		return
	}

	speeds = [4]float64{
		throttle + math.Max(math.Min(+adjRoll+adjPitch+adjYaw, adjMax), -adjMax),
		throttle + math.Max(math.Min(-adjRoll+adjPitch-adjYaw, adjMax), -adjMax),
		throttle + math.Max(math.Min(-adjRoll-adjPitch+adjYaw, adjMax), -adjMax),
		throttle + math.Max(math.Min(+adjRoll-adjPitch-adjYaw, adjMax), -adjMax),
	}
	return
}
