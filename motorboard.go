package godrone

import "os"

const (
	// pwmMax is the max speed the motors can be set to
	pwmMax = 511

	// Commands
	setSpeeds = 0x20
	setLeds   = 0x60
)

// OpenMotorboard returns a new Motorboard driver.
func OpenMotorboard(tty string) (*Motorboard, error) {
	file, err := os.OpenFile(tty, os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}
	return &Motorboard{file: file}, nil
}

// Motorboard implements a motorboard driver for the Parrot AR Drone 2.0. It
// must be used from a single goroutine.
type Motorboard struct {
	file *os.File
}

// WriteSpeeds writes the command for updating the motor speeds. Speeds must be
// int the [0,1] range.
func (m *Motorboard) WriteSpeeds(speeds [4]float64) error {
	var pwms [4]uint16
	for i, speed := range speeds {
		pwms[i] = uint16(speed * pwmMax)
	}

	// see: https://github.com/ardrone/ardrone/blob/master/ardrone/motorboard/motorboard.c
	cmd := []byte{
		byte(setSpeeds | ((pwms[0] & 0x1ff) >> 4)),
		byte(((pwms[0] & 0x1ff) << 4) | ((pwms[1] & 0x1ff) >> 5)),
		byte(((pwms[1] & 0x1ff) << 3) | ((pwms[2] & 0x1ff) >> 6)),
		byte(((pwms[2] & 0x1ff) << 2) | ((pwms[3] & 0x1ff) >> 7)),
		byte(((pwms[3] & 0x1ff) << 1)),
	}
	_, err := m.file.Write(cmd)
	return err
}

// WriteLeds writes the command for updating the LEDs.
// cmd = 011rrrrx xxxggggx (used to be 011grgrg rgrxxxxx in AR Drone 1.0)
// see: https://github.com/ardrone/ardrone/blob/master/ardrone/motorboard/motorboard.c#L243
func (m *Motorboard) WriteLeds(leds [4]LedColor) error {
	cmd := make([]byte, 2)
	cmd[0] = setLeds

	for i, color := range leds {
		if color == LedRed || color == LedOrange {
			cmd[0] = cmd[0] | (1 << (byte(i) + 1))
		}

		if color == LedGreen || color == LedOrange {
			cmd[1] = cmd[1] | (1 << (byte(i) + 1))
		}
	}
	_, err := m.file.Write(cmd)
	return err
}

// Close closes the underlaying tty file.
func (m *Motorboard) Close() error {
	return m.file.Close()
}

// LedColor holds a LED color constant.
type LedColor int

// The available colors for the LEDs
const (
	LedOff    LedColor = iota
	LedRed             = 1
	LedGreen           = 2
	LedOrange          = 3
)

// Leds is a helper method to return an array of LedColor's all set to the same
// value.
func Leds(c LedColor) [4]LedColor {
	return [4]LedColor{c, c, c, c}
}
