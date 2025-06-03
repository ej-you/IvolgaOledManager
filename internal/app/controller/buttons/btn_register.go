package buttons

import (
	"sschmc/internal/pkg/gpiobutton"
)

// BtnEntRisingHandler handles all cases of ENT button rising.
func (b *Buttons) BtnEntRisingHandler() gpiobutton.HandlerFunc {
	return func() {
		switch {
		case b.store.App.IsNone():
			b.screenGreetings()
		case b.store.App.IsGreetings():
			b.screenMenuMain()
		case b.store.App.IsMenuMain():
			b.screenMenuLevel()
		case b.store.App.IsMenuLevel():
			b.screenMessage()
		case b.store.App.IsMessage():
			b.deleteMessage()
		}
	}
}

// BtnEscRisingHandler handles all cases of ESC button rising.
func (b *Buttons) BtnEscRisingHandler() gpiobutton.HandlerFunc {
	return func() {
		switch {
		case b.store.App.IsNone(), b.store.App.IsMenuMain():
			b.screenGreetings()
		case b.store.App.IsGreetings():
			b.screenNone()
		case b.store.App.IsMenuLevel():
			b.screenMenuMain()
		case b.store.App.IsMessage():
			b.screenMenuLevel()
		}
	}
}

// BtnUpRisingHandler handles all cases of UP button rising.
func (b *Buttons) BtnUpRisingHandler() gpiobutton.HandlerFunc {
	return func() {
		switch {
		case b.store.App.IsNone():
			b.screenGreetings()
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
			b.screenGreetings()
		case b.store.App.IsMenuMain():
			b.menuScrollDown(b.store.Menu.GetMain())
		case b.store.App.IsMenuLevel():
			b.menuScrollDown(b.store.Menu.GetLevel())
		case b.store.App.IsMessage():
			b.messageScrollDown(b.store.Message.Get())
		}
	}
}
