package buttons

import (
	"log"

	"sschmc/internal/app/constants"
	"sschmc/internal/pkg/errlog"
	"sschmc/internal/pkg/gpiobutton"
)

// BtnEntRisingHandler handles all cases of ENT button rising.
func (b *Buttons) BtnEntRisingHandler() gpiobutton.HandlerFunc {
	var (
		err       error
		appStatus string
	)
	return func() {
		appStatus = b.store.App.GetStatus()

		switch appStatus {
		case constants.ValueAppStatusNone:
			log.Println("*** ENTER none ***")
			err = b.btnAllGreetings()
		case constants.ValueAppStatusGreetings:
			log.Println("*** ENTER greetings ***")
			err = b.btnEntMainMenu()
		default:
			log.Println("*** ENTER pressed ***")
		}

		if err != nil {
			errlog.Print(err)
		}
	}
}

// btnEscNone clears rendered data and updates app-status in storage to none.
func (b *Buttons) btnEntMainMenu() error {
	b.store.App.SetStatus(constants.ValueAppStatusMenuMain)
	// update render according to new app-status
	b.render <- struct{}{}

	// if err := b.render.Menu(); err != nil {
	// 	return fmt.Errorf("render menu: %w", err)
	// }
	return nil
}
