package godrone

import "time"

// PID implements a simple PID controller.
// see http://en.wikipedia.org/wiki/PID_controller#Pseudocode
type PID struct {
	P           float64
	I           float64
	D           float64
	prevDiff    float64
	prevDesired float64
	integral    float64
}

// UpdateDuration updates the controller with the given value and duration since
// the last update. It returns the new output.
func (p *PID) Update(actual, desired float64, dt time.Duration) float64 {
	var (
		dts        = dt.Seconds()
		diff       = desired - actual
		derivative float64
	)
	if desired != p.prevDesired {
		p.integral = 0
	}
	p.integral += diff * dts
	if dts > 0 {
		derivative = (diff - p.prevDiff) / dts
	}
	p.prevDiff = diff
	p.prevDesired = desired
	return p.P*diff + p.I*p.integral + p.D*derivative
}
