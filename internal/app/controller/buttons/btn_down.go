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
		case b.store.App.IsMessage():
			log.Println("*** DOWN message ***")
			b.btnDownMessage()
		default:
			log.Println("*** DOWN pressed ***")
		}
	}
}

// btnDownMenuMain select the next item in menu.
func (b *Buttons) btnDownMenuMain() {
	// scroll down menu
	menu := b.store.Menu.GetMain()
	menu.SelectNext()
	// update render with new menu view
	b.render <- struct{}{}
}

// btnDownMessage scroll message text down for one line.
func (b *Buttons) btnDownMessage() {
	// scroll down message
	msg := b.store.Message.Get()
	msg.ScrollDown()
	// update render with new message view
	b.render <- struct{}{}
}
