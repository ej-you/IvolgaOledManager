// Package button provides service to output data on display.
package display

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"IvolgaOledManager/internal/pkg/display/ssd1306"
)

// Ensure display service implements interface
var _ Display = (*Service)(nil)

// Display describes all methods for display to output data.
type Display interface {
	// DisplayImage displays image from given path on the OLED-display.
	DisplayImage(path string, x, y int) error
	// DisplayTextLines displays given text lines.
	DisplayTextLines(lines ...TextLine) error
}

// TextLine represents a one text line on display screen.
type TextLine struct {
	// text
	Content string
	// relative part of the total display height occupied by the text line
	RelativeSize int
}

// Service represents a display service. It starts with StartWithShutdown method
// and implements Display interface to output data.
type Service struct {
	device          *ssd1306.SSD1306
	updatesDuration time.Duration
	updateNow       chan struct{} // used to force screen update when updating the renderer
	renderer        Renderer
	mu              sync.RWMutex
}

// Renderer describes object that can be rendered on display through device.
type Renderer interface {
	Render(display Display) error
}

// NewDisplay connects to display device and returns display service.
func NewDisplayService(bus string, width, height int,
	updatesDuration time.Duration) (*Service, error) {

	device, err := ssd1306.NewSSD1306(bus, width, height)
	if err != nil {
		return nil, fmt.Errorf("ssd1306: %w", err)
	}
	return &Service{
		device:          device,
		updatesDuration: updatesDuration,
		updateNow:       make(chan struct{}),
	}, nil
}

// DisplayTextLines displays given text lines.
// It is an adaptor for the same method
// of ssd1306.SSD1306 to implement Display interface.
func (d *Service) DisplayTextLines(lines ...TextLine) error {
	ssd1306Lines := make([]ssd1306.TextLine, 0, len(lines))

	for _, line := range lines {
		ssd1306Lines = append(ssd1306Lines, ssd1306.TextLine{
			Content:      line.Content,
			RelativeSize: line.RelativeSize,
		})
	}
	return d.device.DisplayTextLines(ssd1306Lines...)
}

// DisplayImage displays image from given path on the OLED-display.
func (d *Service) DisplayImage(path string, x, y int) error {
	return d.device.DisplayImage(path, x, y)
}

// SetRenderer sets new renderer object for display.
func (d *Service) SetRenderer(renderer Renderer) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.renderer = renderer

	// send force update signal to chan (if it is not filled)
	select {
	case d.updateNow <- struct{}{}:
	default:
	}
}

// StartWithShutdown starts display.
// Service will stop when given context will be stopped.
// This method is blocking.
func (d *Service) StartWithShutdown(ctx context.Context) error {
	logrus.Info("start display service...")

	// start display updates loop
	d.start(ctx)
	// close display after ctx is done
	if err := d.close(); err != nil {
		return err
	}
	return nil
}

// start starts display updates loop.
func (d *Service) start(ctx context.Context) {
	// init ticker for display periodically updates
	ticker := time.NewTicker(d.updatesDuration)
	defer ticker.Stop()

	var err error
	for {
		select {
		// periodically updates
		case <-ticker.C:
			err = d.render()
		// force update after setting new renderer
		case <-d.updateNow:
			err = d.render()
		// if main context is done
		case <-ctx.Done():
			return
		}
		if err != nil {
			logrus.Errorf("render object: %v", err)
		}
	}
}

// render renders renderer object to display.
func (d *Service) render() error {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if d.renderer == nil {
		return nil
	}
	return d.renderer.Render(d)
}

// close clears display screen and closes display connection.
func (d *Service) close() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if err := d.device.DisplayClear(); err != nil {
		logrus.Errorf("clear screen on display close: %v", err)
	}
	return d.device.Close()
}
