package navdata

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

var ErrSync = errors.New("navdata: could not sync with stream")

const dataSize = 60

type Decoder struct {
	r      io.Reader
	offset int
	buf    []byte
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: r, buf: make([]byte, dataSize)}
}

// Raw returns a raw navdata payload. The returned buffer may be reused by
// later Raw() calls. This is a low-level method, direct usage is not
// recommended.
func (d *Decoder) Raw() ([]byte, error) {
	offset := 0
	for {
		n, err := d.r.Read(d.buf[offset:])
		if err != nil {
			return nil, err
		}
		offset += n
		if offset == len(d.buf) {
			break
		}
	}

	mark := byte(dataSize - 2)
	for i, b := range d.buf {
		if b == mark {
			if i+1 < len(d.buf) && d.buf[i+1] == 0 {
				offset = copy(d.buf, d.buf[i:])
				break
			}
		}

		if i+1 >= len(d.buf) {
			return nil, ErrSync
		}
	}

	// @TODO Remove code duplication with other loop.
	for {
		n, err := d.r.Read(d.buf[offset:])
		if err != nil {
			return nil, err
		}
		offset += n
		if offset == len(d.buf) {
			break
		}
	}

	return d.buf, nil
}

// Decode reads and extracts the next navdata payload into *Data.
func (d *Decoder) Decode(data *Data) error {
	raw, err := d.Raw()
	if err != nil {
		return err
	}

	raw = raw[2:]
	if err := binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, data); err != nil {
		return err
	}

	// @TODO verify checksum (not sure if data.Checksum is correct / how to
	// calculate it)

	return nil
}


// Data as found at https://github.com/RoboticaTUDelft/paparazzi/blob/minor1/sw/airborne/boards/ardrone/navdata.h
// Possibly not correct.
type Data struct {
	Seq uint16

	Ax uint16
	Ay uint16
	Az uint16

	Vx              uint16
	Vy              uint16
	Vz              uint16
	TemperatureAcc  uint16
	TemperatureGyro uint16

	Ultrasound uint16

	UsDebutEcho       uint16
	UsFinEcho         uint16
	UsAssociationEcho uint16
	UsDistanceEcho    uint16

	UsCurveTime  uint16
	UsCurveValue uint16
	UsCurveRef   uint16

	NbEcho uint16

	SumEcho  uint32
	Gradient int16

	FlagEchoIni uint16

	Pressure            int32
	TemperaturePressure int16

	Mx int16
	My int16
	Mz int16

	Checksum uint16
}
