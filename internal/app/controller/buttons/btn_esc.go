package buttons

import (
	"log"

	"sschmc/internal/pkg/gpiobutton"
)

// BtnEscRisingHandler handles all cases of ESC button rising.
func (b *Buttons) BtnEscRisingHandler() gpiobutton.HandlerFunc {
	return func() {
		switch {
		case b.store.App.IsNone(), b.store.App.IsMenuMain():
			b.btnAllGreetings()
		case b.store.App.IsGreetings():
			b.btnEscGreetings()
		case b.store.App.IsMenuLevel():
			b.btnEscMenuLevel()
		case b.store.App.IsMessage():
			b.btnEscMessage()
		default:
			log.Println("*** ESCAPE pressed ***")
		}
	}
}

// btnEscGreetings clears rendered data and updates app-status in storage to none.
func (b *Buttons) btnEscGreetings() {
	b.store.App.SetNone()
	// update render according to new app-status
	b.render <- struct{}{}
}

// btnEscMenuMain updates app-status in storage to menu-main.
func (b *Buttons) btnEscMenuLevel() {
	b.store.App.SetMenuMain()
	// update render
	b.render <- struct{}{}
}

// btnEscMessage updates app-status in storage to menu-main.
func (b *Buttons) btnEscMessage() {
	b.store.App.SetMenuLevel()
	// update render according to new app-status
	b.render <- struct{}{}
}
