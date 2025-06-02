package renderer

import (
	"fmt"

	"sschmc/internal/app/entity"
)

// message renders log message.
func (r *Renderer) message(msg *entity.Message) error {
	drawer, err := r.device.NewTextDrawer()
	if err != nil {
		return fmt.Errorf("create text drawer: %w", err)
	}

	// count limit lines amount
	limit := min(entity.MaxDisplayedItems, len(msg.Lines[msg.FirstLine:]))

	drawer.AddLine("", msg.Datetime())
	for _, line := range msg.Lines[msg.FirstLine : msg.FirstLine+limit] {
		drawer.AddLine("", line)
	}
	drawer.FillEmpty()

	if err := drawer.Draw(); err != nil {
		return fmt.Errorf("display text lines: %w", err)
	}
	return nil
}
