package godrone

import (
	"github.com/felixge/godrone/drivers/navboard"
	"github.com/felixge/godrone/log"
	"time"
)

type Config struct {
	Log                 log.Interface
	Navboard            *navboard.Navboard
	CalibrationDuration time.Duration
}

type stopCmd chan error
type startCmd chan error

func NewFirmware(c Config) *Firmware {
	f := &Firmware{
		log:       c.Log,
		navboard:  c.Navboard,
		cmdCh:     make(chan interface{}),
		navdataCh: make(chan navboard.Data),
	}
	go f.mainLoop()
	return f
}

type Firmware struct {
	log       log.Interface
	navboard  *navboard.Navboard
	cmdCh     chan interface{}
	navdataCh chan navboard.Data
	started   bool
}

func (f *Firmware) Start() error {
	f.log.Info("Starting firmware")
	start := make(startCmd)
	f.cmdCh <- start
	return <-start
}

func (f *Firmware) start(cmd startCmd) {
	if f.started {
		cmd <- f.log.Error("Can't start, already started")
		return
	}

	if err := f.navboard.Calibrate(); err != nil {
		cmd<-err
		return
	}
	//go f.navboardLoop()

	f.started = true
	f.log.Info("Started firmware")
	cmd <- nil
}

func (f *Firmware) Stop() error {
	f.log.Info("Stopping firmware")
	stop := make(stopCmd)
	f.cmdCh <- stop
	return <-stop
}

func (f *Firmware) stop(cmd stopCmd) {
	if !f.started {
		cmd <- f.log.Error("Can't stop, already stopped")
		return
	}

	f.started = false
	f.log.Info("Stopped firmware")
	cmd <- nil
}

func (f *Firmware) Calibrate() error {
	return nil
}

func (f *Firmware) Takeoff() error {
	return nil
}

func (f *Firmware) Land() error {
	return nil
}

func (f *Firmware) mainLoop() {
	defer f.log.Error("Exiting main loop")

	//ahrsCh := f.ahrs.Sub(1)
	//defer f.ahrs.Unsub(ahrsCh)

	for {
		select {
		case data := <-f.navdataCh:
			f.log.Debug("navdata=%+v", data)
		case cmd := <-f.cmdCh:
			f.cmd(cmd)
		}
	}
}

func (f *Firmware) cmd(cmd interface{}) {
	f.log.Debug("Received cmd=%T", cmd)
	switch t := cmd.(type) {
	case startCmd:
		f.start(t)
	case stopCmd:
		f.stop(t)
	default:
		f.log.Error("Unknown cmd=%T", cmd)
	}
}

func (f *Firmware) navboardLoop() {
	for {
		data, err := f.navboard.NextData()
		if err != nil {
			duration := time.Second
			f.log.Debug("Sleeping after navboard error. duration=%s", duration)
			time.Sleep(duration)
		}
		select {
		case f.navdataCh <- data:
		default:
		}
	}
}
