package renderer

import (
	"fmt"

	"sschmc/internal/app/entity"
)

// menu renders menu.
func (r *Renderer) menu(menu *entity.Menu) error {
	drawer := r.device.NewTextDrawer()

	fmt.Println("menu.FirstItem:", menu.FirstItem)
	fmt.Println("menu.SelectedItem:", menu.SelectedItem)

	for idx, menuItem := range menu.Items {
		// skip items before first visible item
		if idx < menu.FirstItem {
			continue
		}
		// add line in current state to drawer
		if idx == menu.SelectedItem {
			drawer.AddLine(entity.SelectedPrefix, menuItem.Title[menuItem.FirstSymbol:])
			// scroll item if need
			menuItem.Scroll()
		} else {
			drawer.AddLine(entity.DefaultPrefix, menuItem.Title[menuItem.FirstSymbol:])
		}
	}
	drawer.FillEmpty()
	// drawer.AddLine(_defaultPrefix, "hello")
	// drawer.AddLine(_defaultPrefix, "{}[]()")
	// drawer.AddLine(_selectedPrefix, "join")
	// drawer.AddLine(_defaultPrefix, "google")
	if err := drawer.Draw(); err != nil {
		return fmt.Errorf("display text lines: %w", err)
	}
	return nil
}
