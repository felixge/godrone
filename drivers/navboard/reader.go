package navboard

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
)

func NewReader(reader io.Reader) *Reader {
	return &Reader{
		bufReader: bufio.NewReader(reader),
	}
}

type Reader struct {
	bufReader *bufio.Reader
}

func (r *Reader) NextData() (raw RawData, err error) {
	var (
		length   uint16
		expected = binary.Size(raw)
		skipped  int
	)

	// Look for the beginning of a navdata packet as indicated by the payload
	// size. This is hacky and will break if parrot increases the payload size,
	// but unfortunately I've been unable with a better sync mechanism, including
	// a very fancy attempt to stop the aquisition, drain the tty buffer in
	// non-blocking mode, and then restart the aquisition. Better ideas are
	// welcome!
	for {
		if err = binary.Read(r.bufReader, binary.LittleEndian, &length); err != nil {
			return
		}
		if int(length) == expected {
			break
		}
		if skipped > expected {
			err = fmt.Errorf("Failed to find payload.")
			return
		}
		skipped += binary.Size(length)
	}
	err = binary.Read(r.bufReader, binary.LittleEndian, &raw)
	return
}
