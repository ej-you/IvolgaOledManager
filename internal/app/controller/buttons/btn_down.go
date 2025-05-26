package buttons

import (
	"log"

	"sschmc/internal/app/constants"
	"sschmc/internal/pkg/errlog"
	"sschmc/internal/pkg/gpiobutton"
)

// BtnDownRisingHandler handles all cases of DOWN button rising.
func (b *Buttons) BtnDownRisingHandler() gpiobutton.HandlerFunc {
	var (
		err       error
		appStatus string
	)
	return func() {
		appStatus = b.store.Get(constants.KeyAppStatus)

		switch appStatus {
		case constants.ValueAppStatusNone:
			err = b.btnAllGreetings()
		default:
			log.Println("*** DOWN pressed ***")
		}

		if err != nil {
			errlog.Print(err)
		}
	}
}
