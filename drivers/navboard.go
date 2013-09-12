package drivers

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"os"
)

type Navboard struct {
	*NavdataDecoder
	file *os.File
}

func NewNavboard(ttyPath string) (*Navboard, error) {
	file, err := os.OpenFile(ttyPath, os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}

	driver := &Navboard{
		file:    file,
		NavdataDecoder: NewNavdataDecoder(file),
	}

	if _, err := file.Write([]byte{3}); err != nil {
		return nil, err
	}

	return driver, nil
}

var ErrSync = errors.New("navdata: could not sync with stream")

const dataSize = 60

type NavdataDecoder struct {
	r      io.Reader
	offset int
	buf    []byte
}

func NewNavdataDecoder(r io.Reader) *NavdataDecoder {
	return &NavdataDecoder{r: r, buf: make([]byte, dataSize)}
}

// Raw returns a raw navdata payload. The returned buffer may be reused by
// later Raw() calls. This is a low-level method, direct usage is not
// recommended.
func (d *NavdataDecoder) Raw() ([]byte, error) {
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

// Decode reads and extracts the next navdata payload into *Navdata.
func (d *NavdataDecoder) Decode(data *Navdata) error {
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

// Navdata as found at https://github.com/RoboticaTUDelft/paparazzi/blob/minor1/sw/airborne/boards/ardrone/navdata.h
// Possibly not correct.
type Navdata struct {
	Seq uint16

	// Accelerometers
	Ax uint16
	Ay uint16
	Az uint16

	// Gyroscopes
	Gx int16
	Gy int16
	Gz int16

	// Everything below is unconfirmed, copied from other sources
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
