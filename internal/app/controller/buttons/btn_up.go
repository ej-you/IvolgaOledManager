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
		case b.store.App.IsMenuMain():
			log.Println("*** UP menu ***")
			b.btnUpMenuMain()
		case b.store.App.IsMessage():
			log.Println("*** UP message ***")
			b.btnUpMessage()
		default:
			log.Println("*** UP pressed ***")
		}
	}
}

// btnUpMenuMain select the previous item in menu.
func (b *Buttons) btnUpMenuMain() {
	// scroll down menu
	menu := b.store.Menu.GetMain()
	menu.SelectPrevious()
	// update render with new menu view
	b.render <- struct{}{}
}

// btnDownMessage scroll message text up for one line.
func (b *Buttons) btnUpMessage() {
	// scroll up message
	msg := b.store.Message.Get()
	msg.ScrollUp()
	// update render with new message view
	b.render <- struct{}{}
}
