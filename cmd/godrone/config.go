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
