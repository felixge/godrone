package attitude

import (
	"fmt"
)

type Data struct {
	Roll  float64
	Pitch float64
	Yaw   float64
}

func (d Data) String() string {
	return fmt.Sprintf(
		"Roll: %+6.1f Pitch: %+6.1f Yaw: %+6.1f",
		d.Roll,
		d.Pitch,
		d.Yaw,
	)
}
