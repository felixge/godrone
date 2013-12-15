package main

import (
	"github.com/BurntSushi/toml"
	"os"
	"path/filepath"
)

// Config holds the user configuration for the GoDrone firmware.
type Config struct {
	NavboardTTY   string
	MotorboardTTY string
	RollPID       []float64
	PitchPID      []float64
	YawPID        []float64
	HttpAddr      string
}

var (
	defaultRollPitchPID = []float64{0.04, 0, 0.002}
)

// DefaultConfig provides sensible defaults in absence of a config file.
var DefaultConfig = Config{
	NavboardTTY:   "/dev/ttyO1",
	MotorboardTTY: "/dev/ttyO0",
	RollPID:       defaultRollPitchPID,
	PitchPID:      defaultRollPitchPID,
	YawPID:        []float64{0.04, 0, 0}, // disabled, needs magnotometer to work well
	HttpAddr:      ":80",
}

// LoadConfig loads the configuration from a toml file. Other file formats may
// be supported in the future as well.
func LoadConfig(file string, config *Config) error {
	if string(file[0]) != "/" {
		wd, err := os.Getwd()
		if err != nil {
			return err
		}
		file = filepath.Join(wd, file)
	}
	_, err := toml.DecodeFile(file, &config)
	return err
}
