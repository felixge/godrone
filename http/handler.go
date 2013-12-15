// Package http implements the HTTP interface for GoDrone.
//
// It currently handles serving the HTML UI and related assets, as well as
// WebSocket clients.
package http

import (
	"code.google.com/p/go.net/websocket"
	"github.com/felixge/godrone/attitude"
	"github.com/felixge/godrone/control"
	"github.com/felixge/godrone/drivers/navboard"
	"github.com/felixge/godrone/http/fs"
	"github.com/felixge/godrone/log"
	"net/http"
	"sync"
)

type Config struct {
	Control *control.Control
	Log     log.Interface
}

type Handler struct {
	lock             sync.Mutex
	config           Config
	websocketHandler http.Handler
	fileHandler      http.Handler
	listeners        []chan update
}

type update struct {
	NavData      navboard.Data
	AttitudeData attitude.Attitude
}

type setpoint struct {
	attitude.Attitude
	Throttle float64
}

func NewHandler(c Config) *Handler {
	h := &Handler{
		config:      c,
		fileHandler: http.FileServer(fs.Fs),
	}
	h.websocketHandler = websocket.Handler(h.handleWebsocket)
	return h
}

func (h *Handler) handleWebsocket(conn *websocket.Conn) {
	var (
		log      = h.config.Log
		ip       = conn.Request().RemoteAddr
		setCh    = make(chan setpoint, 1)
		setErrCh = make(chan error, 1)
	)

	defer conn.Close()

	log.Info("New WebSocket connection. ip=%s", ip)
	defer log.Info("Closed WebSocket connection. ip=%s", ip)

	updateCh := h.sub()
	defer h.unsub(updateCh)

	go func() {
		for {
			var s setpoint
			if err := websocket.JSON.Receive(conn, &s); err != nil {
				setErrCh <- err
				return
			}
			setCh <- s
		}
	}()

	for {
		select {
		case u := <-updateCh:
			if err := websocket.JSON.Send(conn, u); err != nil {
				log.Warn("WebSocket error. err=%s ip=%s", err, ip)
				return
			}
		case s := <-setCh:
			h.config.Control.Set(s.Attitude, s.Throttle)
		case err := <-setErrCh:
			log.Warn("WebSocket error. err=%s ip=%s", err, ip)
			return
		}
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" && r.URL.Path == "/ws" {
		h.websocketHandler.ServeHTTP(w, r)
		return
	}
	h.fileHandler.ServeHTTP(w, r)
}

func (h *Handler) Update(n navboard.Data, a attitude.Attitude) {
	h.pub(update{n, a})
}

func (h *Handler) pub(u update) {
	h.lock.Lock()
	defer h.lock.Unlock()

	for _, ch := range h.listeners {
		select {
		case ch <- u:
		default:
		}
	}
}

func (h *Handler) sub() chan update {
	ch := make(chan update)
	h.lock.Lock()
	defer h.lock.Unlock()
	h.listeners = append(h.listeners, ch)
	return ch
}

func (h *Handler) unsub(ch chan update) {
	h.lock.Lock()
	defer h.lock.Unlock()
	for i, chEntry := range h.listeners {
		if ch == chEntry {
			h.listeners = append(h.listeners[:i], h.listeners[i+1:]...)
			return
		}
	}
	panic("failed to unsub")
}
