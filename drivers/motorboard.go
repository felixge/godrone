package drivers

import (
	"fmt"
	"os"
	"sync"
	"time"
)

type Motorboard struct {
	file        *os.File
	speeds      [4]int
	leds        [4]LedColor
	ledsChanged bool
	mutex       sync.RWMutex
}

func NewMotorboard(ttyPath string) (*Motorboard, error) {
	driver := &Motorboard{}
	err := driver.open(ttyPath)
	if err != nil {
		return nil, err
	}
	go driver.loop()
	return driver, nil
}

func (c *Motorboard) open(path string) error {
	file, err := os.OpenFile(path, os.O_RDWR, 0)
	if err != nil {
		return err
	}
	c.file = file
	return nil
}

func (c *Motorboard) loop() {
	hz := 200
	sleepTime := (1000 / time.Duration(hz)) * time.Millisecond

	for {
		c.mutex.Lock()
		c.updateSpeeds()
		if c.ledsChanged {
			c.updateLeds()
			c.ledsChanged = false
		}
		c.mutex.Unlock()

		time.Sleep(sleepTime)
	}
}

func (c *Motorboard) SetLed(led int, color LedColor) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.leds[led] = color
	c.ledsChanged = true
}

func (c *Motorboard) SetLeds(color LedColor) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for i := 0; i < len(c.leds); i++ {
		c.leds[i] = color
	}
	c.ledsChanged = true
}

func (c *Motorboard) SetSpeeds(speed int) error {
	for i := 0; i < len(c.speeds); i++ {
		if err := c.SetSpeed(i, speed); err != nil {
			return err
		}
	}
	return nil
}

func (c *Motorboard) Speed(motorId int) (int, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if motorId >= len(c.speeds) {
		return 0, fmt.Errorf("unknown motor: %d", motorId)
	}

	return c.speeds[motorId], nil
}

func (c *Motorboard) SetSpeed(motorId int, speed int) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if motorId >= len(c.speeds) {
		return fmt.Errorf("unknown motor: %d", motorId)
	}

	c.speeds[motorId] = speed
	return nil
}


func (c *Motorboard) MotorCount() int {
	return len(c.speeds)
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
func (m *Motorboard) ledCmd() []byte {
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
func (m *Motorboard) pwmCmd() []byte {
	cmd := make([]byte, 5)
	cmd[0] = byte(0x20 | ((m.speeds[0] & 0x1ff) >> 4))
	cmd[1] = byte(((m.speeds[0] & 0x1ff) << 4) | ((m.speeds[1] & 0x1ff) >> 5))
	cmd[2] = byte(((m.speeds[1] & 0x1ff) << 3) | ((m.speeds[2] & 0x1ff) >> 6))
	cmd[3] = byte(((m.speeds[2] & 0x1ff) << 2) | ((m.speeds[3] & 0x1ff) >> 7))
	cmd[4] = byte(((m.speeds[3] & 0x1ff) << 1))
	return cmd
}

func (m *Motorboard) updateSpeeds() error {
	_, err := m.file.Write(m.pwmCmd())
	return err
}

func (m *Motorboard) updateLeds() error {
	_, err := m.file.Write(m.ledCmd())
	return err
}
