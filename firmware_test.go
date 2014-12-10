package godrone

import "testing"

func BenchmarkControlLoop(b *testing.B) {
	b.ReportAllocs()
	mb, err := OpenMotorboard("/dev/null")
	if err != nil {
		b.Fatal(err)
	}
	f, _ := NewCustomFirmware(&mockNavboard{}, mb)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f.Control()
	}
}

type mockNavboard struct{}

func (b *mockNavboard) Read() (data Navdata, err error) {
	return Navdata{
		Seq:      1,
		AccRoll:  2000,
		AccPitch: 2000,
		AccYaw:   8000,
	}, nil
}
