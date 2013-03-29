package attitude

import (
	"github.com/felixge/godrone/src/navdata"
	"sort"
)

// 1G = 9.81 ms/s^2
const oneG = float64(9.81)

type Attitude struct {
	navDriver *navdata.Driver
	navData   navdata.Data
	attData   Data
	// The initial bias for the accelerometers
	axBias    int
	ayBias    int
	azBias    int
	// The sensor output for 1G
	aOneG     int
}

func NewAttitude() (*Attitude, error) {
	driver, err := navdata.NewDriver(navdata.DefaultTTYPath)
	if err != nil {
		panic(err)
	}
	a := &Attitude{navDriver: driver}
	if err := a.init(); err != nil {
		return nil, err
	}
	return a, nil
}

func (a *Attitude) init() error {
	count := 10
	samples := struct{ X, Y, Z sort.IntSlice }{
		make(sort.IntSlice, count),
		make(sort.IntSlice, count),
		make(sort.IntSlice, count),
	}

	for i := 0; i < count; i++ {
		if err := a.navDriver.Decode(&a.navData); err != nil {
			return err
		}

		samples.X[i] = int(a.navData.Ax)
		samples.Y[i] = int(a.navData.Ay)
		samples.Z[i] = int(a.navData.Az)
	}

	sort.Sort(samples.X)
	sort.Sort(samples.Y)
	sort.Sort(samples.Z)

	// Get the median
	a.axBias = samples.X[count/2]
	a.ayBias = samples.Y[count/2]
	a.azBias = samples.Z[count/2]

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

	return &a.attData, nil
}

type Data struct {
	// Acceleration in m/s^2
	Ax float64
	Ay float64
	Az float64
}
