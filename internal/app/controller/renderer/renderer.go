// Package renderer provides functions to set up renderer
// for periodically update SSD1306 OLED display.
package renderer

import (
	"context"
	"fmt"
	"log"
	"time"

	"sschmc/internal/app/repo/storage"
	"sschmc/internal/pkg/errlog"
	"sschmc/internal/pkg/ssd1306"
)

type Renderer struct {
	device             *ssd1306.SSD1306
	greetingsImgPath   string
	menuUpdateDuration time.Duration
	store              *storage.RepoStorageManager
	needUpdate         <-chan struct{}
}

func New(bus, greetingsImgPath string, menuUpdateDuration time.Duration,
	store *storage.RepoStorageManager, needUpdate <-chan struct{}) (*Renderer, error) {

	oled, err := ssd1306.NewSSD1306(bus)
	if err != nil {
		return nil, fmt.Errorf("connect to oled: %w", err)
	}

	return &Renderer{
		device:             oled,
		greetingsImgPath:   greetingsImgPath,
		needUpdate:         needUpdate,
		menuUpdateDuration: menuUpdateDuration,
		store:              store,
	}, nil
}

// StartWithShutdown starts renderer in background and wait for
// context is done for gracefully shutdown output device.
// This function is blocking.
func (r *Renderer) StartWithShutdown(ctx context.Context) {
	// start renderer loop
	startCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go r.start(startCtx)

	// if main context is done
	<-ctx.Done()
	if err := r.close(); err != nil {
		errlog.Print(err)
	}
}

func (r *Renderer) start(ctx context.Context) {
	// init ticker for menus periodically updates
	ticker := time.NewTicker(r.menuUpdateDuration)
	defer ticker.Stop()
	ticker.Stop()

	var err error
	for {
		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
			err = r.update()

		case <-r.needUpdate:
			err = r.update()
			// reset ticker for menus
			if r.store.App.IsMenuAny() {
				ticker.Reset(r.menuUpdateDuration)
			} else {
				ticker.Stop()
			}
		}
		if err != nil {
			errlog.Print(err)
		}
	}
}

// clear clears image.
func (r *Renderer) clear() error {
	return r.device.DisplayClear()
}

// close clears image and closes display connection.
func (r *Renderer) close() error {
	if err := r.clear(); err != nil {
		return err
	}
	return r.device.Close()
}

// update updates image according to app-status.
func (r *Renderer) update() error {
	switch {
	case r.store.App.IsNone():
		return r.clear()
	case r.store.App.IsGreetings():
		return r.greetings()
	case r.store.App.IsMenuMain():
		return r.menu(r.store.Menu.GetMain())
	// first main menu branch
	case r.store.App.IsMenuLogs():
		return r.menu(r.store.Menu.GetLogs())
	case r.store.App.IsMenuLevel():
		return r.menu(r.store.Menu.GetLevel())
	case r.store.App.IsMessage():
		return r.message(r.store.Message.Get())
	// second main menu branch
	case r.store.App.IsMenuStation():
		return r.menu(r.store.Menu.GetStation())
	case r.store.App.IsSensor():
		return r.sensor(r.store.Sensor.Get())
	default:
		log.Println("WARNING: no one render rule found")
	}
	return nil
}
