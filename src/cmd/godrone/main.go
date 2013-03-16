package main

import (
	"github.com/felixge/godrone/src/navdata"
	"log"
)

func main() {
	log.SetFlags(log.Ltime | log.Lmicroseconds)

	driver, err := navdata.NewDriver(navdata.DefaultTTYPath)
	if err != nil {
		panic(err)
	}

	var data navdata.Data
	for {
		if err := driver.Decode(&data); err != nil {
			if err == navdata.ErrSync {
				log.Printf("%s\n", err)
				continue
			}
			panic(err)
		}

		log.Printf("%#v\n", data)
	}
}
