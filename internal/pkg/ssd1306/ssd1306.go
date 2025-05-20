// Package ssd1306 contains functions and types to controlling SSD1306 OLED display via I2C.
package ssd1306

import (
	"fmt"

	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/devices/v3/ssd1306"
	"periph.io/x/host/v3"
)

type SSD1306 struct {
	device    *ssd1306.Dev
	busCloser func() error
}

func NewSSD1306(bus string) (*SSD1306, error) {
	// initialise all relevant drivers
	if _, err := host.Init(); err != nil {
		return nil, fmt.Errorf("init drivers: %w", err)
	}

	busCloser, err := i2creg.Open(bus)
	if err != nil {
		return nil, fmt.Errorf("open bus: %w", err)
	}
	device, err := ssd1306.NewI2C(busCloser, &ssd1306.DefaultOpts)
	if err != nil {
		return nil, fmt.Errorf("open device: %w", err)
	}

	return &SSD1306{
		device:    device,
		busCloser: busCloser.Close,
	}, nil
}

// Close closes OLED bus.
func (s SSD1306) Close() error {
	return s.busCloser()
}
