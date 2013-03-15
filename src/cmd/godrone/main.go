package main

import (
	"fmt"
	"github.com/felixge/godrone/src/navdata"
)

func main() {
	driver, err := navdata.NewDriver(navdata.DefaultTTYPath)
	if err != nil {
		panic(err)
	}

	var data navdata.Data
	for {
		if err := driver.Decode(&data); err != nil {
			panic(err)
		}

		fmt.Printf("%#v\n", data)
	}
}
