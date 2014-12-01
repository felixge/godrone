package godrone

import "fmt"

// PRY holds pitch, roll, yaw information. The unit depends on the context
// where it is used.
type PRY struct {
	Pitch float64
	Roll  float64
	Yaw   float64
}

func (p PRY) String() string {
	return fmt.Sprintf("%7.2f P %7.2f R %7.2f Y", p.Pitch, p.Roll, p.Yaw)
}

func (p1 PRY) Equal(p2 PRY) bool {
	return p1.Pitch == p2.Pitch &&
		p1.Roll == p2.Roll &&
		p1.Yaw == p2.Yaw
}
