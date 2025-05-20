package ssd1306

import (
	"fmt"
)

// DisplayClear clears OLED display (turn off all pixels).
func (s SSD1306) DisplayClear() error {
	if err := s.device.Halt(); err != nil {
		return fmt.Errorf("clear display: %w", err)
	}
	return nil
}
