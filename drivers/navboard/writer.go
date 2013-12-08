package navboard

import (
	"fmt"
	"io"
)

type command byte

const (
	start command = 1 << iota
	stop
	resync
)

var commands = map[command]string{
	start:  "Start",
	stop:   "Stop",
	resync: "Resync",
}

func (c command) String() string {
	if s, ok := commands[c]; ok {
		return s
	} else {
		return fmt.Sprintf("unkown cmd: %d", c)
	}
}

func (c command) Bytes() []byte {
	return []byte{byte(c)}
}

func newWriter(w io.Writer) *writer {
	return &writer{w: w}
}

type writer struct {
	w io.Writer
}

func (w *writer) WriteCommand(cmd command) (err error) {
	_, err = w.w.Write(cmd.Bytes())
	return
}
