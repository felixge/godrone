package main

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"github.com/felixge/godrone/src/attitude"
	"github.com/felixge/godrone/src/motorboard"
	"github.com/felixge/godrone/src/navdata"
	"log"
	"net/http"
	"strconv"
	"sync"
)

var clients []*websocket.Conn
var clientsLock sync.Mutex
var motorsLock sync.Mutex
var motors *motorboard.Driver
var speeds [4]int

func main() {
	log.SetFlags(log.Ltime | log.Lmicroseconds)

	go serveHttp()

	log.Printf("Initializing motorboard ...")
	motors, err := motorboard.NewDriver(motorboard.DefaultTTYPath)
	if err != nil {
		panic(err)
	}
	motors.SetLeds(motorboard.LedRed)

	navDriver, err := navdata.NewDriver(navdata.DefaultTTYPath)
	if err != nil {
		panic(err)
	}

	log.Printf("Initializing attitude ...")
	att, err := attitude.NewAttitude(navDriver)
	if err != nil {
		panic(err)
	}

	log.Printf("Starting main loop ...")
	motors.SetLeds(motorboard.LedGreen)


	for {
		data, err := att.Update()
		if err != nil {
			panic(err)
		}

		_ = data
		motors.Speeds[0] = speeds[0]
		motors.Speeds[1] = speeds[1]
		motors.Speeds[2] = speeds[2]
		motors.Speeds[3] = speeds[3]

		//rollError := data.Roll / 90
		//if rollError >= 0 {
			//motors.Speeds[0] += int(rollError * float64(2048))
			//motors.Speeds[3] += int(rollError * float64(2048))
		//} else if rollError < 0 {
			//motors.Speeds[1] += int(-rollError * float64(2048))
			//motors.Speeds[2] += int(-rollError * float64(2048))
		//}

		//pitchError := data.Pitch / 90
		//if pitchError >= 0 {
			//motors.Speeds[0] += int(pitchError * float64(2048))
			//motors.Speeds[1] += int(pitchError * float64(2048))
		//} else if pitchError < 0 {
			//motors.Speeds[2] += int(-pitchError * float64(2048))
			//motors.Speeds[3] += int(-pitchError * float64(2048))
		//}

		//if motors.Speeds[0] > 511 {
		//motors.Speeds[0] = 511
		//}
		//if motors.Speeds[1] > 511 {
		//motors.Speeds[1] = 511
		//}
		//if motors.Speeds[2] > 511 {
		//motors.Speeds[2] = 511
		//}
		//if motors.Speeds[3] > 511 {
		//motors.Speeds[3] = 511
		//}

		//fmt.Printf("%f: %d, %d\n", data.Roll, motors.Speeds[0], motors.Speeds[1])

		//_ = data
		//motorsLock.Lock()
		//if err := motors.UpdateSpeeds(); err != nil {
			//panic(err)
		//}
		//motorsLock.Unlock()
		//if err := motors.UpdateLeds(); err != nil {
		//panic(err)
		//}
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
		websocket.Message.Receive(ws, &d)
		val, err := strconv.ParseInt(d, 10, 32)
		if err != nil {
			panic(err)
		}

		for i := 0; i < len(speeds); i++ {
			speeds[i] = int(val)
		}
		fmt.Printf("received: %#v\n", d)
	}
}
