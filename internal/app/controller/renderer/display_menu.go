package renderer

import (
	"fmt"
)

// menu renders menu.
func (r *Renderer) menu() error {
	// if err := d.device.DisplayText("Tgoogle hello world * {}\n\"@ json {} O", 0, 0); err != nil {
	if err := r.device.DisplayText("Line0 []\nLine1 {}\nLine2 ()\nLine3 $\nTgoogle1hello2wordNEW LINE!!!", 0, 0); err != nil {
		return fmt.Errorf("display image: %w", err)
	}
	return nil
}
