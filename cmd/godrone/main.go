package main

import (
	"log"
	"time"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
	log.Printf("Godrone started")
	navboard, err := OpenNavboard()
	if err != nil {
		log.Fatalf("Failed to open navboard: %s", err)
	}
	filter := NewComplementary()
	var (
		navdata  Navdata
		attitude PRY
		dt       = time.Second / 200
	)
	for {
		if err := navboard.Read(&navdata); err != nil {
			log.Printf("Failed to read navdata: %s", err)
			continue
		}
		filter.Update(&attitude, navdata, dt)
		log.Printf("%s", attitude)
	}
}
