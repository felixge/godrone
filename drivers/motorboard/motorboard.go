package motorboard

import (
	"os"
	"sync"
)

type Motorboard struct {
	lock   sync.RWMutex
	speeds [4]float64
	file   *os.File
	writer *writer
	leds   [4]LedColor
}

func NewMotorboard(tty string) (*Motorboard, error) {
	m := &Motorboard{}
	if err := m.open(tty); err != nil {
		return nil, err
	}
	return m, nil
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

func (m *Motorboard) SetSpeeds(speeds [4]float64) error {
	return m.writer.WriteSpeeds(speeds)
}

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

func (m *Motorboard) Close() error {
	return m.file.Close()
}
