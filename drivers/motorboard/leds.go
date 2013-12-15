package motorboard

// LedColor holds a LED color constant.
type LedColor int

// The available colors for the LEDs
const (
	LedOff    LedColor = iota
	LedRed             = 1
	LedGreen           = 2
	LedOrange          = 3
)

// Leds is a helper method to return an array of LedColor's all set to the same
// value.
func Leds(c LedColor) [4]LedColor {
	return [4]LedColor{c, c, c, c}
}
