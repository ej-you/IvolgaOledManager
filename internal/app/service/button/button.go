// Package button provides service to handle buttons pressings.
package button

import (
	"context"
	"fmt"
	"sync"

	"IvolgaOledManager/config"
	"IvolgaOledManager/internal/pkg/button"
)

const _buttonsAmount = 4 // amount of buttons

var (
	ButtonEsc  button.Name = "esc"  // escape button
	ButtonUp   button.Name = "up"   // up button
	ButtonDown button.Name = "down" // down button
	ButtonEnt  button.Name = "ent"  // enter button
)

// Buttons is a map of button services.
type Buttons map[button.Name]*button.GPIOButton

// New sets up all buttons and returns them as a map.
func New(cfg *config.Config) (Buttons, error) {
	var err error
	btns := make(Buttons, _buttonsAmount)

	btns[ButtonEsc], err = button.New(ButtonEsc,
		cfg.Buttons.Escape, cfg.Hardware.Buttons.CheckAliveTimeout)
	if err != nil {
		return nil, err
	}
	btns[ButtonUp], err = button.New(ButtonUp,
		cfg.Buttons.Up, cfg.Hardware.Buttons.CheckAliveTimeout)
	if err != nil {
		return nil, err
	}
	btns[ButtonDown], err = button.New(ButtonDown,
		cfg.Buttons.Down, cfg.Hardware.Buttons.CheckAliveTimeout)
	if err != nil {
		return nil, err
	}
	btns[ButtonEnt], err = button.New(ButtonEnt,
		cfg.Buttons.Enter, cfg.Hardware.Buttons.CheckAliveTimeout)
	if err != nil {
		return nil, err
	}
	return btns, nil
}

// StartWithShutdown starts all buttons to handle pressings.
// Given context is used for all buttons.
// This method is blocking.
func (b Buttons) StartWithShutdown(ctx context.Context) error {
	btnsCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	errChan := make(chan error, 1)

	// start all buttons
	for key, val := range b {
		wg.Add(1)
		go func(btnName button.Name, btn *button.GPIOButton) {
			defer wg.Done()
			if err := btn.StartWithShutdown(btnsCtx); err != nil {
				// use default to prevent goroutine blocking if errChan is filled
				select {
				case errChan <- fmt.Errorf("start %s btn: %w", btnName, err):
					cancel()
				default:
				}
			}
		}(key, val)
	}
	// wait for all buttons
	wg.Wait()

	select {
	// return error if was occurred
	case err := <-errChan:
		return err
	default:
		return nil
	}
}

// Ready signals that the service is ready-to-use.
func (b Buttons) Ready() <-chan struct{} {
	done := make(chan struct{})

	go func() {
		defer close(done)
		var wg sync.WaitGroup

		for _, val := range b {
			wg.Add(1)
			go func(btn *button.GPIOButton) {
				defer wg.Done()
				// wait for button
				<-btn.Ready()
			}(val)
		}
		// wait for all buttons
		wg.Wait()
	}()

	return done
}

// Handlers is a map of handlers for buttons.
type Handlers map[button.Name]button.HandlerFunc

// SetRisingHandlers sets new handlers for buttons' rising.
func (b Buttons) SetRisingHandlers(handlers Handlers) {
	var handler button.HandlerFunc
	var found bool
	// iterate buttons
	for btnName, btn := range b {
		// use default handler if handler for button is not specified
		if handler, found = handlers[btnName]; !found {
			handler = button.DefaultHandlerFunc
		}
		// set handler
		btn.SetRisingHandler(handler)
	}
}

// SetFallingHandlers sets new handlers for buttons' falling.
func (b Buttons) SetFallingHandlers(handlers Handlers) {
	var handler button.HandlerFunc
	var found bool
	// iterate buttons
	for btnName, btn := range b {
		// use default handler if handler for button is not specified
		if handler, found = handlers[btnName]; !found {
			handler = button.DefaultHandlerFunc
		}
		// set handler
		btn.SetRisingHandler(handler)
	}
}
