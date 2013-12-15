package navboard

import (
	"io"
)

type command byte

// Navboard control commands. Found in Parrot SDK:
// https://projects.ardrone.org/embedded/ardrone-api/d1/d4d/ardrone__common__config_8h.html#9fc16d3fe3d71ffea3841efded3132ee
// https://projects.ardrone.org/embedded/ardrone-api/d3/d19/ardrone__common__config_8h-source.html
const (
	// command to start acquisition with ADC
	startAcq command = 1
	// command to stop acquisition with ADC
	stopAcq = 2
	// command to resync acquisition with ADC
	resync = 3
	// command to ADC send a test frame (123456)
	test = 4
	// command to ADC send his number : version (MSB) subversion (LSB)
	version = 5
	// set the ultrasound at 22,22Hz
	selectUltrasound22hz = 7
	// set the ultrasound at 25Hz
	selectUltrasound25hz = 8
	// command to ADC to send the calibration
	sendCalibre = 13
	// command to ADC to receved a new calibration
	recevedCalibre = 14
	// get the hard version of the navboard
	getHardVersion = 15
	// enabled the separation of sources ultrasound
	activeSeparation = 16
	// disables the ultrasound source separation
	stopSeparation = 17
	// command to ADC to receved the prod data
	recevedProd = 18
	// command to ADC to send the prod data
	sendProd = 19
	// command to ADC to send PWM ultrasond in continue
	activeEtalonage = 20
	// command to ADC to stop send PWM ultrasond in continue
	activeUltrason     = 21
	activeTestUltrason = 22
)

func (c command) Bytes() []byte {
	return []byte{byte(c)}
}

func newWriter(w io.Writer) *writer {
	return &writer{w: w}
}

type writer struct {
	w io.Writer
}

// WriteCommand writes the given cmd.
func (w *writer) WriteCommand(cmd command) (err error) {
	_, err = w.w.Write(cmd.Bytes())
	return
}
