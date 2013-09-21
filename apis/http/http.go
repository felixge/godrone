package http

import (
	"code.google.com/p/go.net/websocket"
	"encoding/json"
	"fmt"
	"github.com/felixge/godrone/drivers"
	"github.com/felixge/godrone/log"
	"io/ioutil"
	"net"
	"net/http"
	"regexp"
	"strconv"
)

var files http.FileSystem

type responseWriter struct {
	http.ResponseWriter
	log log.Logger
}

func (r *responseWriter) writeJSON(val interface{}) {
	data, err := json.MarshalIndent(val, "", "  ")
	if err != nil {
		r.log.Err("writeJSON: %s for %#v", err, val)
		return
	}
	data = append(data, '\n')
	if _, err := r.Write(data); err != nil {
		r.log.Err("writeJSON: %s", err)
		return
	}
}

func NewHttpAPI(port int, m *drivers.Motorboard, n *drivers.Navboard, log log.Logger) (*HttpAPI, error) {
	api := &HttpAPI{motorboard: m, navboard: n, log: log}
	mux := http.NewServeMux()
	mux.HandleFunc("/motors/", api.motors)
	mux.HandleFunc("/navdata/", api.navdata)
	mux.Handle("/ws", websocket.Handler(api.websocket))
	mux.Handle("/", http.FileServer(files))
	api.mux = mux

	addr := fmt.Sprintf(":%d", port)
	server := &http.Server{
		Addr:    addr,
		Handler: http.HandlerFunc(api.serveHTTP),
	}
	api.server = server

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	api.listener = listener

	return api, nil
}

type HttpAPI struct {
	listener   net.Listener
	server     *http.Server
	mux        *http.ServeMux
	log        log.Logger
	motorboard *drivers.Motorboard
	navboard   *drivers.Navboard
}

func (h *HttpAPI) Serve() error {
	return h.server.Serve(h.listener)
}

func (h *HttpAPI) serveHTTP(w http.ResponseWriter, r *http.Request) {
	h.log.Debug("%s %s", r.Method, r.URL)
	h.mux.ServeHTTP(w, r)
}

var motorsRegexp = regexp.MustCompile("^/motors/([0-9]+)$")

func (h *HttpAPI) motors(hw http.ResponseWriter, r *http.Request) {
	w := &responseWriter{hw, h.log}

	if r.Method != "GET" && r.Method != "PUT" {
		w.Header().Set("Allow", "GET,PUT")
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "invalid method: %s\n", r.Method)
		return
	}

	var (
		allMotors = r.URL.Path == "/motors/"
		motorId   = -1
	)

	if !allMotors {
		m := motorsRegexp.FindStringSubmatch(r.URL.Path)
		if len(m) != 2 {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "unknown motor\n")
			return
		}

		mId, err := strconv.ParseInt(m[1], 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "invalid motor: %s\n", m[1])
			return
		}
		motorId = int(mId)
	}

	if r.Method == "PUT" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "could not read body: %s\n", err)
			return
		}

		speed, err := strconv.ParseFloat(string(body), 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "invalid speed: %s\n", err)
			return
		}

		for i := 0; i < h.motorboard.MotorCount(); i++ {
			if !allMotors && i != motorId {
				continue
			}

			if err := h.motorboard.SetSpeed(i, float64(speed)); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "could not set speed: %s\n", err)
				return
			}
		}
	}

	speeds := make(map[string]float64, h.motorboard.MotorCount())
	for i := 0; i < h.motorboard.MotorCount(); i++ {
		if !allMotors && i != motorId {
			continue
		}

		speed, err := h.motorboard.Speed(i)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "could not get speed: %s\n", err)
			return
		}

		speeds[strconv.FormatInt(int64(i), 10)] = speed
	}

	h.log.Debug("motor speeds: %#v", speeds)
	w.writeJSON(speeds)
}

func (h *HttpAPI) navdata(hw http.ResponseWriter, r *http.Request) {
	w := &responseWriter{hw, h.log}

	if r.Method != "GET" {
		w.Header().Set("Allow", "GET")
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "invalid method: %s\n", r.Method)
		return
	}

	navdata, err := h.navboard.Get()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "could not get navdata: %s\n", err)
		return
	}

	h.log.Debug("navdata: %#v", navdata)
	w.writeJSON(navdata)
}

func (h *HttpAPI) websocket(ws *websocket.Conn) {
	addr := ws.RemoteAddr()

	defer ws.Close()
	defer h.log.Info("Disconnected websocket client: %s", addr)

	h.log.Info("New websocket client: %s", addr)
	navdataCh, errCh := h.navboard.Subscribe()
	defer close(navdataCh)
	defer close(errCh)
	h.log.Info("Subscribed")

	for {
		var sendErr error
		select {
		case navdata := <-navdataCh:
			h.log.Debug("navdata: %#v", navdata)
			sendErr = websocket.JSON.Send(ws, navdata)
		case err := <-errCh:
			sendErr = websocket.JSON.Send(ws, err)
		}

		if sendErr != nil {
			h.log.Err("Error writing to websocket client: %s, %s", addr, sendErr)
			break
		}
	}
}
