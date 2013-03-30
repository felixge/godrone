package motorboard

import (
	"os"
)

type Driver struct {
	file   *os.File
	Speeds [4]int
	leds   []LedColor
}

func NewDriver() (*Driver, error) {
	driver := &Driver{
		leds: make([]LedColor, 4),
	}
	err := driver.Open("/dev/ttyO0")
	if err != nil {
		return nil, err
	}
	return driver, nil
}

func (c *Driver) Open(path string) error {
	file, err := os.OpenFile(path, os.O_RDWR, 0)
	if err != nil {
		return err
	}
	c.file = file
	return nil
}

func (c *Driver) SetLed(motorId int, color LedColor) {
	c.leds[motorId] = color
}

type LedColor int

const (
	LedOff    LedColor = iota
	LedRed             = 1
	LedGreen           = 2
	LedOrange          = 3
)

// cmd = 011rrrrx xxxggggx (used to be 011grgrg rgrxxxxx in AR Drone 1.0)
// see: https://github.com/ardrone/ardrone/blob/master/ardrone/motorboard/motorboard.c#L243
func (m *Driver) ledCmd() []byte {
	cmd := make([]byte, 2)
	cmd[0] = 0x60

	for i, color := range m.leds {
		if color == LedRed || color == LedOrange {
			cmd[0] = cmd[0] | (1 << (byte(i) + 1))
		}

		if color == LedGreen || color == LedOrange {
			cmd[1] = cmd[1] | (1 << (byte(i) + 1))
		}
	}
	return cmd
}

// see: https://github.com/ardrone/ardrone/blob/master/ardrone/motorboard/motorboard.c
func (m *Driver) pwmCmd() []byte {
	cmd := make([]byte, 5)
	cmd[0] = byte(0x20 | ((m.Speeds[0] & 0x1ff) >> 4))
	cmd[1] = byte(((m.Speeds[0] & 0x1ff) << 4) | ((m.Speeds[1] & 0x1ff) >> 5))
	cmd[2] = byte(((m.Speeds[1] & 0x1ff) << 3) | ((m.Speeds[2] & 0x1ff) >> 6))
	cmd[3] = byte(((m.Speeds[2] & 0x1ff) << 2) | ((m.Speeds[3] & 0x1ff) >> 7))
	cmd[4] = byte(((m.Speeds[3] & 0x1ff) << 1))
	return cmd
}

func (m *Driver) UpdateSpeeds() error {
	_, err := m.file.Write(m.pwmCmd())
	return err
}

func (m *Driver) UpdateLeds() error {
	_, err := m.file.Write(m.ledCmd())
	return err
}
