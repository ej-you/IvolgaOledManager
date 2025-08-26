// Package ssd1306 provides interface to controlling SSD1306 OLED-display via I2C.
package ssd1306

import (
	"fmt"

	"github.com/pkg/errors"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/devices/v3/ssd1306"
)

// SSD1306 represents an OLED-display.
type SSD1306 struct {
	device       *ssd1306.Dev
	busCloser    func() error
	screenWidth  int
	screenHeight int
}

// NewSSD1306 returns a new instance of SSD1306.
func NewSSD1306(bus string, screenWidth, screenHeight int) (*SSD1306, error) {
	busCloser, err := i2creg.Open(bus)
	if err != nil {
		return nil, fmt.Errorf("open bus: %w", err)
	}
	device, err := ssd1306.NewI2C(busCloser, &ssd1306.DefaultOpts)
	if err != nil {
		return nil, fmt.Errorf("open device: %w", err)
	}

	instance := &SSD1306{
		device:       device,
		busCloser:    busCloser.Close,
		screenWidth:  screenWidth,
		screenHeight: screenHeight,
	}
	if err := instance.DisplayClear(); err != nil {
		return nil, fmt.Errorf("clear oled on startup: %w", err)
	}
	return instance, nil
}

// Close closes OLED bus used by OLED-display.
func (s *SSD1306) Close() error {
	if err := s.device.Halt(); err != nil {
		return fmt.Errorf("clear oled: %w", err)
	}
	return errors.Wrap(s.busCloser(), "close oled bus")
}
