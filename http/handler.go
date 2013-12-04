package http

import (
	"code.google.com/p/go.net/websocket"
	"github.com/codegangsta/martini"
	"github.com/felixge/godrone"
	"github.com/felixge/godrone/http/fs"
	"github.com/felixge/godrone/log"
	"net/http"
)

type Control struct {
	Pitch    float64
	Roll     float64
	Yaw      float64
	Vertical float64
}

func NewHandler(fw *godrone.Firmware, log log.Interface) http.Handler {
	m := martini.Classic()

	m.Get("/ws", websocket.Handler(func(conn *websocket.Conn) {
		defer conn.Close()

		for {
			var c Control
			if err := websocket.JSON.Receive(conn, &c); err != nil {
				log.Error("Could not receive control data. err=%s", err)
				return
			}
			log.Debug("Control data: %#v", c)
		}
	}).ServeHTTP)

	m.NotFound(http.FileServer(fs.Fs).ServeHTTP)
	return m
}
