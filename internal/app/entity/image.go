package entity

import (
	"fmt"

	"IvolgaOledManager/internal/pkg/display"
)

// Image is a type to render image on display.
type Image struct {
	// path to image
	ImagePath string
}

// Render implements display.Renderer. It renders image on display.
func (g *Image) Render(device display.Display) error {
	if err := device.DisplayImage(g.ImagePath, 0, 0); err != nil {
		return fmt.Errorf("display image: %w", err)
	}
	return nil
}
