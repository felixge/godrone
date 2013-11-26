package navboard

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"
	"io/ioutil"
	"os"
	"syscall"
)

func NewReader(reader io.Reader) *Reader {
	return &Reader{
		bufReader: bufio.NewReader(reader),
	}
}

type Reader struct {
	bufReader *bufio.Reader
}

func (r *Reader) Drain() (err error) {
	_, err = io.Copy(ioutil.Discard, r.bufReader)
	if perr, ok := err.(*os.PathError); ok {
		if perr.Err == syscall.EAGAIN {
			err = nil
		}
	}
	return
}

func (r *Reader) NextData() (data Data, err error) {
	var payloadLength uint16
	if err = binary.Read(r.bufReader, binary.LittleEndian, &payloadLength); err != nil {
		return
	}
	payload := bytes.NewBuffer(make([]byte, 0, payloadLength))
	if _, err = io.CopyN(payload, r.bufReader, int64(payloadLength)); err != nil {
		return
	}
	err = binary.Read(payload, binary.LittleEndian, &data)
	return
}
