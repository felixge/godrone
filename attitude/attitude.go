package attitude

import (
	"fmt"
)

// Attitude describes the drones position relative to the ground.
type Attitude struct {
	Roll     float64 // degree
	Pitch    float64 // degree
	Yaw      float64 // degree
	Altitude float64 // meters
}

// String returns a human readable version of the Attitude.
func (d Attitude) String() string {
	return fmt.Sprintf(
		"Roll: %+6.1f Pitch: %+6.1f Yaw: %+6.1f Altitude: %+6.1f",
		d.Roll,
		d.Pitch,
		d.Yaw,
		d.Altitude,
	)
}
