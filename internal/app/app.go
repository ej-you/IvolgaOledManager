// Package app provides App interface for run full application.
package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"syscall"

	"periph.io/x/host/v3"

	"sschmc/config"
	"sschmc/internal/app/controller/buttons"
	"sschmc/internal/app/controller/renderer"
	"sschmc/internal/pkg/errlog"
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
	// init display
	display, err := renderer.NewDisplay(a.cfg.Hardware.Oled.Bus, a.cfg.App.GreetingsImgPath)
	if err != nil {
		return fmt.Errorf("init display: %w", err)
	}
	// defer display closing
	defer func() {
		if err := display.Close(); err != nil {
			errlog.Print(err)
		}
	}()

	// if err := display.Menu(); err != nil {
	// 	return fmt.Errorf("display menu: %w", err)
	// }
	// // if err := oled.Greetings(); err != nil {
	// // 	return fmt.Errorf("display greetings: %w", err)
	// // }
	// time.Sleep(10 * time.Second)
	// if err := display.Clear(); err != nil {
	// 	return fmt.Errorf("clear display: %w", err)
	// }
	// return nil

	// handle shutdown process signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	// ctx for all buttons
	buttonsContext, buttonsCancel := context.WithCancel(context.Background())
	// create gracefully shutdown task
	go func() {
		defer buttonsCancel()
		handledSignal := <-quit
		log.Printf("Get %q signal. Shutdown app...", handledSignal.String())
	}()

	// init buttons
	btns, err := buttons.New(
		a.cfg.Hardware.Buttons.Escape,
		a.cfg.Hardware.Buttons.Up,
		a.cfg.Hardware.Buttons.Down,
		a.cfg.Hardware.Buttons.Enter,
		time.Second,
		a.store,
		display,
	)
	if err != nil {
		return err
	}
	// start handle buttons rising/falling
	btns.HandleAll(buttonsContext)
	log.Println("App shutdown successfully!")
	return nil
}
