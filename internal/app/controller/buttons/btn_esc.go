package buttons

import (
	"fmt"
	"log"

	"sschmc/internal/app/constants"
	"sschmc/internal/pkg/errlog"
	"sschmc/internal/pkg/gpiobutton"
)

// BtnEscRisingHandler handles all cases of ESC button rising.
func (b *Buttons) BtnEscRisingHandler() gpiobutton.HandlerFunc {
	var (
		err       error
		appStatus string
	)
	return func() {
		appStatus = b.store.Get(constants.KeyAppStatus)

		switch appStatus {
		case constants.ValueAppStatusNone, constants.ValueAppStatusMenuMain:
			log.Println("*** ESCAPE none/menu ***")
			err = b.btnAllGreetings()
		case constants.ValueAppStatusGreetings:
			log.Println("*** ESCAPE greetings ***")
			err = b.btnEscNone()
		default:
			log.Println("*** ESCAPE pressed ***")
		}

		if err != nil {
			errlog.Print(err)
		}
	}
}

// btnEscNone clears rendered data and updates app-status in storage to none.
func (b *Buttons) btnEscNone() error {
	if err := b.render.Clear(); err != nil {
		return fmt.Errorf("clear rendered: %w", err)
	}
	b.store.Set(constants.KeyAppStatus, constants.ValueAppStatusNone)
	return nil
}
