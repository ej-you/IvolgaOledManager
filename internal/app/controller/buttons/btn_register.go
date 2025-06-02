package buttons

import (
	"sschmc/internal/pkg/gpiobutton"
)

// BtnEntRisingHandler handles all cases of ENT button rising.
func (b *Buttons) BtnEntRisingHandler() gpiobutton.HandlerFunc {
	return func() {
		switch {
		case b.store.App.IsNone():
			b.toGreetings()
		case b.store.App.IsGreetings():
			b.toMenuMain()
		case b.store.App.IsMenuMain():
			b.toMenuLevel()
		case b.store.App.IsMenuLevel():
			b.toMessage()
		}
	}
}

// BtnEscRisingHandler handles all cases of ESC button rising.
func (b *Buttons) BtnEscRisingHandler() gpiobutton.HandlerFunc {
	return func() {
		switch {
		case b.store.App.IsNone(), b.store.App.IsMenuMain():
			b.toGreetings()
		case b.store.App.IsGreetings():
			b.backToNone()
		case b.store.App.IsMenuLevel():
			b.backToMenuMain()
		case b.store.App.IsMessage():
			b.backToMenuLevel()
		}
	}
}

// BtnUpRisingHandler handles all cases of UP button rising.
func (b *Buttons) BtnUpRisingHandler() gpiobutton.HandlerFunc {
	return func() {
		switch {
		case b.store.App.IsNone():
			b.toGreetings()
		case b.store.App.IsMenuMain():
			b.menuScrollUp(b.store.Menu.GetMain())
		case b.store.App.IsMenuLevel():
			b.menuScrollUp(b.store.Menu.GetLevel())
		case b.store.App.IsMessage():
			b.messageScrollUp(b.store.Message.Get())
		}
	}
}

// BtnDownRisingHandler handles all cases of DOWN button rising.
func (b *Buttons) BtnDownRisingHandler() gpiobutton.HandlerFunc {
	return func() {
		switch {
		case b.store.App.IsNone():
			b.toGreetings()
		case b.store.App.IsMenuMain():
			b.menuScrollDown(b.store.Menu.GetMain())
		case b.store.App.IsMenuLevel():
			b.menuScrollDown(b.store.Menu.GetLevel())
		case b.store.App.IsMessage():
			b.messageScrollDown(b.store.Message.Get())
		}
	}
}
