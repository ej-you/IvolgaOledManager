// Package button provides service to output data on display.
package display

import (
	"IvolgaOledManager/config"
	"IvolgaOledManager/internal/pkg/ssd1306"
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

var _ Device = (*ssd1306.SSD1306)(nil)

// Display device for display service.
type Device interface {
	Close() error
	DisplayClear() error
	DisplayImage(path string, x, y int) error
	DisplayTextLines(lines ...ssd1306.TextLine) error
}

// Renderer describes object that can be rendered on display.
type Renderer interface {
	Render(device Device) error
}

// Display service.
type Display struct {
	device          Device
	updatesDuration time.Duration

	updateNow chan struct{} // used to force screen update when updating the renderer
	renderer  Renderer
	mu        sync.RWMutex
}

// NewDisplay connects to display device and returns display service.
func NewDisplay(cfg *config.Config) (*Display, error) {
	device, err := ssd1306.NewSSD1306(cfg.Hardware.Oled.Bus, cfg.Oled.Width, cfg.Oled.Height)
	if err != nil {
		return nil, fmt.Errorf("ssd1306: %w", err)
	}
	return &Display{
		device:          device,
		updatesDuration: cfg.Hardware.Oled.UpdatesDuration,
		updateNow:       make(chan struct{}),
	}, nil
}

// SetRenderer sets new renderer object for display.
func (d *Display) SetRenderer(renderer Renderer) {
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
func (d *Display) StartWithShutdown(ctx context.Context) error {
	// start display updates loop
	d.start(ctx)
	// close display after ctx is done
	if err := d.close(); err != nil {
		return err
	}
	return nil
}

// start starts display updates loop.
func (d *Display) start(ctx context.Context) {
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
			logrus.Errorf("Render object: %v", err)
		}
	}
}

// close clears display screen and closes display connection.
func (d *Display) close() error {
	if err := d.device.DisplayClear(); err != nil {
		logrus.Errorf("clear screen on display close: %v", err)
	}
	return d.device.Close()
}

// render renders renderer object to display.
func (d *Display) render() error {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if d.renderer == nil {
		return nil
	}
	return d.renderer.Render(d.device)
}
