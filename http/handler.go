package http

import (
	"code.google.com/p/go.net/websocket"
	"github.com/codegangsta/martini"
	"github.com/felixge/godrone/attitude"
	"github.com/felixge/godrone/control"
	"github.com/felixge/godrone/http/fs"
	"github.com/felixge/godrone/log"
	"net/http"
)

type Setpoint struct {
	attitude.Data
	Throttle float64
}

func NewHandler(c *control.Control, log log.Interface) http.Handler {
	m := martini.Classic()

	m.Get("/ws", websocket.Handler(func(conn *websocket.Conn) {
		defer conn.Close()

		for {
			var s Setpoint
			if err := websocket.JSON.Receive(conn, &s); err != nil {
				log.Error("Could not receive setpoint. err=%s", err)
				return
			}

			c.Set(s.Data, s.Throttle)
		}
	}).ServeHTTP)

	m.NotFound(http.FileServer(fs.Fs).ServeHTTP)
	return m
}
