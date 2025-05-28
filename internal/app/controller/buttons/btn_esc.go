package buttons

import (
	"log"

	"sschmc/internal/pkg/gpiobutton"
)

// BtnEscRisingHandler handles all cases of ESC button rising.
func (b *Buttons) BtnEscRisingHandler() gpiobutton.HandlerFunc {
	return func() {
		switch {
		case b.store.App.IsNone(): //, b.store.App.IsMenuMain():
			log.Println("*** ESCAPE none/menu ***")
			b.btnAllGreetings()
		case b.store.App.IsGreetings():
			log.Println("*** ESCAPE greetings ***")
			b.btnEscGreetings()
		case b.store.App.IsMenuMain():
			log.Println("*** ESCAPE UP menu ***")
			b.btnEscMenuMain()
		case b.store.App.IsMessage():
			log.Println("*** ESCAPE message ***")
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

// btnEscNone updates app-status in storage to greetings.
func (b *Buttons) btnEscMenuMain() {
	b.store.App.SetGreetings()
	// update render
	b.render <- struct{}{}
}

// btnEscMessage updates app-status in storage to menu-main.
func (b *Buttons) btnEscMessage() {
	b.store.App.SetMenuMain()
	// update render according to new app-status
	b.render <- struct{}{}
}
