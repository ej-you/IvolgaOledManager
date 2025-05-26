package renderer

import (
	"fmt"
)

// Message renders log message.
func (d display) Message() error {
	if err := d.device.DisplayText("Tgoogle1hello2wordNEW LINE!!!", 0, 0); err != nil {
		return fmt.Errorf("display image: %w", err)
	}
	return nil
}
