package motorboard

import (
	"os"
	"sync"
	"time"
)

const DefaultTTYPath = "/dev/ttyO0"

type Driver struct {
	file        *os.File
	Speeds      [4]int
	leds        [4]LedColor
	ledsChanged bool
	mutex       sync.Mutex
}

func NewDriver(ttyPath string) (*Driver, error) {
	driver := &Driver{}
	err := driver.open(ttyPath)
	if err != nil {
		return nil, err
	}
	go driver.loop()
	return driver, nil
}

func (c *Driver) open(path string) error {
	file, err := os.OpenFile(path, os.O_RDWR, 0)
	if err != nil {
		return err
	}
	c.file = file
	return nil
}

func (c *Driver) loop() {
	hz := 200
	sleepTime := (1000 / time.Duration(hz)) * time.Millisecond

	for {
		c.mutex.Lock()
		if c.ledsChanged {
			c.updateLeds()
			c.ledsChanged = false
		}
		c.mutex.Unlock()

		time.Sleep(sleepTime)
	}
}

func (c *Driver) SetLed(led int, color LedColor) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.leds[led] = color
	c.ledsChanged = true
}

func (c *Driver) SetLeds(color LedColor) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for i := 0; i < len(c.leds); i++ {
		c.leds[i] = color
	}
	c.ledsChanged = true
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

func (m *Driver) updateSpeeds() error {
	_, err := m.file.Write(m.pwmCmd())
	return err
}

func (m *Driver) updateLeds() error {
	_, err := m.file.Write(m.ledCmd())
	return err
}
