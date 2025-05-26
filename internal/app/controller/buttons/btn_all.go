package buttons

import (
	"fmt"

	"sschmc/internal/app/constants"
)

// btnAllGreetings renders greetings and updates app-status in storage.
func (b *Buttons) btnAllGreetings() error {
	if err := b.render.Greetings(); err != nil {
		return fmt.Errorf("render greetings: %w", err)
	}
	b.store.Set(constants.KeyAppStatus, constants.ValueAppStatusGreetings)
	return nil
}
