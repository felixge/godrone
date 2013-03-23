package main

import (
	"code.google.com/p/go.net/websocket"
	"io"
	"log"
	"net/http"
)

func main() {
	log.SetFlags(log.Ltime | log.Lmicroseconds)

	http.Handle("/ws", websocket.Handler(handleWs))

	addr := ":80"
	log.Printf("serving clients at %s", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}

func handleWs(ws *websocket.Conn) {
	io.Copy(ws, ws)
}
