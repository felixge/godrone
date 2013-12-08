package motorboard

import (
	"io"
)

const (
	// pwmMax is the max speed the motors can be set to
	pwmMax = 511

	// Commands
	setSpeeds = 0x20
	setLeds   = 0x60
)

func newWriter(w io.Writer) *writer {
	return &writer{w: w}
}

type writer struct {
	w io.Writer
}

// WriteSpeeds writes the given [0,1] ranged speeds
func (w *writer) WriteSpeeds(speeds [4]float64) (err error) {
	var pwms [4]uint16
	for i, speed := range speeds {
		pwms[i] = uint16(speed * pwmMax)
	}

	// see: https://github.com/ardrone/ardrone/blob/master/ardrone/motorboard/motorboard.c
	cmd := []byte{
		byte(setSpeeds | ((pwms[0] & 0x1ff) >> 4)),
		byte(((pwms[0] & 0x1ff) << 4) | ((pwms[1] & 0x1ff) >> 5)),
		byte(((pwms[1] & 0x1ff) << 3) | ((pwms[2] & 0x1ff) >> 6)),
		byte(((pwms[2] & 0x1ff) << 2) | ((pwms[3] & 0x1ff) >> 7)),
		byte(((pwms[3] & 0x1ff) << 1)),
	}
	_, err = w.w.Write(cmd)
	return
}

// cmd = 011rrrrx xxxggggx (used to be 011grgrg rgrxxxxx in AR Drone 1.0)
// see: https://github.com/ardrone/ardrone/blob/master/ardrone/motorboard/motorboard.c#L243
func (w *writer) WriteLeds(leds [4]LedColor) (err error) {
	cmd := make([]byte, 2)
	cmd[0] = setLeds

	for i, color := range leds {
		if color == LedRed || color == LedOrange {
			cmd[0] = cmd[0] | (1 << (byte(i) + 1))
		}

		if color == LedGreen || color == LedOrange {
			cmd[1] = cmd[1] | (1 << (byte(i) + 1))
		}
	}
	_, err = w.w.Write(cmd)
	return
}
