package navboard

import (
	"github.com/felixge/godrone/log"
	"os"
	"syscall"
)

var (
	DefaultTTY = "/dev/ttyO1"
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
	if data, err = n.reader.NextData(); err != nil {
		n.log.Error("Failed to read data. err=%s", err)
	} else {
		//n.log.Debug("Read data=%+v", data)
	}
	return
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
	n.log.Debug("Writing stop command")
	if err = n.writer.WriteCommand(Stop); err != nil {
		return
	}
	n.log.Debug("Setting O_NONBLOCK")
	if _, err = n.fcntl(syscall.F_SETFL, syscall.O_NONBLOCK); err != nil {
		return
	}
	n.log.Debug("Draining tty")
	if err = n.reader.Drain(); err != nil {
		return
	}
	n.log.Debug("Setting O_ASYNC")
	if _, err = n.fcntl(syscall.F_SETFL, syscall.O_ASYNC); err != nil {
		return
	}
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

func (n *Navboard) fcntl(cmd int, arg int) (val int, err error) {
	return fcntl(n.file.Fd(), cmd, arg)
}

func fcntl(fd uintptr, cmd int, arg int) (val int, err error) {
	v, _, e := syscall.Syscall(syscall.SYS_FCNTL, fd, uintptr(cmd), uintptr(arg))
	val = int(v)
	if e != 0 {
		err = e
	}
	return val, err
}
