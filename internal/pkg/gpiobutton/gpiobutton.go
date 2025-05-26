// Package gpiobutton provides gpio button initialization and setting up
// the button rising/pulling handlers. All buttons are with an external pull up resistor,
// so default button value is HIGH.
// It contains GPIOButton interface.
package gpiobutton

import (
	"context"
	"errors"
	"fmt"
	"time"

	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
)

var _ GPIOButton = (*gpioButton)(nil)

// Function for handle button rising/falling.
type HandlerFunc func()

type GPIOButton interface {
	HandleWithShutdown(ctx context.Context, risingHandler, fallingHandler HandlerFunc) error
}

type gpioButton struct {
	gpioPin           gpio.PinIO
	checkAliveTimeout time.Duration
	state             gpio.Level
}

func New(gpioName string, checkAliveTimeout time.Duration) (GPIOButton, error) {
	// get button by GPIO name
	gpioPin := gpioreg.ByName(gpioName)
	if gpioPin == nil {
		return nil, errors.New("find gpio button by name")
	}
	// set up input for button
	if err := gpioPin.In(gpio.PullNoChange, gpio.BothEdges); err != nil {
		return nil, fmt.Errorf("set up input for gpio button: %w", err)
	}

	return &gpioButton{
		gpioPin:           gpioPin,
		checkAliveTimeout: checkAliveTimeout,
		state:             gpio.High,
	}, nil
}

// HandleWithShutdown sets up handlers for button and
// gracefully shutdown GPIO button after context is done.
func (b *gpioButton) HandleWithShutdown(ctx context.Context,
	risingHandler, fallingHandler HandlerFunc) error {

	for {
		// check context is done
		select {
		case <-ctx.Done():
			return b.shutdown()
		default:
			if !b.edgeOccured() {
				continue
			}
		}
		// handle button falling if level is HIGH else rising
		if b.state {
			fallingHandler()
		} else {
			risingHandler()
		}
	}
}

// edgeOccured returns true if real edge is occured, not timeout wait.
func (b *gpioButton) edgeOccured() bool {
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
func (b *gpioButton) shutdown() error {
	if err := b.gpioPin.Halt(); err != nil {
		return fmt.Errorf("halt gpio button: %w", err)
	}
	return nil
}
