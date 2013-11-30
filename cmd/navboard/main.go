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
	//if err := navboard.Calibrate(); err != nil {
	//panic(err)
	//}
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
			log.Error("Could not read navdata. err=%s", err)
			continue
		}
		fmt.Printf(
			"\rAX: %05d AY: %05d AZ: %05d GX: %05d GY: %05d GZ: %05d",
			data.Raw.Ax,
			data.Raw.Ay,
			data.Raw.Az,
			data.Raw.Gx,
			data.Raw.Gy,
			data.Raw.Gz,
		)
	}
}
