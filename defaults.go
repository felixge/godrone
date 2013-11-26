package godrone

import (
	"github.com/felixge/godrone/drivers/navboard"
	"github.com/felixge/log"
	"time"
)

var (
	DefaultLogger              = log.DefaultLogger
	DefaultCalibrationDuration = time.Second
	DefaultNavboard            = navboard.NewNavboard(navboard.DefaultTTY, DefaultLogger)
	DefaultConfig              = Config{
		Log:                 DefaultLogger,
		Navboard:            DefaultNavboard,
		CalibrationDuration: DefaultCalibrationDuration,
	}
	DefaultFirmware = NewFirmware(DefaultConfig)
)
