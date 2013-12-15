package navboard

import (
	"fmt"
	"github.com/felixge/godrone/imu"
	"github.com/felixge/godrone/log"
	"math"
	"os"
)

const (
	DefaultTTY = "/dev/ttyO1"
)

// gyroGains are the measured gains for converting the gyroscope output into
// deg/sec.
//
// from: /data/config.ini on drone
// gyros_gains                    = { 1.0569232e-03 -1.0664322e-03 -1.0732636e-03 }
var gyroGains = [3]float64{16, -16, -16}

func NewNavboard(tty string, log log.Interface) *Navboard {
	return &Navboard{
		tty: tty,
		log: log,
	}
}

type Navboard struct {
	reader      *reader
	writer      *writer
	tty         string
	file        *os.File
	log         log.Interface
	calibration calibration
}

// @TODO Turn into ReadData, taking a pointer
func (n *Navboard) NextData() (data Data, err error) {
	defer func() {
		if err != nil {
			n.Close()
		}
	}()

	if err = n.open(); err != nil {
		return
	}

	if data.Raw, err = n.reader.NextData(); err != nil {
		n.log.Error("Failed to read data. err=%s", err)
		return
	}

	data.Data = data.Raw.ImuData().Sub(n.calibration.Offsets).Div(n.calibration.Gains)
	// taken from ardrone project, could not find out the sensor model / verify
	// the details on this yet. It also seems like the minimum value is always
	// 30cm.
	// see https://github.com/ardrone/ardrone/blob/master/ardrone/navboard/navboard.c#L107
	data.Data.UsAltitude = float64(data.Raw.Ultrasound&0x7fff) * 0.000340

	return
}

type calibration struct {
	Offsets imu.Data
	Gains   imu.Data
}

func (n *Navboard) Calibrate() error {
	var (
		samples                      = make([]imu.Floats, 0, 40)
		retries                      = 100
		sums, sqrSums, means, stdevs imu.Floats
		names                        = []string{"Ax", "Ay", "Az", "Gx", "Gy", "Gz"}
	)

	for len(samples) < cap(samples) {
		data, err := n.NextData()
		if err != nil {
			if retries <= 0 {
				return err
			}
			retries--
			continue
		}

		values := data.Raw.ImuData().Floats()
		samples = append(samples, values)

		for i, val := range values {
			sums[i] += val
		}
	}

	for i := 0; i < len(means); i++ {
		means[i] = sums[i] / float64(len(samples))
	}

	for _, values := range samples {
		for i, v := range values {
			sqrSums[i] += (v - means[i]) * (v - means[i])
		}
	}

	for i := 0; i < len(stdevs); i++ {
		stdevs[i] = math.Sqrt(sqrSums[i] / float64(len(samples)))
	}

	n.log.Debug("calibration means: %#v", means)
	n.log.Debug("calibration stdevs: %#v", stdevs)

	for i, stdev := range stdevs {
		// stddev is usually around 1 for all sensors. 10 is an empirical value
		// that indicates there is too much sensor noise (drone is moving or
		// sensors are going crazy).
		if stdev > 20 {
			return fmt.Errorf("Standard deviation too high: std=%.2f sensor=%s", stdev, names[i])
		}
	}

	var o, g imu.Floats
	o = means

	// The drone should sitting on a level floor at this point, so we assume that
	// Az is measuring 1G, and we can calculate the gain for that by substracting
	// the avg between Ax and Ay from Az. This assumes that the gain is identical
	// for all sensors, but that's much more convenient than trying to measure
	// the gain of each sensor individually.
	ag := o[imu.Az] - (o[imu.Ax]+o[imu.Ay])/2
	g[imu.Ax], g[imu.Ay], g[imu.Az] = ag, ag, -ag
	o[imu.Az] -= ag

	// @TODO Figure out gains for gyroscopes
	g[imu.Gx], g[imu.Gy], g[imu.Gz] = gyroGains[0], gyroGains[1], gyroGains[2]

	n.calibration.Offsets = imu.NewData(o)
	n.calibration.Gains = imu.NewData(g)

	n.log.Debug("calibration offsets: %+v", n.calibration.Offsets)
	n.log.Debug("calibration gains: %+v", n.calibration.Gains)

	return nil
}

func (n *Navboard) open() (err error) {
	defer func() {
		if err != nil {
			n.log.Error("Could not open tty. tty=%s err=%#v", n.tty, err)
		}
	}()

	if n.file != nil {
		return
	}

	n.log.Debug("Opening tty=%s", n.tty)
	if n.file, err = os.OpenFile(n.tty, os.O_RDWR, 0); err != nil {
		return
	}
	n.writer = newWriter(n.file)
	n.reader = newReader(n.file)
	n.log.Debug("Writing start command")
	if err = n.writer.WriteCommand(start); err != nil {
		return
	}
	n.log.Debug("Opened tty=%s", n.tty)
	return
}

func (n *Navboard) Close() (err error) {
	n.log.Debug("Closing tty=%s", n.tty)
	if n.file != nil {
		err = n.file.Close()
	}
	n.file = nil
	n.reader = nil
	n.writer = nil
	return
}
