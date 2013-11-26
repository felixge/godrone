package main

import (
	"github.com/felixge/godrone"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log := godrone.DefaultConfig.Log
	firmware := godrone.DefaultFirmware
	if err := firmware.Start(); err != nil {
		panic(err)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT)
	sig := <-sigCh
	log.Info("Received signal=%s, shutting down", sig)
	defer firmware.Stop()
}
