package main

import (
	"github.com/felixge/godrone"
	"github.com/felixge/godrone/http"
	gohttp "net/http"
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
	defer firmware.Stop()

	h := http.NewHandler(firmware, log)
	go gohttp.ListenAndServe(":80", h)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT)
	sig := <-sigCh
	log.Info("Received signal=%s, shutting down", sig)
}
