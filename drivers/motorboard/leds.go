package motorboard

type LedColor int

const (
	LedOff    LedColor = iota
	LedRed             = 1
	LedGreen           = 2
	LedOrange          = 3
)

func Leds(c LedColor) [4]LedColor {
	return [4]LedColor{c, c, c, c}
}
