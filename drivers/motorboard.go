package drivers

import (
	"fmt"
	"github.com/felixge/godrone/log"
	"github.com/felixge/godrone/util"
	"os"
	"sync"
	"time"
)

const PWM_MAX = float64(511)

type Motorboard struct {
	file        *os.File
	pwms        [4]int
	leds        [4]LedColor
	ledsChanged bool
	mutex       sync.RWMutex
	log         log.Logger
	timer       *util.LoopTimer
}

func NewMotorboard(ttyPath string, log log.Logger) (*Motorboard, error) {
	timer := util.NewLoopTimer("motorboard", log)
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
	// navboard runs at ~200HZ
	interval := time.Second / 400
	for {
		start := time.Now()
		m.timer.Tick()
		m.mutex.RLock()
		m.updateSpeeds()
		if m.ledsChanged {
			m.updateLeds()
			m.ledsChanged = false
		}
		m.mutex.RUnlock()

		sleep := interval - time.Since(start)
		if sleep > 0 {
			time.Sleep(sleep)
		}
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

func (m *Motorboard) Speed(motorId int) (float64, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if motorId >= len(m.pwms) {
		return 0, fmt.Errorf("unknown motor: %d", motorId)
	}

	return float64(m.pwms[motorId]) / PWM_MAX, nil
}

func (m *Motorboard) SetSpeed(motorId int, speed float64) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if motorId >= len(m.pwms) {
		return fmt.Errorf("unknown motor: %d", motorId)
	}

	m.pwms[motorId] = int(speed * PWM_MAX)
	return nil
}

func (m *Motorboard) MotorCount() int {
	return len(m.pwms)
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
	cmd[0] = byte(0x20 | ((m.pwms[0] & 0x1ff) >> 4))
	cmd[1] = byte(((m.pwms[0] & 0x1ff) << 4) | ((m.pwms[1] & 0x1ff) >> 5))
	cmd[2] = byte(((m.pwms[1] & 0x1ff) << 3) | ((m.pwms[2] & 0x1ff) >> 6))
	cmd[3] = byte(((m.pwms[2] & 0x1ff) << 2) | ((m.pwms[3] & 0x1ff) >> 7))
	cmd[4] = byte(((m.pwms[3] & 0x1ff) << 1))
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
