package renderer

import (
	"fmt"

	"sschmc/internal/app/entity"
)

// menu renders menu.
func (r *Renderer) menu(menu *entity.Menu) error {
	drawer, err := r.device.NewTextDrawer()
	if err != nil {
		return fmt.Errorf("create text drawer: %w", err)
	}

	drawer.AddLine("", menu.Title)
	for idx, menuItem := range menu.Items {
		// skip items before first visible item
		if idx < menu.FirstItem {
			continue
		}
		// add line in current state to drawer
		if idx == menu.SelectedItem {
			drawer.AddLine(entity.SelectedPrefix, menuItem.FormattedTitle())
			// scroll item if need
			menuItem.Scroll()
		} else {
			drawer.AddLine(entity.DefaultPrefix, menuItem.FormattedTitle())
		}
	}
	drawer.FillEmpty()
	if err := drawer.Draw(); err != nil {
		return fmt.Errorf("display text lines: %w", err)
	}
	return nil
}
