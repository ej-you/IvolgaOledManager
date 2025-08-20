package buttons

import (
	"IvolgaOledManager/internal/pkg/gpiobutton"
)

// BtnEntRisingHandler handles all cases of ENT button rising.
func (b *Buttons) BtnEntRisingHandler() gpiobutton.HandlerFunc {
	return func() {
		switch {
		case b.store.App.IsNone():
			b.screenGreetings()
		case b.store.App.IsGreetings(), b.store.App.IsStationResult():
			b.screenStationResult()
		}
	}
}

// BtnEscRisingHandler handles all cases of ESC button rising.
func (b *Buttons) BtnEscRisingHandler() gpiobutton.HandlerFunc {
	return func() {
		switch {
		case b.store.App.IsNone(), b.store.App.IsStationResult():
			b.screenGreetings()
		case b.store.App.IsGreetings():
			b.screenNone()
		}
	}
}

// BtnUpRisingHandler handles all cases of UP button rising.
func (b *Buttons) BtnUpRisingHandler() gpiobutton.HandlerFunc {
	return func() {
		switch {
		case b.store.App.IsNone():
			b.screenGreetings()
		case b.store.App.IsStationResult():
			b.nextStationResult()
		}
	}
}

// BtnDownRisingHandler handles all cases of DOWN button rising.
func (b *Buttons) BtnDownRisingHandler() gpiobutton.HandlerFunc {
	return func() {
		switch {
		case b.store.App.IsNone():
			b.screenGreetings()
		case b.store.App.IsStationResult():
			b.previousStationResult()
		}
	}
}
