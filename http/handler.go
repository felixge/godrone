// Package http implements the HTTP interface for GoDrone.
//
// It currently handles serving the HTML UI and related assets, as well as
// WebSocket clients.
package http

import (
	"code.google.com/p/go.net/websocket"
	"encoding/json"
	"fmt"
	"github.com/felixge/godrone/attitude"
	"github.com/felixge/godrone/control"
	"github.com/felixge/godrone/drivers/navboard"
	"github.com/felixge/godrone/http/fs"
	"github.com/felixge/godrone/log"
	"net/http"
	"sync"
)

// Config holds the arguments required to create a Handler.
type Config struct {
	Control *control.Control
	Log     log.Interface
	Version string
}

// Handler provides a http.Handler.
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

// NewHandler returns a new handler.
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

func (h *Handler) handleConfig(w http.ResponseWriter, r *http.Request) {
	config := map[string]interface{}{
		"version": h.config.Version,
	}
	data, err := json.Marshal(config)
	if err != nil {
		err = h.config.Log.Error("Could marshal JSON config. err=%s", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%s", err)
		return
	}
	w.Header().Set("Content-Type", "text/javascript; charset=UTF-8")
	fmt.Fprintf(w, "window.Config = %s;", data)
}

// ServeHTTP implements the http.Handler interface. It acts as a router
// dispatching websocket / asset requests.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" && r.URL.Path == "/ws" {
		h.websocketHandler.ServeHTTP(w, r)
		return
	}
	if r.Method == "GET" && r.URL.Path == "/js/config.js" {
		h.handleConfig(w, r)
		return
	}
	h.fileHandler.ServeHTTP(w, r)
}

// Update informs all connected clients about new navboard.Data and
// attitude.Data.
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
