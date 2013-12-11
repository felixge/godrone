package attitude

import (
	"fmt"
)

type Data struct {
	Roll     float64 // degree
	Pitch    float64 // degree
	Yaw      float64 // degree
	Altitude float64 // meters
}

func (d Data) String() string {
	return fmt.Sprintf(
		"Roll: %+6.1f Pitch: %+6.1f Yaw: %+6.1f Altitude: %+6.1f",
		d.Roll,
		d.Pitch,
		d.Yaw,
		d.Altitude,
	)
}
