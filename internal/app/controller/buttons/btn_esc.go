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
			b.btnEscNone()
		case b.store.App.IsMenuMain():
			log.Println("*** ESCAPE (UP) menu ***")
			b.btnEscMenuMain()
		default:
			log.Println("*** ESCAPE pressed ***")
		}
	}
}

// btnEscNone clears rendered data and updates app-status in storage to none.
func (b *Buttons) btnEscNone() {
	b.store.App.SetNone()
	// update render according to new app-status
	b.render <- struct{}{}
}

// TODO: temp
func (b *Buttons) btnEscMenuMain() {
	// scroll down menu
	menu := b.store.Menu.GetMenuMain()
	menu.SelectPrevious()
	// update render with new menu view
	b.render <- struct{}{}
}
