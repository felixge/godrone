package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

const (
	// packetSize is the size of a single navdata payload on the tty stream.
	packetSize = 0x3a
)

var packetHeader = []byte{packetSize, 0x00}

// OpenNavboard opens the navboard tty file at the given location and returns a
// Navboard struct on success.
func OpenNavboard(tty string) (*Navboard, error) {
	n := &Navboard{buf: &bytes.Buffer{}}
	return n, n.open(tty)
}

// Navboard provides access to the navboard. Must be used from a single
// goroutine.
type Navboard struct {
	file   *os.File
	reader *bufio.Reader
	buf    *bytes.Buffer
}

// open opens the navboard tty file.
func (b *Navboard) open(ttyPath string) error {
	file, err := os.OpenFile(ttyPath, os.O_RDWR, 0)
	if err != nil {
		return err
	}
	b.file = file
	b.reader = bufio.NewReader(file)
	return nil
}

// Read reads the next packet of navdata.
func (b *Navboard) Read() (data Navdata, err error) {
	// The loop below is used to sync with the packet stream received from the
	// navboard. This has to be done as the first bytes we read will often be in
	// the middle of a packet. A better approach would be to flush any buffered
	// data from the tty before the first read, but previous attempts of doing
	// this with ioctl and TCFLSH have not been successful. So for now we use the
	// known packetHeader as a marker to find in the stream. This can fail if
	// this value occurs in the middle of a packet and/or if Parrot adds
	// additional data to the packet in the future. Improvements are welcome!
	// Note: I may have had the parrot firmware running while experimenting with
	// TCFLSH, so this approach should be retried to make sure.
	i := 0
	for {
		var c byte
		c, err = b.reader.ReadByte()
		if err != nil {
			return
		}
		if c == packetHeader[i] {
			i++
			if i == len(packetHeader) {
				break
			}
		} else {
			i = 0
		}
	}

	if _, err = io.CopyN(b.buf, b.reader, packetSize); err != nil {
		return
	}
	sum := uint16(0)
	buf := b.buf.Bytes()
	for i := 0; i < len(buf)-2; i += 2 {
		sum += uint16(buf[i]) + (uint16(buf[i+1]) << 8)
	}
	if err = binary.Read(b.buf, binary.LittleEndian, &data); err != nil {
		return
	}
	if sum != data.Checksum {
		err = fmt.Errorf("Bad checksum. expected=%d got=%d", data.Checksum, sum)
		return
	}
	return
}

// Close closes the underlaying tty file.
func (d *Navboard) Close() error {
	err := d.file.Close()
	d.file = nil
	d.reader = nil
	return err
}

// Navdata holds the navboard data as read from the tty file. Based on
// https://github.com/paparazzi/paparazzi/blob/master/sw/airborne/boards/ardrone/navdata.h
// but with some adjustements for values that seem to be signed rather than
// unsigned.
type Navdata struct {
	Seq uint16

	// Accelerometers
	AccRoll  uint16
	AccPitch uint16
	AccYaw   uint16

	// Gyroscopes
	GyroRoll  int16
	GyroPitch int16
	GyroYaw   int16

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
