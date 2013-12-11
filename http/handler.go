package http

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"github.com/codegangsta/martini"
	"github.com/felixge/godrone/attitude"
	"github.com/felixge/godrone/control"
	"github.com/felixge/godrone/http/fs"
	"github.com/felixge/godrone/log"
	"net/http"
	"time"
)

func NewHandler(c *control.Control, log log.Interface) http.Handler {
	m := martini.Classic()
	var setpoint attitude.Data

	go func() {
		for {
			fmt.Printf("setpoint: %s\n", setpoint)
			time.Sleep(time.Second)
		}
	}()

	m.Get("/ws", websocket.Handler(func(conn *websocket.Conn) {
		defer conn.Close()

		for {
			var cmd attitude.Data
			if err := websocket.JSON.Receive(conn, &cmd); err != nil {
				log.Error("Could not receive setpoint cmd. err=%s", err)
				return
			}

			setpoint.Roll = cmd.Roll
			setpoint.Pitch = cmd.Pitch
			setpoint.Yaw = cmd.Yaw
			setpoint.Altitude += cmd.Altitude * 0.01
			if setpoint.Altitude < 0 {
				setpoint.Altitude = 0
			} else if setpoint.Altitude > 3 {
				setpoint.Altitude = 3
			}

			c.Set(setpoint)
		}
	}).ServeHTTP)

	m.NotFound(http.FileServer(fs.Fs).ServeHTTP)
	return m
}
