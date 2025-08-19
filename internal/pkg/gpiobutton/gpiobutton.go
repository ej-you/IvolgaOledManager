// Package gpiobutton provides gpio button initialization and setting up
// the button rising/pulling handlers. All buttons are with an external pull up resistor,
// so default button value is HIGH.
package gpiobutton

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
)

// Function to handle button rising/falling.
type HandlerFunc func(ctx context.Context)

type GPIOButton struct {
	gpioPin           gpio.PinIO
	checkAliveTimeout time.Duration
	state             gpio.Level

	handlersCtx    context.Context
	risingHandler  HandlerFunc
	fallingHandler HandlerFunc
}

// New sets up new gpio button and returns it. Rising and falling handlers are empty.
// Use SetRisingHandler/SetFallingHandler to set up handlers for button.
func New(gpioName string, checkAliveTimeout time.Duration) (*GPIOButton, error) {
	// get button by GPIO name
	gpioPin := gpioreg.ByName(gpioName)
	if gpioPin == nil {
		return nil, fmt.Errorf("find gpio button by name %s", gpioName)
	}
	// set up input for button
	if err := gpioPin.In(gpio.PullNoChange, gpio.BothEdges); err != nil {
		return nil, fmt.Errorf("set up input for gpio button %s: %w", gpioName, err)
	}

	logrus.Infof("Button %s was initialized successfully", gpioName)
	return &GPIOButton{
		gpioPin:           gpioPin,
		checkAliveTimeout: checkAliveTimeout,
		state:             gpio.High,
		handlersCtx:       context.Background(),
		risingHandler:     func(_ context.Context) {},
		fallingHandler:    func(_ context.Context) {},
	}, nil
}

// SetRisingHandler sets new handler when the button is rise.
func (b *GPIOButton) SetRisingHandler(ctx context.Context, handler HandlerFunc) {
	b.handlersCtx = ctx
	b.risingHandler = handler
}

// SetFallingHandler sets new handler when the button is fall.
func (b *GPIOButton) SetFallingHandler(ctx context.Context, handler HandlerFunc) {
	b.handlersCtx = ctx
	b.fallingHandler = handler
}

// HandleWithShutdown sets up handlers for button and
// gracefully shutdown GPIO button after context is done.
func (b *GPIOButton) HandleWithShutdown(ctx context.Context) error {
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
		// handle button falling if level is HIGH else rising
		if b.state {
			b.fallingHandler(b.handlersCtx)
		} else {
			b.risingHandler(b.handlersCtx)
		}
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

// shutdown stopped gpio button.
func (b *GPIOButton) shutdown() error {
	if err := b.gpioPin.Halt(); err != nil {
		return fmt.Errorf("halt gpio button: %w", err)
	}
	return nil
}
