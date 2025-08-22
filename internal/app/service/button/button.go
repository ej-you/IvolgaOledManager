// Package button provides service to handle buttons pressings.
package button

import (
	"context"
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"

	"IvolgaOledManager/config"
	"IvolgaOledManager/internal/pkg/gpiobutton"
)

const _buttonsAmount = 4 // amount of buttons

var _ Button = (*gpiobutton.GPIOButton)(nil)

// One button from buttons service.
type Button interface {
	HandleWithShutdown(ctx context.Context) error
	SetRisingHandler(ctx context.Context, handler gpiobutton.HandlerFunc)
	SetFallingHandler(ctx context.Context, handler gpiobutton.HandlerFunc)
}

// Buttons service.
type Buttons map[string]Button

// NewButtons sets up all buttons and returns them as a map.
func NewButtons(cfg *config.Config) (*Buttons, error) {
	var err error
	btns := make(Buttons, _buttonsAmount)

	btns["up"], err = gpiobutton.New(cfg.Buttons.Up, cfg.Hardware.Buttons.CheckAliveTimeout)
	if err != nil {
		return nil, fmt.Errorf("up btn: %w", err)
	}
	btns["down"], err = gpiobutton.New(cfg.Buttons.Down, cfg.Hardware.Buttons.CheckAliveTimeout)
	if err != nil {
		return nil, fmt.Errorf("down btn: %w", err)
	}
	btns["esc"], err = gpiobutton.New(cfg.Buttons.Escape, cfg.Hardware.Buttons.CheckAliveTimeout)
	if err != nil {
		return nil, fmt.Errorf("esc btn: %w", err)
	}
	btns["ent"], err = gpiobutton.New(cfg.Buttons.Enter, cfg.Hardware.Buttons.CheckAliveTimeout)
	if err != nil {
		return nil, fmt.Errorf("ent btn: %w", err)
	}
	return &btns, nil
}

// StartWithShutdown starts all buttons to handle pressings.
// Given context is used for all buttons.
// This method is blocking.
func (b *Buttons) StartWithShutdown(ctx context.Context) error {
	btnsCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	errChan := make(chan error, 1)

	// start all buttons
	for key, val := range *b {
		logrus.Infof("start button:%s service...", key)

		wg.Add(1)
		go func(btnName string, btn Button) {
			defer wg.Done()
			if err := btn.HandleWithShutdown(btnsCtx); err != nil {
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
