package drivers

import (
	"github.com/felixge/godrone/log"
	"time"
)

type loopTimer struct {
	name     string
	lastHz   time.Time
	lastTick time.Time
	counter  int
	min      time.Duration
	max      time.Duration
	log      log.Logger
}

func newLoopTimer(name string, log log.Logger) *loopTimer {
	return &loopTimer{
		name: name,
		log:  log,
	}
}

func (l *loopTimer) Tick() {
	if l.lastHz.IsZero() {
		l.lastHz = time.Now()
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
	l.lastTick = time.Now()

	if time.Since(l.lastHz) >= time.Second {
		hz := float64(l.counter) / time.Since(l.lastHz).Seconds()
		l.log.Debug("%s hz: %f (min: %s, max: %s)", l.name, hz, l.min, l.max)
		l.counter = 0
		l.min = 0
		l.max = 0
		l.lastHz = time.Now()
	}
	l.counter++

}
