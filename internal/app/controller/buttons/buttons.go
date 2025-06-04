// Package button provides functions for set up gpio buttons on input.
package buttons

import (
	"context"
	"fmt"
	"sync"
	"time"

	"sschmc/internal/app/repo/db"
	"sschmc/internal/app/repo/storage"
	"sschmc/internal/pkg/errlog"
	"sschmc/internal/pkg/gpiobutton"
)

type Buttons struct {
	btnEsc    gpiobutton.GPIOButton
	btnUp     gpiobutton.GPIOButton
	btnDown   gpiobutton.GPIOButton
	btnEnt    gpiobutton.GPIOButton
	msgRepoDB db.MessageRepoDB
	store     *storage.RepoStorageManager
	render    chan<- struct{}
}

// NewButtons returns new Buttons struct pointer.
// First four params are names for GPIO buttons.
// The checkAliveTimeout param is a duration to see if button alive.
// The store param is an app key-value storage.
// The render param is a chan to send tasks for renderer to output data.
func New(btnEscName, btnUpName, btnDownName, btnEntName string, checkAliveTimeout time.Duration,
	dbStorage db.MessageRepoDB, store *storage.RepoStorageManager,
	render chan<- struct{}) (*Buttons, error) {

	var err error
	buttons := &Buttons{
		msgRepoDB: dbStorage,
		store:     store,
		render:    render,
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
	// Set greetings on startup
	b.screenGreetings()

	var wg sync.WaitGroup
	// start handlers for every button in a separate goroutines simultaneously
	wg.Add(4)
	go func() {
		defer wg.Done()
		err := b.btnEsc.HandleWithShutdown(ctx, b.BtnEscRisingHandler(), func() {})
		if err != nil {
			errlog.Print(err)
		}
	}()
	go func() {
		defer wg.Done()
		err := b.btnUp.HandleWithShutdown(ctx, b.BtnUpRisingHandler(), func() {})
		if err != nil {
			errlog.Print(err)
		}
	}()
	go func() {
		defer wg.Done()
		err := b.btnDown.HandleWithShutdown(ctx, b.BtnDownRisingHandler(), func() {})
		if err != nil {
			errlog.Print(err)
		}
	}()
	go func() {
		defer wg.Done()
		err := b.btnEnt.HandleWithShutdown(ctx, b.BtnEntRisingHandler(), func() {})
		if err != nil {
			errlog.Print(err)
		}
	}()

	// wait until all gpio are closed
	wg.Wait()
}
