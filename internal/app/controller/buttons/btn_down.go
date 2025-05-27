package buttons

import (
	"log"

	"sschmc/internal/pkg/gpiobutton"
)

// BtnDownRisingHandler handles all cases of DOWN button rising.
func (b *Buttons) BtnDownRisingHandler() gpiobutton.HandlerFunc {
	return func() {
		switch {
		case b.store.App.IsNone():
			b.btnAllGreetings()
		case b.store.App.IsMenuMain():
			log.Println("*** DOWN menu ***")
			b.btnDownMenuMain()
		default:
			log.Println("*** DOWN pressed ***")
		}
	}
}

// btnEscNone clears rendered data and updates app-status in storage to none.
func (b *Buttons) btnDownMenuMain() {
	// scroll down menu
	menu := b.store.Menu.GetMenuMain()
	menu.SelectNext()
	// update render with new menu view
	b.render <- struct{}{}
}
