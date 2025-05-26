// Package button provides functions for set up gpio buttons on input.
package buttons

import (
	"context"
	"fmt"
	"sync"
	"time"

	"sschmc/internal/pkg/gpiobutton"
	"sschmc/internal/pkg/storage"
)

type Buttons struct {
	store   storage.Storage
	render  chan<- struct{}
	btnEsc  gpiobutton.GPIOButton
	btnUp   gpiobutton.GPIOButton
	btnDown gpiobutton.GPIOButton
	btnEnt  gpiobutton.GPIOButton
}

// NewButtons returns new Buttons struct pointer.
// First four params are names for GPIO buttons.
// The checkAliveTimeout param is a duration to see if button alive.
// The store param is an app key-value storage.
// The render param is a renderer for output data.
func New(btnEscName, btnUpName, btnDownName, btnEntName string, checkAliveTimeout time.Duration,
	store storage.Storage, render chan<- struct{}) (*Buttons, error) {

	var err error
	buttons := &Buttons{
		store:  store,
		render: render,
	}

	// init all buttons
	buttons.btnEsc, err = gpiobutton.New(btnEscName, checkAliveTimeout)
	if err != nil {
		return nil, fmt.Errorf("init escape button: %w", err)
	}
	buttons.btnUp, err = gpiobutton.New(btnUpName, checkAliveTimeout)
	if err != nil {
		return nil, fmt.Errorf("init up button: %w", err)
	}
	buttons.btnDown, err = gpiobutton.New(btnDownName, checkAliveTimeout)
	if err != nil {
		return nil, fmt.Errorf("init down button: %w", err)
	}
	buttons.btnEnt, err = gpiobutton.New(btnEntName, checkAliveTimeout)
	if err != nil {
		return nil, fmt.Errorf("init enter button: %w", err)
	}

	return buttons, nil
}

// HandleAll sets up all button handlers.
// Given context is used for all buttons.
// This function is blocking.
func (b *Buttons) HandleAll(ctx context.Context) {
	var wg sync.WaitGroup

	// start handlers for every button in a separate goroutines simultaneously
	wg.Add(4)
	go func() {
		defer wg.Done()
		b.btnEsc.HandleWithShutdown(ctx, b.BtnEscRisingHandler(), func() {})
	}()
	go func() {
		defer wg.Done()
		b.btnUp.HandleWithShutdown(ctx, b.BtnUpRisingHandler(), func() {})
	}()
	go func() {
		defer wg.Done()
		b.btnDown.HandleWithShutdown(ctx, b.BtnDownRisingHandler(), func() {})
	}()
	go func() {
		defer wg.Done()
		b.btnEnt.HandleWithShutdown(ctx, b.BtnEntRisingHandler(), func() {})
	}()

	// wait until all gpio are closed
	wg.Wait()
}
