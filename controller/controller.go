package controller

import (
	"github.com/felixge/godrone/controller/drivers"
	"github.com/felixge/pidctrl"
)

type Config struct{
	MotorboardTTY string
	NavboardTTY   string
}

func NewController(c Config, l log) {
	
}


type Controller struct {
	motorboard *drivers.Motorboard
	navboard   *drivers.Navboard
	leds *Leds
	xc *pidctrl.Pidctrl
}

func (c *Controller) Leds() *Leds {
	return c.leds
}


func (c *Controller) SetThrust(value float64) {
}




	leds := motorboard.Leds()
	for i := 0; i < 20; i++ {
		for l, _ := range leds {
			if i % 2 == 0 {
				leds[l] = drivers.LedOff
			} else {
				leds[l] = drivers.LedGreen
			}
		}
		motorboard.SetLeds(leds)
		time.Sleep(50 * time.Millisecond)
	}
