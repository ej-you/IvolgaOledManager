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
			b.btnAllGreetings()
		case b.store.App.IsGreetings():
			b.btnEntMenuMain()
		case b.store.App.IsMenuMain():
			b.btnEntMenuLevel()
		case b.store.App.IsMenuLevel():
			b.btnEntMessage()
		default:
			log.Println("*** ENTER pressed ***")
		}
	}
}

// btnEntMenuMain update status to menu-main and creates render task.
func (b *Buttons) btnEntMenuMain() {
	// create main menu and save it to storage
	mainMenu := &entity.Menu{
		Title: "Main menu",
		Items: []*entity.MenuItem{
			entity.NewMenuItem("Привет, мир", ""),
			entity.NewMenuItem("Какого чёрта ты делаешь, чувак?", ""),
			entity.NewMenuItem("Golang", ""),
			entity.NewMenuItem("func{}", ""),
			entity.NewMenuItem("sample", ""),
			entity.NewMenuItem("Как good оно работает!", ""),
		},
	}
	b.store.Menu.SetMain(mainMenu)

	b.store.App.SetMenuMain()
	// update render according to new app-status
	b.render <- struct{}{}
}

// btnEntMenuLevel update status to menu-level and creates render task.
func (b *Buttons) btnEntMenuLevel() {
	// create level menu and save it to storage
	levelMenu := &entity.Menu{
		Title: "Level menu",
		Items: []*entity.MenuItem{
			entity.NewMenuItem("menu 2", ""),
			entity.NewMenuItem("ну и ещё пункт для приличия", ""),
		},
	}
	b.store.Menu.SetLevel(levelMenu)

	b.store.App.SetMenuLevel()
	// update render according to new app-status
	b.render <- struct{}{}
}

// btnEntMessage clears rendered data and updates app-status in storage to none.
func (b *Buttons) btnEntMessage() {
	menu := b.store.Menu.GetLevel()

	var text string
	if menu.SelectedItem == 0 {
		text = "This is a very short message"
	} else {
		text = `The quick brown fox jumped over the lazy dog.
		Какого чёрта ты делаешь, чувак?
		This is electromagnetically!
		A big crocodile died empty-fanged, gulping horribly in jerking kicking little
		motions. Nonchalant old Peter Quinn ruthlessly shot the under-water vermin with
		Xavier yelling Zap!`
	}

	// create main menu and save it to storage
	msg := &entity.Message{
		Content: text,
	}
	msg.Format()
	b.store.Message.Set(msg)

	b.store.App.SetMessage()
	// update render according to new app-status
	b.render <- struct{}{}
}
