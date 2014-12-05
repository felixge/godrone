package godrone

import (
	"fmt"
	"io"
	"strings"
	"time"
)

// We know the navboard runs at 200 Hz. Assuming that the navboard's
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
	return NewCustomFirmware(navboard, motorboard)
}

// NewCustomFirmware returns the same defaults as NewFirmware, but allows
// the caller to provide their own MotorLedWriter and NavdataReader
// implementations.
func NewCustomFirmware(n NavdataReader, m MotorLedWriter) (*Firmware, error) {
	return &Firmware{
		Navboard:   n,
		Motorboard: m,
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
			SonarGain: 0.2,
			SonarMax:  4,
		},
		// yes, more experimenting :)
		Controller: &Controller{
			RotationBand: 0.3,
			ThrottleMin:  0.4,
			Pitch:        PID{P: 0.02, I: 0.01, D: 0},
			Roll:         PID{P: 0.02, I: 0.01, D: 0},
			Yaw:          PID{P: 0.02, I: 0, D: 0},
			Altitude:     PID{P: 0.1, I: 0.2, D: 0.01},
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
	// Navboard holds the navboard interface.
	Navboard NavdataReader
	// Motorboard holds the motorboard interface.
	Motorboard MotorLedWriter
	// Calibrator holds the calibrator.
	Calibrator Calibrator
	// Calibration holds the current calibration.
	Calibration Calibration
	// Navdata holds the navdata read by the tick function.
	Navdata Navdata
	// Sensors holds the sensor values calculated by the Control function.
	Sensors Sensors
	// Filter estimates the placement of the drone based on navdata.
	Filter Filter
	// Controller tries to achieve the desired placement using the actuators.
	Controller *Controller
	// Motors holds the motor speeds applied by the Control function.
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

// Control calls Observe and uses the gathered data to manipulate the drones
// actuators in order to achieve the desired placement.
func (f *Firmware) Control() error {
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
	if mb, ok := f.Motorboard.(io.Closer); ok {
		if err := mb.Close(); err != nil {
			errors = append(errors, err.Error())
		}
	}
	if nb, ok := f.Navboard.(io.Closer); ok {
		if err := nb.Close(); err != nil {
			errors = append(errors, err.Error())
		}
	}
	if len(errors) == 0 {
		return nil
	}
	return fmt.Errorf(strings.Join(errors, ", "))
}
