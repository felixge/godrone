package navdata

import (
	"errors"
	"io"
	"log"
)

var ErrSyncFail = errors.New("navdata: could find start of stream")

const dataSize = 60

type Decoder struct {
	r      io.Reader
	synced bool
	buf []byte
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: r, buf: make([]byte, dataSize)}
}

func (d *Decoder) Decode(data *Data) error {
	if !d.synced {
		log.Printf("navdata: syncing ...")
		if err := d.sync(); err != nil {
			return err
		}
		d.synced = true
		log.Printf("navdata: synced")
	}

	if _, err := d.r.Read(d.buf); err != nil {
		// we're probably out of sync now
		d.synced = false
		return err
	}

	log.Printf("navdata: %#x", d.buf)

	return nil
}

// sync 
func (d *Decoder) sync() error {
	d.buf = d.buf[:0]

	_, err := d.r.Read(d.buf)
	if err != nil {
		return err
	}

	mark := byte(dataSize - 2)
	for i, b := range d.buf {
		if b == mark {
			if i + 1 < len(d.buf) && d.buf[i+1] == 0 {
				n := copy(d.buf, d.buf[i:])
				d.buf = d.buf[0:n]
				return nil
			}
		}
	}

	return ErrSyncFail
}

type Data struct {
}
