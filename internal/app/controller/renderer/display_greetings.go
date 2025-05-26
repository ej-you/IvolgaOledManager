package renderer

import (
	"fmt"
)

// greetings renders greetings image.
func (r *Renderer) greetings() error {
	if err := r.device.DisplayImage(r.greetingsImgPath, 0, 0); err != nil {
		return fmt.Errorf("display image: %w", err)
	}
	return nil
}
