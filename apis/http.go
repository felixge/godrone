package apis

import (
	"fmt"
	"github.com/felixge/godrone/drivers"
	"io/ioutil"
	"net"
	"net/http"
	"regexp"
	"strconv"
)

func NewHttpAPI(port int, m *drivers.Motorboard) (*HttpAPI, error) {
	api := &HttpAPI{motorboard: m}
	mux := http.NewServeMux()
	mux.HandleFunc("/motors/", api.motors)

	addr := fmt.Sprintf(":%d", port)
	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}
	api.server = server

	listener, err := net.Listen("tcp4", addr)
	if err != nil {
		return nil, err
	}
	api.listener = listener

	return api, nil
}

type HttpAPI struct {
	listener   net.Listener
	server     *http.Server
	motorboard *drivers.Motorboard
}

func (h *HttpAPI) Serve() error {
	return h.server.Serve(h.listener)
}

var motorsRegexp = regexp.MustCompile("^/motors/([0-9]+)$")

func (h *HttpAPI) motors(w http.ResponseWriter, r *http.Request) {
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

		speed, err := strconv.ParseInt(string(body), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "invalid speed: %s\n", err)
			return
		}

		for i := 0; i < h.motorboard.MotorCount(); i++ {
			if !allMotors && i != motorId {
				continue
			}

			if err := h.motorboard.SetSpeed(i, int(speed)); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "could not set speed: %s\n", err)
				return
			}
		}
	}

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

		fmt.Fprintf(w, "%d\n", speed)
	}
}
