package navboard

import (
	"github.com/felixge/godrone/imu"
)

// Data holds the navboard data after adjusting for sensitivity / bias.
type Data struct {
	imu.Data
	Raw RawData
}

// RawNavdata holds the navboard data as read from the tty file.
// Based on https://github.com/RoboticaTUDelft/paparazzi/blob/minor1/sw/airborne/boards/ardrone/navdata.h
// but adjusted a little (some values seem to be signed rather than unsigned)
type RawData struct {
	Seq uint16

	// Accelerometers
	Ax uint16
	Ay uint16
	Az uint16

	// Gyroscopes
	Gx int16
	Gy int16
	Gz int16

	TemperatureAcc  uint16
	TemperatureGyro uint16

	Ultrasound int16

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

// ImuData extracts the unadjusted imu.Data contained in the raw data.
func (r RawData) ImuData() imu.Data {
	return imu.Data{
		Ax: float64(r.Ax),
		Ay: float64(r.Ay),
		Az: float64(r.Az),
		Gx: float64(r.Gx),
		Gy: float64(r.Gy),
		Gz: float64(r.Gz),
	}
}
