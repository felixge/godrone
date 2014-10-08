package main

type Sensors struct {
	// Acceleration in m/s^2
	Acc PRY
	// Rotation in deg/s
	Gyro PRY
	// Altitude in m
	Sonar float64
}
