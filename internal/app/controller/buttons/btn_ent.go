package buttons

import (
	"fmt"
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
		appStatus = b.store.Get(constants.KeyAppStatus)

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
	if err := b.render.Menu(); err != nil {
		return fmt.Errorf("render menu: %w", err)
	}
	b.store.Set(constants.KeyAppStatus, constants.ValueAppStatusMenuMain)
	return nil
}
