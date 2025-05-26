package renderer

import (
	"fmt"
)

// Greetings renders greetings image.
func (d display) Greetings() error {
	if err := d.device.DisplayImage(d.greetingsImgPath, 0, 0); err != nil {
		return fmt.Errorf("display image: %w", err)
	}
	return nil
}
