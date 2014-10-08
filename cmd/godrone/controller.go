package main

import "math"

import "time"

const (
	rotationBand  = 0.3
	throttleHover = 0.45
	throttleMin   = 0.1
	throttleMax   = 1 - rotationBand
)

type Controller struct {
	Pitch        PID
	Roll         PID
	Yaw          PID
	Altitude     PID
	LandAltitude float64
}

func (c *Controller) Apply(desired State, actual *State, dt time.Duration, motorboard *Motorboard) error {
	var speeds [4]float64
	if desired.Fly || !desired.Fly && actual.Altitude > c.LandAltitude {
		var (
			a        = actual
			d        = desired
			pitchOut = c.Roll.Update(a.Orientation.Pitch, d.Orientation.Pitch, dt)
			rollOut  = c.Roll.Update(a.Orientation.Roll, d.Orientation.Roll, dt)
			yawOut   = c.Yaw.Update(a.Orientation.Yaw, d.Orientation.Yaw, dt)
			altOut   = c.Altitude.Update(a.Altitude, d.Altitude, dt)
		)
		throttle := math.Max(throttleMin, math.Min(throttleMax, throttleHover+altOut))
		speeds = [4]float64{
			throttle + clipBand(+rollOut+pitchOut+yawOut, rotationBand),
			throttle + clipBand(-rollOut+pitchOut-yawOut, rotationBand),
			throttle + clipBand(-rollOut-pitchOut+yawOut, rotationBand),
			throttle + clipBand(+rollOut-pitchOut-yawOut, rotationBand),
		}
		actual.Fly = true
	} else {
		actual.Fly = false
	}
	if err := motorboard.WriteSpeeds(speeds); err != nil {
		return err
	}
	return nil
}

func clipBand(val, band float64) float64 {
	return band/2 + clip(val, band/2)
}

func clip(val, max float64) float64 {
	return math.Max(math.Min(val, max), -max)
}
