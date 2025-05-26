// Package app provides App interface for run full application.
package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"periph.io/x/host/v3"

	"sschmc/config"
	"sschmc/internal/app/controller/buttons"
	"sschmc/internal/app/controller/renderer"
	"sschmc/internal/pkg/storage"
)

var _ App = (*app)(nil)

type App interface {
	Run() error
}

// App implementation.
type app struct {
	cfg   *config.Config
	store storage.Storage
}

// New returns App interface.
func New(cfg *config.Config) (App, error) {
	// initialise all relevant drivers
	if _, err := host.Init(); err != nil {
		return nil, fmt.Errorf("init drivers: %w", err)
	}
	return &app{
		cfg:   cfg,
		store: storage.NewMap(),
	}, nil
}

// Run starts full application.
func (a app) Run() error {
	// ctx for app
	appContext, appCancel := context.WithCancel(context.Background())
	// channel for display updates
	updateDisplay := make(chan struct{})

	// handle shutdown process signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	// create gracefully shutdown task
	go func() {
		defer appCancel()
		handledSignal := <-quit
		log.Printf("Get %q signal. Shutdown app...", handledSignal.String())
	}()

	// init display
	render, err := renderer.New(a.cfg.Hardware.Oled.Bus, a.cfg.App.GreetingsImgPath,
		updateDisplay, a.store)
	if err != nil {
		return fmt.Errorf("init display: %w", err)
	}

	// init buttons
	btns, err := buttons.New(
		a.cfg.Hardware.Buttons.Escape,
		a.cfg.Hardware.Buttons.Up,
		a.cfg.Hardware.Buttons.Down,
		a.cfg.Hardware.Buttons.Enter,
		time.Second,
		a.store,
		updateDisplay,
	)
	if err != nil {
		return err
	}

	// start renderer
	go render.StartWithShutdown(appContext)
	// start handle buttons rising/falling
	btns.HandleAll(appContext)
	log.Println("App shutdown successfully!")
	return nil
}
