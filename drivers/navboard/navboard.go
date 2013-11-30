package navboard

import (
	"github.com/felixge/godrone/log"
	"os"
)

const (
	DefaultTTY = "/dev/ttyO1"
	Hz         = 200
)

func NewNavboard(tty string, log log.Interface) *Navboard {
	return &Navboard{
		tty: tty,
		log: log,
	}
}

type Navboard struct {
	reader *Reader
	writer *Writer
	tty    string
	file   *os.File
	log    log.Interface
}

func (n *Navboard) NextData() (data Data, err error) {
	defer func() {
		if err != nil {
			n.close()
		}
	}()

	if err = n.open(); err != nil {
		return
	}

	if data.Raw, err = n.reader.NextData(); err != nil {
		n.log.Error("Failed to read data. err=%s", err)
	} else {
		//n.log.Debug("Read data=%+v", data)
	}
	return
}

type calibration struct {
	Samples int64

	Ax0 int64
	Ay0 int64
	Az0 int64
	A1G int64 // Output for 1g

	Gx0 int64
	Gy0 int64
	Gz0 int64
}

func (n *Navboard) Calibrate() error {
	return nil
	//var (
		//samples = int64(40)
		//calib   calibration
	//)
	//n.log.Info("Calibrating. samples=%s", samples)
	//for calib.Samples < samples {
		//data, err := n.NextData()
		//if err != nil {
			//continue
		//}

		//calib.Samples++

		//calib.Ax0 += int64(data.Ax)
		//calib.Ay0 += int64(data.Ay)
		//calib.Az0 += int64(data.Az)

		//calib.Gx0 += int64(data.Gx)
		//calib.Gy0 += int64(data.Gy)
		//calib.Gz0 += int64(data.Gz)
	//}

	//calib.Ax0 /= calib.Samples
	//calib.Ay0 /= calib.Samples
	//calib.Az0 /= calib.Samples
	//calib.A1G = -(calib.Az0 - (calib.Ax0+calib.Ay0)/2)
	//calib.Az0 -= calib.A1G

	//calib.Gx0 /= calib.Samples
	//calib.Gy0 /= calib.Samples
	//calib.Gz0 /= calib.Samples

	//n.log.Info("Done calibrating. calibration=%+v", calib)
	//return nil
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

	n.log.Info("Opening tty=%s", n.tty)
	if n.file, err = os.OpenFile(n.tty, os.O_RDWR, 0); err != nil {
		return
	}
	n.writer = NewWriter(n.file)
	n.reader = NewReader(n.file)
	n.log.Debug("Writing start command")
	if err = n.writer.WriteCommand(Start); err != nil {
		return
	}
	n.log.Debug("Opened tty=%s", n.tty)
	return
}

func (n *Navboard) close() {
	n.log.Debug("Closing tty=%s", n.tty)
	if n.file != nil {
		n.file.Close()
	}
	n.file = nil
	n.reader = nil
	n.writer = nil
}
