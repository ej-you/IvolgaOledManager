package renderer

import (
	"fmt"
)

const (
	_valueAppStatusMenuMain = "menu-main"
	_defaultPrefix          = "  "
	_selectedPrefix         = "> "
)

// menu renders menu.
func (r *Renderer) menu() error {
	// menu := r.store.Menu.Get(_valueAppStatusMenuMain)

	// if err := r.device.DisplayText("Line0 []\nLine1 {}\nLine2 ()\nLine3 $\nTgoogle1hello2wordNEW LINE!!!", 0, 0); err != nil {
	// 	return fmt.Errorf("display image: %w", err)
	// }

	drawer := r.device.NewTextDrawer()
	drawer.AddLine(_defaultPrefix, "hello")
	drawer.AddLine(_defaultPrefix, "{}[]()")
	drawer.AddLine(_selectedPrefix, "join")
	drawer.AddLine(_defaultPrefix, "google")
	if err := drawer.Draw(); err != nil {
		return fmt.Errorf("display text lines: %w", err)
	}
	return nil
}
