package renderer

import (
	"fmt"

	"IvolgaOledManager/internal/app/entity"
)

// station renders station result.
func (r *Renderer) station(statRes *entity.StationResult) error {
	drawer, err := r.device.NewTextDrawer()
	if err != nil {
		return fmt.Errorf("create text drawer: %w", err)
	}

	drawer.AddLine("", statRes.Title())
	drawer.AddLine("", "")
	drawer.AddLine("", statRes.ResultText())
	drawer.FillEmpty()

	if err := drawer.Draw(); err != nil {
		return fmt.Errorf("display text lines: %w", err)
	}
	return nil
}
