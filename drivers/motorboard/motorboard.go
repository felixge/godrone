package motorboard

import (
	"os"
)

// NewMotorboard returns a new Motorboard, or an error if opening the tty file
// failed.
func NewMotorboard(tty string) (*Motorboard, error) {
	m := &Motorboard{}
	if err := m.open(tty); err != nil {
		return nil, err
	}
	return m, nil
}

// Motorboard implements a motorboard driver for the Parrot AR Drone 2.0. It
// must be used from a single goroutine.
type Motorboard struct {
	speeds [4]float64
	file   *os.File
	writer *writer
	leds   [4]LedColor
}

func (m *Motorboard) open(tty string) (err error) {
	if m.file != nil {
		return
	}

	m.file, err = os.OpenFile(tty, os.O_RDWR, 0)
	if err != nil {
		return
	}
	m.writer = newWriter(m.file)
	return
}

// SetSpeeds updates the motors with the given speeds. This method has to be
// called frequently (usually at the same rate sensor data is read from the
// navboard), otherwise the motors will stop.
func (m *Motorboard) SetSpeeds(speeds [4]float64) error {
	return m.writer.WriteSpeeds(speeds)
}

// SetLeds changes the colors of the LEDs below the motors.
func (m *Motorboard) SetLeds(leds [4]LedColor) (err error) {
	if leds == m.leds {
		return
	}
	err = m.writer.WriteLeds(leds)
	if err != nil {
		m.leds = leds
	}
	return
}

// Close closes the underlaying tty file.
func (m *Motorboard) Close() error {
	return m.file.Close()
}
