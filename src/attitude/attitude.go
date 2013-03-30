package attitude

import (
	"github.com/felixge/godrone/src/navdata"
	"math"
	"sort"
)

// 1G = 9.81 ms/s^2
const oneG = float64(9.81)

type Attitude struct {
	// Input
	navDriver *navdata.Driver
	navData   navdata.Data

	// Output
	attData Data

	// Initial Accelerometer Bias
	axBias int
	ayBias int
	azBias int

	// Initial Gyroscope Bias
	gxBias int
	gyBias int
	gzBias int

	// The sensor output for 1G
	aOneG int
}

func NewAttitude(navDriver *navdata.Driver) (*Attitude, error) {
	a := &Attitude{navDriver: navDriver}
	if err := a.init(); err != nil {
		return nil, err
	}
	return a, nil
}

func (a *Attitude) init() error {
	count := 100
	samples := struct {
		Ax, Ay, Az sort.IntSlice
		Gx, Gy, Gz sort.IntSlice
	}{
		make(sort.IntSlice, count),
		make(sort.IntSlice, count),
		make(sort.IntSlice, count),
		make(sort.IntSlice, count),
		make(sort.IntSlice, count),
		make(sort.IntSlice, count),
	}

	for i := 0; i < count; i++ {
		if err := a.navDriver.Decode(&a.navData); err != nil {
			return err
		}

		samples.Ax[i] = int(a.navData.Ax)
		samples.Ay[i] = int(a.navData.Ay)
		samples.Az[i] = int(a.navData.Az)
		samples.Gx[i] = int(a.navData.Gx)
		samples.Gy[i] = int(a.navData.Gy)
		samples.Gz[i] = int(a.navData.Gz)
	}

	sort.Sort(samples.Ax)
	sort.Sort(samples.Ay)
	sort.Sort(samples.Az)
	sort.Sort(samples.Gx)
	sort.Sort(samples.Gy)
	sort.Sort(samples.Gz)

	// Get the median
	a.axBias = samples.Ax[count/2]
	a.ayBias = samples.Ay[count/2]
	a.azBias = samples.Az[count/2]
	a.gxBias = samples.Gx[count/2]
	a.gyBias = samples.Gy[count/2]
	a.gzBias = samples.Gz[count/2]

	// The drone is supposed to be on a flat surface when this code runs, and all
	// accelerometer seem to have the same output range. So the difference
	// between the Z axis and the X or Y axis should be the output for 1G. This
	// assumes that the sensor output is linear, which I hope to be the case : ).
	a.aOneG = a.azBias - a.axBias
	a.azBias = a.azBias - a.aOneG

	return nil
}

func (a *Attitude) Update() (*Data, error) {
	if err := a.navDriver.Decode(&a.navData); err != nil {
		return nil, err
	}

	a.attData.Ax = (float64(a.navData.Ax) - float64(a.axBias)) / float64(a.aOneG) * oneG
	a.attData.Ay = (float64(a.navData.Ay) - float64(a.ayBias)) / float64(a.aOneG) * oneG
	a.attData.Az = (float64(a.navData.Az) - float64(a.azBias)) / float64(a.aOneG) * oneG

	a.attData.Gx = int(a.navData.Gx) - a.gxBias
	a.attData.Gy = int(a.navData.Gy) - a.gyBias
	a.attData.Gz = int(a.navData.Gz) - a.gzBias

	//a.attData.Roll += a.attData.Gx
	//a.attData.Pitch += a.attData.Gy
	//a.attData.Yaw += a.attData.Gz

	a.attData.Pitch = math.Atan2(a.attData.Az, a.attData.Ax)*(180/math.Pi) - 90
	a.attData.Roll = math.Atan2(a.attData.Az, a.attData.Ay)*(180/math.Pi) - 90

	return &a.attData, nil
}

type Data struct {
	// Acceleration in m/s^2
	Ax float64
	Ay float64
	Az float64

	// Gyroscope data (unit unclear)
	Gx int
	Gy int
	Gz int

	// Results
	Roll  float64
	Pitch float64
	Yaw   float64
}
