package ssd1306

import (
	"fmt"
	"image"

	"periph.io/x/devices/v3/ssd1306/image1bit"
)

// DisplayClear clears OLED display.
func (s *SSD1306) DisplayClear() error {
	black := image1bit.NewVerticalLSB(image.Rect(0, 0, _displayWidth, _displayHeight))
	if err := s.device.Draw(s.device.Bounds(), black, image.Pt(0, 0)); err != nil {
		return fmt.Errorf("clear display: %w", err)
	}
	return nil
}
