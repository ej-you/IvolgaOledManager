// Package display provides functions for output different information
// to SSD1306 OLED display.
package display

import (
	"fmt"

	"github.com/pkg/errors"

	"sschmc/internal/pkg/ssd1306"
)

var _ Display = (*display)(nil)

// Interface with functions for all output cases.
type Display interface {
	Close() error

	Greetings() error
}

type display struct {
	device           *ssd1306.SSD1306
	greetingsImgPath string
}

func New(bus, greetingsImgPath string) (Display, error) {
	oled, err := ssd1306.NewSSD1306(bus)
	if err != nil {
		return nil, fmt.Errorf("connect to oled: %w", err)
	}

	return &display{
		device:           oled,
		greetingsImgPath: greetingsImgPath,
	}, nil
}

// Close closes display connection.
func (d display) Close() error {
	err := d.device.Close()
	return errors.Wrap(err, "close display device")
}

// Display greetings image.
func (d display) Greetings() error {
	if err := d.device.DisplayImage(d.greetingsImgPath, -32, 0); err != nil {
		return fmt.Errorf("display image: %w", err)
	}
	return nil
}
