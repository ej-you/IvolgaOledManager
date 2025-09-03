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
	// will be closed if the service was completely started and is ready-to-use now
	ready chan struct{}

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
		ready:           make(chan struct{}),
		device:          device,
		updatesDuration: updatesDuration,
		updateNow:       make(chan struct{}),
	}, nil
}

// DisplayTextLines displays given text lines.
// It is an adaptor for the same method
// of ssd1306.SSD1306 to implement Display interface.
func (s *Service) DisplayTextLines(lines ...TextLine) error {
	ssd1306Lines := make([]ssd1306.TextLine, 0, len(lines))

	for _, line := range lines {
		ssd1306Lines = append(ssd1306Lines, ssd1306.TextLine{
			Content:      line.Content,
			RelativeSize: line.RelativeSize,
		})
	}
	return s.device.DisplayTextLines(ssd1306Lines...)
}

// DisplayImage displays image from given path on the OLED-display.
func (s *Service) DisplayImage(path string, x, y int) error {
	return s.device.DisplayImage(path, x, y)
}

// SetRenderer sets new renderer object for display.
func (s *Service) SetRenderer(renderer Renderer) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.renderer = renderer

	// send force update signal to chan (if it is not filled)
	select {
	case s.updateNow <- struct{}{}:
	default:
	}
}

// StartWithShutdown starts display.
// Service will stop when given context will be stopped.
// This method is blocking.
func (s *Service) StartWithShutdown(ctx context.Context) error {
	logrus.Info("start display...")
	defer logrus.Info("stop display: ok")

	// start display updates loop
	s.start(ctx)
	// close display after ctx is done
	if err := s.close(); err != nil {
		return err
	}
	return nil
}

// Ready signals that the service is ready-to-use.
func (s *Service) Ready() <-chan struct{} {
	return s.ready
}

// start starts display updates loop.
func (s *Service) start(ctx context.Context) {
	// init ticker for display periodically updates
	ticker := time.NewTicker(s.updatesDuration)
	defer ticker.Stop()

	// notify that service is ready-to-use
	close(s.ready)

	var err error
	for {
		select {
		// periodically updates
		case <-ticker.C:
			err = s.render()
		// force update after setting new renderer
		case <-s.updateNow:
			err = s.render()
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
func (s *Service) render() error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.renderer == nil {
		return nil
	}
	return s.renderer.Render(s)
}

// close clears display screen and closes display connection.
func (s *Service) close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.device.DisplayClear(); err != nil {
		logrus.Errorf("clear screen on display close: %v", err)
	}
	return s.device.Close()
}
