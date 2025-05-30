// Package ssd1306 contains functions and types to controlling SSD1306 OLED display via I2C.
package ssd1306

import (
	"fmt"

	"github.com/pkg/errors"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/devices/v3/ssd1306"
)

const (
	_displayWidth  = 128 // oled width
	_displayHeight = 64  // oled height
)

type SSD1306 struct {
	device    *ssd1306.Dev
	busCloser func() error
}

func NewSSD1306(bus string) (*SSD1306, error) {
	busCloser, err := i2creg.Open(bus)
	if err != nil {
		return nil, fmt.Errorf("open bus: %w", err)
	}
	device, err := ssd1306.NewI2C(busCloser, &ssd1306.DefaultOpts)
	if err != nil {
		return nil, fmt.Errorf("open device: %w", err)
	}

	instance := &SSD1306{
		device:    device,
		busCloser: busCloser.Close,
	}
	if err := instance.DisplayClear(); err != nil {
		return nil, fmt.Errorf("clear device on startapp: %w", err)
	}

	return instance, nil
}

// Close clears OLED display (turn off all pixels) and closes OLED bus.
func (s *SSD1306) Close() error {
	if err := s.device.Halt(); err != nil {
		return fmt.Errorf("clear ssd1306 oled: %w", err)
	}
	return errors.Wrap(s.busCloser(), "close ssd1306 oled")
}
