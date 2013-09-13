package drivers

import (
	"fmt"
	"github.com/felixge/godrone/log"
	"os"
	"sync"
)

type Motorboard struct {
	file        *os.File
	speeds      [4]int
	leds        [4]LedColor
	ledsChanged bool
	mutex       sync.RWMutex
	log         log.Logger
	timer       *loopTimer
}

func NewMotorboard(ttyPath string, log log.Logger) (*Motorboard, error) {
	timer := newLoopTimer("motorboard", log)
	motorboard := &Motorboard{log: log, timer: timer}
	err := motorboard.open(ttyPath)
	if err != nil {
		return nil, err
	}
	go motorboard.loop()
	return motorboard, nil
}

func (m *Motorboard) open(path string) error {
	file, err := os.OpenFile(path, os.O_RDWR, 0)
	if err != nil {
		return err
	}
	m.file = file
	return nil
}

func (m *Motorboard) loop() {
	for {
		m.timer.Tick()
		m.mutex.RLock()
		m.updateSpeeds()
		if m.ledsChanged {
			m.updateLeds()
			m.ledsChanged = false
		}
		m.mutex.RUnlock()
	}
}

func (m *Motorboard) SetLed(led int, color LedColor) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.leds[led] = color
	m.ledsChanged = true
}

func (m *Motorboard) SetLeds(color LedColor) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for i := 0; i < len(m.leds); i++ {
		m.leds[i] = color
	}
	m.ledsChanged = true
}

func (m *Motorboard) Speed(motorId int) (int, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if motorId >= len(m.speeds) {
		return 0, fmt.Errorf("unknown motor: %d", motorId)
	}

	return m.speeds[motorId], nil
}

func (m *Motorboard) SetSpeed(motorId int, speed int) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if motorId >= len(m.speeds) {
		return fmt.Errorf("unknown motor: %d", motorId)
	}

	m.speeds[motorId] = speed
	return nil
}

func (m *Motorboard) MotorCount() int {
	return len(m.speeds)
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
