package godrone

import "math"

import "time"

const (
	rotationBand  = 0.3
	throttleHover = 0.45
	throttleMin   = 0.1
	throttleMax   = 1 - rotationBand
)

type Controller struct {
	Pitch    PID
	Roll     PID
	Yaw      PID
	Altitude PID
}

func (c *Controller) Control(actual, desired Placement, dt time.Duration) [4]float64 {
	var speeds [4]float64
	var (
		pitchOut = c.Roll.Update(actual.Pitch, desired.Pitch, dt)
		rollOut  = c.Roll.Update(actual.Roll, desired.Roll, dt)
		yawOut   = c.Yaw.Update(actual.Yaw, desired.Yaw, dt)
		altOut   = c.Altitude.Update(actual.Altitude, desired.Altitude, dt)
	)
	throttle := math.Max(throttleMin, math.Min(throttleMax, throttleHover+altOut))
	speeds = [4]float64{
		throttle + clipBand(+rollOut+pitchOut+yawOut, rotationBand),
		throttle + clipBand(-rollOut+pitchOut-yawOut, rotationBand),
		throttle + clipBand(-rollOut-pitchOut+yawOut, rotationBand),
		throttle + clipBand(+rollOut-pitchOut-yawOut, rotationBand),
	}
	return speeds
}

func clipBand(val, band float64) float64 {
	return band/2 + clip(val, band/2)
}

func clip(val, max float64) float64 {
	return math.Max(math.Min(val, max), -max)
}
