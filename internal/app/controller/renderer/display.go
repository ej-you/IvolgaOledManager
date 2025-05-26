package renderer

import (
	"fmt"

	"sschmc/internal/pkg/ssd1306"
)

var _ Renderer = (*display)(nil)

// Renderer implementation.
type display struct {
	device           *ssd1306.SSD1306
	greetingsImgPath string
}

func NewDisplay(bus, greetingsImgPath string) (Renderer, error) {
	oled, err := ssd1306.NewSSD1306(bus)
	if err != nil {
		return nil, fmt.Errorf("connect to oled: %w", err)
	}

	return &display{
		device:           oled,
		greetingsImgPath: greetingsImgPath,
	}, nil
}

// Close clears image and closes display connection.
func (d display) Close() error {
	if err := d.Clear(); err != nil {
		return err
	}
	return d.device.Close()
}

// Clear clears image.
func (d display) Clear() error {
	return d.device.DisplayClear()
}
