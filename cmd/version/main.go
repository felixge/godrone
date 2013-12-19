// Command version is a utility that prints the current godrone version and
// exits.
package main

import (
	"fmt"
	"github.com/felixge/godrone"
)

func main() {
	fmt.Println(godrone.Version)
}
