package buttons

import (
	"sschmc/internal/app/constants"
)

// btnAllGreetings renders greetings and updates app-status in storage.
func (b *Buttons) btnAllGreetings() error {
	b.store.App.SetStatus(constants.ValueAppStatusGreetings)
	// update render according to new app-status
	b.render <- struct{}{}

	// if err := b.render.Greetings(); err != nil {
	// 	return fmt.Errorf("render greetings: %w", err)
	// }
	return nil
}
