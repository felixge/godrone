package main

import (
	"fmt"
	"github.com/felixge/godrone"
	"os"
	"os/signal"
)

var (
	navboard = godrone.DefaultNavboard
	log      = godrone.DefaultLogger
)

func main() {
	if err := navboard.Calibrate(); err != nil {
		log.Fatal(err)
	}
	defer navboard.Close()

	go debug()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	<-sigCh
}

func debug() {
	for {
		data, err := navboard.NextData()
		if err != nil {
			continue
		}
		fmt.Printf("\r%s --> %s", data.Raw.ImuData(), data.Data)
	}
}
