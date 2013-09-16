package util

import (
	"github.com/felixge/godrone/log"
	"time"
)

type LoopTimer struct {
	name     string
	lastLog  time.Time
	lastTick time.Time
	counter  int
	min      time.Duration
	max      time.Duration
	log      log.Logger
}

func NewLoopTimer(name string, log log.Logger) *LoopTimer {
	return &LoopTimer{
		name: name,
		log:  log,
	}
}

func (l *LoopTimer) Tick() {
	if l.lastLog.IsZero() {
		l.lastLog = time.Now()
	}

	if !l.lastTick.IsZero() {
		tickDuration := time.Since(l.lastTick)
		if tickDuration < l.min || l.min == 0 {
			l.min = tickDuration
		}
		if tickDuration > l.max {
			l.max = tickDuration
		}
	}

	if time.Since(l.lastLog) >= 10*time.Second {
		hz := float64(l.counter) / time.Since(l.lastLog).Seconds()
		l.log.Debug("%s hz: %f (min: %s, max: %s)", l.name, hz, l.min, l.max)
		l.counter = 0
		l.min = 0
		l.max = 0
		l.lastLog = time.Now()
	}
	l.counter++
	l.lastTick = time.Now()
}
