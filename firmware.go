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

	log.Debug("Initializing navboard")
	navboard, err := drivers.NewNavboard(c.NavboardTTY)
	if err != nil {
		return nil, log.Emergency("Could not initialize navboard: %s", err)
	}

	firmware := &Firmware{
		config:   &c,
		log:      log,
		navboard: navboard,
	}
	return firmware, nil
}

type Firmware struct {
	config   *Config
	log      log.Logger
	navboard *drivers.Navboard
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
