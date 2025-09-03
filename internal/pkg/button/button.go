// Package button provides gpio button initialization and setting up
// the button rising/pulling handlers. All buttons are with an external pull up resistor,
// so default button value is HIGH.
package button

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
)

// Ensure GPIO-button implements interface
var _ Button = (*GPIOButton)(nil)

// Button describes all methods for button to manage input data.
type Button interface {
	// SetRisingHandler sets new handler for button rising.
	SetRisingHandler(handler HandlerFunc)
	// SetFallingHandler sets new handler for button falling.
	SetFallingHandler(handler HandlerFunc)
}

// Name represents a name of a button.
type Name string

// HandlerFunc is a function to handle button rising/falling.
type HandlerFunc func() error

// DefaultHandlerFunc is an empty handler func.
var DefaultHandlerFunc HandlerFunc = func() error { return nil }

// GPIOButton represents a button based on GPIO.
type GPIOButton struct {
	// will be closed if the service was completely started and is ready-to-use now
	ready chan struct{}

	serviceName       Name
	gpioPin           gpio.PinIO
	checkAliveTimeout time.Duration
	state             gpio.Level

	mu             sync.Mutex
	risingHandler  HandlerFunc
	fallingHandler HandlerFunc
}

// New sets up new GPIO button and returns it. Rising and falling handlers are empty.
// Use SetRisingHandler/SetFallingHandler to set up handlers for button.
func New(serviceName Name, gpioName string, checkAliveTimeout time.Duration) (*GPIOButton, error) {
	// get button by GPIO name
	gpioPin := gpioreg.ByName(gpioName)
	if gpioPin == nil {
		return nil, fmt.Errorf("find button:%s gpio by name %s", serviceName, gpioName)
	}
	// set up input for button
	if err := gpioPin.In(gpio.PullNoChange, gpio.BothEdges); err != nil {
		return nil, fmt.Errorf("set up input for gpio button:%s %s: %w", serviceName, gpioName, err)
	}

	return &GPIOButton{
		ready:             make(chan struct{}),
		serviceName:       serviceName,
		gpioPin:           gpioPin,
		checkAliveTimeout: checkAliveTimeout,
		state:             gpio.High,
		risingHandler:     DefaultHandlerFunc,
		fallingHandler:    DefaultHandlerFunc,
	}, nil
}

// SetRisingHandler sets new handler for button rising.
func (b *GPIOButton) SetRisingHandler(handler HandlerFunc) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.risingHandler = handler
}

// SetFallingHandler sets new handler for button falling.
func (b *GPIOButton) SetFallingHandler(handler HandlerFunc) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.fallingHandler = handler
}

// StartWithShutdown starts GPIO button handling and
// gracefully shutdown it after context is done.
func (b *GPIOButton) StartWithShutdown(ctx context.Context) error {
	logrus.Infof("start %s...", b.serviceName)
	defer logrus.Infof("stop %s: ok", b.serviceName)

	// notify that service is ready-to-use
	close(b.ready)
	for {
		// check context is done
		select {
		case <-ctx.Done():
			return b.shutdown()
		default:
			if !b.edgeOccurred() {
				continue
			}
		}
		// handle button pressing
		if err := b.handle(); err != nil {
			logrus.Errorf("%s: %v", b.serviceName, err)
		}
	}
}

// Ready signals that the service is ready-to-use.
func (b *GPIOButton) Ready() <-chan struct{} {
	return b.ready
}

// handle handles button falling if level is HIGH else rising.
func (b *GPIOButton) handle() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	// run handler
	if b.state {
		return b.fallingHandler()
	} else {
		return b.risingHandler()
	}
}

// edgeOccurred returns true if real edge is occurred, not timeout wait.
func (b *GPIOButton) edgeOccurred() bool {
	// wait for rising up or falling down of button and
	// check the button every checkAliveTimeout duration to see if it alive
	b.gpioPin.WaitForEdge(b.checkAliveTimeout)

	// check that new state not equals to old state
	newState := b.gpioPin.Read()
	if newState == b.state {
		return false
	}
	// set new state for button
	b.state = newState
	return true
}

// shutdown stops gpio button.
func (b *GPIOButton) shutdown() error {
	if err := b.gpioPin.Halt(); err != nil {
		return fmt.Errorf("halt gpio button: %w", err)
	}
	return nil
}
