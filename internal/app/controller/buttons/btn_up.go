package buttons

import (
	"log"
	"sschmc/internal/pkg/gpiobutton"
)

// BtnUpRisingHandler handles all cases of UP button rising.
func (b *Buttons) BtnUpRisingHandler() gpiobutton.HandlerFunc {
	return func() {
		switch {
		case b.store.App.IsNone():
			b.btnAllGreetings()
		default:
			log.Println("*** UP pressed ***")
		}
	}
}
