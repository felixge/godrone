package godrone

import "fmt"

type Placement struct {
	// Orientation in deg
	PRY
	// Altitude in meters
	Altitude float64
}

func (p Placement) String() string {
	return fmt.Sprintf("%s %7.2f A", p.PRY.String(), p.Altitude)
}

func (p1 Placement) Equal(p2 Placement) bool {
	return p1.PRY.Equal(p2.PRY) && p1.Altitude == p2.Altitude
}
