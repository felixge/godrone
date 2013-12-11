package imu

import (
	"fmt"
	"strings"
)

type Sensor int

func (s Sensor) String() string {
	return names[s]
}

const (
	Ax Sensor = iota
	Ay
	Az
	Gx
	Gy
	Gz
	Sensors
)

var names = map[Sensor]string{
	Ax: "Ax",
	Ay: "Ay",
	Az: "Az",
	Gx: "Gx",
	Gy: "Gy",
	Gz: "Gz",
}

type Floats [Sensors]float64

type Data struct {
	Ax, Ay, Az float64
	Gx, Gy, Gz float64
	UsAltitude  float64
}

func NewData(f Floats) (d Data) {
	d.Ax, d.Ay, d.Az = f[Ax], f[Ay], f[Az]
	d.Gx, d.Gy, d.Gz = f[Gx], f[Gy], f[Gz]
	return
}

func (d Data) Floats() (f Floats) {
	f[Ax], f[Ay], f[Az] = d.Ax, d.Ay, d.Az
	f[Gx], f[Gy], f[Gz] = d.Gx, d.Gy, d.Gz
	return
}

func (d Data) Sub(v Data) (r Data) {
	df := d.Floats()
	vf := v.Floats()
	rf := Floats{}

	for i := 0; i < int(Sensors); i++ {
		rf[i] = df[i] - vf[i]
	}
	return NewData(rf)
}

func (d Data) Div(v Data) (r Data) {
	df := d.Floats()
	vf := v.Floats()
	rf := Floats{}

	for i := 0; i < int(Sensors); i++ {
		rf[i] = df[i] / vf[i]
	}
	return NewData(rf)
}

func (d Data) String() string {
	var (
		f      = d.Floats()
		format = make([]string, len(f))
		args   = make([]interface{}, len(f))
	)
	for i := 0; i < len(f); i++ {
		format[i] = "%+ 8.2f " + Sensor(i).String()
		args[i] = interface{}(f[i])
	}
	return fmt.Sprintf(strings.Join(format, " "), args...)
}
