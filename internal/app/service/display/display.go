// Package display provides service to output data on display.
package display

import (
	"IvolgaOledManager/config"
	"IvolgaOledManager/internal/pkg/display"
)

// New returns a new instance of display service.
func New(cfg *config.Config) (*display.Service, error) {
	return display.NewDisplayService(
		cfg.Hardware.Oled.Bus,
		cfg.Hardware.Oled.Width,
		cfg.Hardware.Oled.Height,
		cfg.Hardware.Oled.UpdatesDuration,
	)
}
