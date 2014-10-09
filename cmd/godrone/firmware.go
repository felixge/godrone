package main

import (
	"fmt"
	"strings"
	"time"
)

// We know the navboard runs at 200 Hz. Assuming that the navboards
// internal clock is more accurate than the system clock, we should get
// better precision from hardcoding this assumption.
var dt = time.Second / 200

// NewFirmware returns a firmware with a reasonable default configuration. The
// returned firmware itself does not do anything unless the user calls methods
// on it.
func NewFirmware() (*Firmware, error) {
	navboard, err := OpenNavboard("/dev/ttyO1")
	if err != nil {
		return nil, fmt.Errorf("Failed to open navboard: %s", err)
	}
	motorboard, err := OpenMotorboard("/dev/ttyO0")
	if err != nil {
		return nil, fmt.Errorf("Failed to open navboard: %s", err)
	}
	return &Firmware{
		Navboard:   navboard,
		Motorboard: motorboard,
		// found by experimenting :)
		// overwritten by calibration (except for gyro/sonar scale)
		Calibration: Calibration{
			AccZeros:   PRY{Pitch: 2038, Roll: 2070, Yaw: 2048},
			AccScale:   PRY{Pitch: 1, Roll: 1, Yaw: 1},
			GyroZeros:  PRY{Pitch: 16, Roll: -27, Yaw: 1.5},
			GyroScale:  PRY{Pitch: -16, Roll: 16, Yaw: 16},
			SonarScale: 3500,
		},
		Calibrator: Calibrator{
			Samples:   200,
			MaxStdDev: 5,
		},
		// also found by experimenting :)
		Filter: Filter{
			GyroGain:  0.98,
			AccGain:   0.02,
			SonarGain: 0.1,
			SonarMax:  4,
		},
		// yes, more experimenting :)
		Controller: &Controller{
			Pitch:    PID{P: 0.02, I: 0.0001, D: 0},
			Roll:     PID{P: 0.02, I: 0.0001, D: 0},
			Yaw:      PID{P: 0.02, I: 0, D: 0},
			Altitude: PID{P: 0.2, I: 0.05, D: 0.05},
		},
	}, nil
}

// Firmware provides an interface for writing custom drone firmwares. It does
// not support concurrent access, and offers many ways for the user to shoot
// himself in the foot when modifying anything but the Desired state.
//
// @TODO Many interesting things could be achieved by turning some of the
// fields below into integer types.
type Firmware struct {
	// Navboard holds the navboard.
	Navboard *Navboard
	// Motorboard holds the motorboard.
	Motorboard *Motorboard
	// Calibrator holds the calibrator.
	Calibrator Calibrator
	// Calibration holds the current calibration.
	Calibration Calibration
	// Navdata holds the navdata read by the tick function.
	Navdata Navdata
	// Sensors holds the sensor values calculated by the Fly function
	Sensors Sensors
	// Filter estimates the placement of the drone based on navdata.
	Filter Filter
	// Controller tries to achieve the desired placement using the actuators.:w
	Controller *Controller
	// Motors holds the motor speeds applied by the Fly function
	Motors [4]float64
	// Actual holds the placement of the drone as estimated by the filter.
	Actual Placement
	// Desired holds the desired placement the controller is trying to achieve.
	Desired Placement
}

// Observe reads navdata and uses it to update sensor and placement estimates.
func (f *Firmware) Observe() error {
	var err error
	f.Navdata, err = f.Navboard.Read()
	if err != nil {
		return fmt.Errorf("Failed to read navdata: %s", err)
	}
	f.Sensors = f.Calibration.Convert(f.Navdata)
	f.Filter.Update(&f.Actual, f.Sensors, dt)
	return nil
}

// Fly calls Observe and uses the gathered data to manipulate the drones
// actuators in order to achieve the desired placement.
func (f *Firmware) Fly() error {
	if err := f.Observe(); err != nil {
		return err
	}
	f.Motors = f.Controller.Control(f.Actual, f.Desired, dt)
	if err := f.Motorboard.WriteSpeeds(f.Motors); err != nil {
		return fmt.Errorf("Failed to write speeds: %s", err)
	}
	return nil
}

// Calibrate can be called while the drone is on a level surface to estimate
// sensor calibration.
func (f *Firmware) Calibrate() error {
	return f.Calibrator.Calibrate(f.Navboard, &f.Calibration)
}

// Close cleans up resources allocated by the firmware.
func (f *Firmware) Close() error {
	var errors []string
	if err := f.Motorboard.Close(); err != nil {
		errors = append(errors, err.Error())
	}
	if err := f.Navboard.Close(); err != nil {
		errors = append(errors, err.Error())
	}
	if len(errors) == 0 {
		return nil
	}
	return fmt.Errorf(strings.Join(errors, ", "))
}
