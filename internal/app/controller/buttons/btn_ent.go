package buttons

import (
	"log"

	"sschmc/internal/app/entity"
	"sschmc/internal/pkg/gpiobutton"
)

// BtnEntRisingHandler handles all cases of ENT button rising.
func (b *Buttons) BtnEntRisingHandler() gpiobutton.HandlerFunc {
	return func() {
		switch {
		case b.store.App.IsNone():
			log.Println("*** ENTER none ***")
			b.btnAllGreetings()
		case b.store.App.IsGreetings():
			log.Println("*** ENTER greetings ***")
			b.btnEntMenuMain()
		default:
			log.Println("*** ENTER pressed ***")
		}
	}
}

// btnEscNone clears rendered data and updates app-status in storage to none.
func (b *Buttons) btnEntMenuMain() {
	// create main menu and save it to storage
	mainMenu := &entity.Menu{
		Title: "Main menu",
		Items: []*entity.MenuItem{
			{Title: "hello"},
			{Title: "What a hell are you doing?"},
			{Title: "Golang"},
			{Title: "func{}"},
			{Title: "sample"},
			{Title: "B A A A A A A A C"},
		},
		SelectedItem: 1,
	}
	b.store.Menu.SetMenuMain(mainMenu)

	b.store.App.SetMenuMain()
	// update render according to new app-status
	b.render <- struct{}{}
}
