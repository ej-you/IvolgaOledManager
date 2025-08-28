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
	SetRisingHandler(ctx context.Context, handler HandlerFunc)
	// SetFallingHandler sets new handler for button falling.
	SetFallingHandler(ctx context.Context, handler HandlerFunc)
}

// Name represents a name of a button.
type Name string

var (
	ButtonEsc  Name = "esc"  // escape button
	ButtonUp   Name = "up"   // up button
	ButtonDown Name = "down" // down button
	ButtonEnt  Name = "ent"  // enter button
)

// HandlerFunc is a function to handle button rising/falling.
type HandlerFunc func(ctx context.Context)

// DefaultHandlerFunc is an empty handler func.
var DefaultHandlerFunc HandlerFunc = func(_ context.Context) {}

// GPIOButton represents a button based on GPIO.
type GPIOButton struct {
	// will be closed if the service was completely started and is ready-to-use now
	ready chan struct{}

	Name              Name
	gpioPin           gpio.PinIO
	checkAliveTimeout time.Duration
	state             gpio.Level

	mu             sync.Mutex
	handlersCtx    context.Context
	risingHandler  HandlerFunc
	fallingHandler HandlerFunc
}

// New sets up new GPIO button and returns it. Rising and falling handlers are empty.
// Use SetRisingHandler/SetFallingHandler to set up handlers for button.
func New(name Name, gpioName string, checkAliveTimeout time.Duration) (*GPIOButton, error) {
	// get button by GPIO name
	gpioPin := gpioreg.ByName(gpioName)
	if gpioPin == nil {
		return nil, fmt.Errorf("find button:%s gpio by name %s", name, gpioName)
	}
	// set up input for button
	if err := gpioPin.In(gpio.PullNoChange, gpio.BothEdges); err != nil {
		return nil, fmt.Errorf("set up input for gpio button:%s %s: %w", name, gpioName, err)
	}

	return &GPIOButton{
		ready:             make(chan struct{}),
		Name:              name,
		gpioPin:           gpioPin,
		checkAliveTimeout: checkAliveTimeout,
		state:             gpio.High,
		handlersCtx:       context.Background(),
		risingHandler:     DefaultHandlerFunc,
		fallingHandler:    DefaultHandlerFunc,
	}, nil
}

// SetRisingHandler sets new handler for button rising.
func (b *GPIOButton) SetRisingHandler(ctx context.Context, handler HandlerFunc) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.handlersCtx = ctx
	b.risingHandler = handler
}

// SetFallingHandler sets new handler for button falling.
func (b *GPIOButton) SetFallingHandler(ctx context.Context, handler HandlerFunc) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.handlersCtx = ctx
	b.fallingHandler = handler
}

// StartWithShutdown starts GPIO button handling and
// gracefully shutdown it after context is done.
func (b *GPIOButton) StartWithShutdown(ctx context.Context) error {
	logrus.Infof("start button:%s service...", b.Name)

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
		b.handle()
	}
}

// Ready signals that the service is ready-to-use.
func (b *GPIOButton) Ready() <-chan struct{} {
	return b.ready
}

// handle handles button falling if level is HIGH else rising.
func (b *GPIOButton) handle() {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.state {
		b.fallingHandler(b.handlersCtx)
	} else {
		b.risingHandler(b.handlersCtx)
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
