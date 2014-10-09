package main

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
