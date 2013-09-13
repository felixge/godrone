package godrone

import (
	"github.com/felixge/godrone/drivers"
	"github.com/felixge/godrone/log"
	"os"
)

func NewFirmware(c Config) (*Firmware, error) {
	log, err := log.NewLogger(c.LogLevel, c.LogTimeFormat, os.Stdout)
	if err != nil {
		return nil, err
	}

	log.Info("Initializing firmware")

	log.Debug("Initializing navboard: %s", c.NavboardTTY)
	navboard, err := drivers.NewNavboard(c.NavboardTTY)
	if err != nil {
		return nil, log.Emergency("Could not initialize navboard: %s", err)
	}
	log.Debug("Initialized navboard")

	log.Debug("Initializing motorboard: %s", c.MotorboardTTY)
	motorboard, err := drivers.NewMotorboard(c.MotorboardTTY)
	if err != nil {
		return nil, log.Emergency("Could not initialize motorboard: %s", err)
	}
	log.Debug("Initialized motorboard")

	firmware := &Firmware{
		config:     &c,
		log:        log,
		navboard:   navboard,
		motorboard: motorboard,
	}
	log.Info("Initialized firmware")
	return firmware, nil
}

type Firmware struct {
	config     *Config
	log        log.Logger
	navboard   *drivers.Navboard
	motorboard *drivers.Motorboard
}

// Loop causes the firmware to take control over the nav
func (f *Firmware) Loop() error {
	f.log.Info("Starting main loop")
	return nil
}

//navboard, err := drivers.NewNavboard(navdata.DefaultTTYPath)
//if err != nil {
//panic(err)
//}

//var data navdata.Data
//for {
//if err := driver.Decode(&data); err != nil {
//if err == navdata.ErrSync {
//log.Printf("%s\n", err)
//continue
//}
//panic(err)
//}

//log.Printf("%#v\n", data)
//}
