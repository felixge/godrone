package main

type State struct {
	// Orientation in deg
	Orientation PRY
	// Altitude in meters
	Altitude float64
	// Fly determines whether the drone should cut of the engines when very close
	// to the ground or not.
	Fly bool
}
