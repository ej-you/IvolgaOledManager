package renderer

import (
	"fmt"
)

// Menu renders menu.
func (d display) Menu() error {
	// if err := d.device.DisplayText("Tgoogle hello world * {}\n\"@ json {} O", 0, 0); err != nil {
	if err := d.device.DisplayText("Line0 []\nLine1 {}\nLine2 ()\nLine3 $\nTgoogle1hello2wordNEW LINE!!!", 0, 0); err != nil {
		return fmt.Errorf("display image: %w", err)
	}
	return nil
}
