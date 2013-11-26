package navboard

import (
	"fmt"
	"io"
)

type Command byte

const (
	Start Command = 1 << iota
	Stop
	Resync
)

var commands = map[Command]string{
	Start:  "Start",
	Stop:   "Stop",
	Resync: "Resync",
}

func (c Command) String() string {
	if s, ok := commands[c]; ok {
		return s
	} else {
		return fmt.Sprintf("unkown cmd: %d", c)
	}
}

func (c Command) Bytes() []byte {
	return []byte{byte(c)}
}

func NewWriter(writer io.Writer) *Writer {
	return &Writer{writer: writer}
}

type Writer struct {
	writer io.Writer
}

func (w *Writer) WriteCommand(cmd Command) (err error) {
	_, err = w.writer.Write(cmd.Bytes())
	return
}
