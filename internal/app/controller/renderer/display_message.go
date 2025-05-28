package renderer

import (
	"fmt"

	"sschmc/internal/app/entity"
)

// message renders log message.
func (r *Renderer) message(msg *entity.Message) error {
	drawer := r.device.NewTextDrawer()

	// count limit lines amount
	limit := min(entity.MaxDisplayedLines, len(msg.Lines[msg.FirstLine:]))

	for _, line := range msg.Lines[msg.FirstLine : msg.FirstLine+limit] {
		drawer.AddLine("", line)
	}
	drawer.FillEmpty()

	if err := drawer.Draw(); err != nil {
		return fmt.Errorf("display text lines: %w", err)
	}
	return nil
}
