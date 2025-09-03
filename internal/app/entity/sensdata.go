// Package entity contains all app entities.
package entity

import (
	"fmt"

	"IvolgaOledManager/internal/pkg/display"
)

// Sensdata is a data from station sensor.
type Sensdata struct {
	// measurement name to output
	Title string
	// sensor data value to output
	Data string

	// raw sensor data
	RawData float64
}

// Render implements display.Renderer. It renders sensor data on display.
func (s *Sensdata) Render(device display.Display) error {
	// output text screen
	err := device.DisplayTextLines(
		display.TextLine{Content: s.Title, RelativeSize: 2},
		display.TextLine{Content: "", RelativeSize: 1},
		display.TextLine{Content: s.Data, RelativeSize: 4},
		display.TextLine{Content: "", RelativeSize: 1},
	)
	if err != nil {
		return fmt.Errorf("display text screen: %w", err)
	}
	return nil
}
