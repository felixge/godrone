package main

import (
	"flag"
	"log"
)
import "net/http"

import (
	"path/filepath"
	"runtime"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
	var addr = flag.String("addr", ":8080", "The addr to listen on.")
	flag.Parse()
	dir := flag.Arg(0)
	if dir == "" {
		_, srcFile, _, _ := runtime.Caller(0)
		dir = filepath.Join(filepath.Dir(srcFile), "public")
	}
	log.Printf("Listening on: %s", *addr)
	if err := http.ListenAndServe(*addr, http.FileServer(http.Dir(dir))); err != nil {
		log.Fatalf("%s", err)
	}
}
