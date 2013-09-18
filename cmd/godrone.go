package main

import (
	"github.com/felixge/godrone"
)

func main() {
	// Initialize firmware
	firmware, err := godrone.NewFirmware(godrone.DefaultConfig)
	if err != nil {
		panic(err)
	}

	// Run the firmware
	if err := firmware.Run(); err != nil {
		panic(err)
	}
}
