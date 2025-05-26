// Package renderer provides Renderer interface with functions for
// output data to SSD1306 OLED display.
package renderer

// Interface with functions for all output cases.
type Renderer interface {
	Close() error
	Clear() error

	Greetings() error
	Menu() error
	Message() error
}
