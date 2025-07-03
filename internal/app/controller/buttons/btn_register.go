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
		case b.store.App.IsGreetings():
			b.screenMenuMain()
		// first main menu branch
		case b.store.App.IsMenuMain():
			b.screenMenuLogsOrStation()
		case b.store.App.IsMenuLogs():
			b.screenMenuLevel()
		case b.store.App.IsMenuLevel():
			b.screenMessage()
		case b.store.App.IsMessage():
			b.deleteMessage()
		// second main menu branch
		case b.store.App.IsMenuStation():
			b.screenSensor()
		case b.store.App.IsSensor():
			b.updateStation()
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
		// first main menu branch
		case b.store.App.IsMenuLogs():
			b.screenMenuMain()
		case b.store.App.IsMenuLevel():
			b.screenMenuLogs()
		case b.store.App.IsMessage():
			b.screenMenuLevel()
		// second main menu branch
		case b.store.App.IsMenuStation():
			b.screenMenuMain()
		case b.store.App.IsSensor():
			b.screenMenuStation()
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
		// first main menu branch
		case b.store.App.IsMenuLogs():
			b.menuScrollUp(b.store.Menu.GetLogs())
		case b.store.App.IsMenuLevel():
			b.menuScrollUp(b.store.Menu.GetLevel())
		case b.store.App.IsMessage():
			b.messageScrollUp(b.store.Message.Get())
		// second main menu branch
		case b.store.App.IsMenuStation():
			b.menuScrollUp(b.store.Menu.GetStation())
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
		// first main menu branch
		case b.store.App.IsMenuLogs():
			b.menuScrollDown(b.store.Menu.GetLogs())
		case b.store.App.IsMenuLevel():
			b.menuScrollDown(b.store.Menu.GetLevel())
		case b.store.App.IsMessage():
			b.messageScrollDown(b.store.Message.Get())
		// second main menu branch
		case b.store.App.IsMenuStation():
			b.menuScrollDown(b.store.Menu.GetStation())
		}
	}
}
