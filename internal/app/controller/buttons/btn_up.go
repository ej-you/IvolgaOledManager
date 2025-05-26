package buttons

import (
	"log"
	"sschmc/internal/app/constants"
	"sschmc/internal/pkg/errlog"
	"sschmc/internal/pkg/gpiobutton"
)

// BtnUpRisingHandler handles all cases of UP button rising.
func (b *Buttons) BtnUpRisingHandler() gpiobutton.HandlerFunc {
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
			log.Println("*** UP pressed ***")
		}

		if err != nil {
			errlog.Print(err)
		}
	}
}
