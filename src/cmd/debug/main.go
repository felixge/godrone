package main

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"github.com/felixge/godrone/src/attitude"
	"github.com/felixge/godrone/src/navdata"
	"log"
	"net/http"
	"sync"
)

var clients []*websocket.Conn
var clientsLock sync.Mutex

func main() {
	log.SetFlags(log.Ltime | log.Lmicroseconds)

	go serveHttp()

	log.Printf("Initializing sensors ...")
	driver, err := navdata.NewDriver(navdata.DefaultTTYPath)
	if err != nil {
		panic(err)
	}

	att, err := attitude.NewAttitude(driver)
	if err != nil {
		panic(err)
	}

	log.Printf("Starting main loop ...")

	i := 0
	for {
		data, err := att.Update()
		if err != nil {
			panic(err)
		}

		i++
		//fmt.Printf("%f | %f\n", data.Roll, data.Pitch)
		if i % 10 == 0 {
			fmt.Printf("0:%f\n", data.Roll)
			fmt.Printf("1:%f\n", data.Pitch)
		}
		//fmt.Printf("%f | %f | %f\n", data.Ax, data.Ay, data.Az)
	}
}

func serveHttp() {
	http.Handle("/ws", websocket.Handler(handleWs))
	addr := ":80"
	log.Printf("serving clients at %s", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}

func handleWs(ws *websocket.Conn) {
	log.Printf("New client: %s", ws.RemoteAddr().String())
	clientsLock.Lock()
	clients = append(clients, ws)
	clientsLock.Unlock()

	var d string
	for {
		websocket.Message.Receive(ws, &d);
	}
}
