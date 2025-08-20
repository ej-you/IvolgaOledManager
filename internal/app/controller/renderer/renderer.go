// Package renderer provides functions to set up renderer
// for periodically update SSD1306 OLED display.
package renderer

import (
	"context"
	"fmt"
	"log"
	"time"

	"IvolgaOledManager/internal/app/repo/storage"
	"IvolgaOledManager/internal/pkg/drawer"
	"IvolgaOledManager/internal/pkg/errlog"
)

const (
	_displayWidth  = 128 // oled display width in pixels
	_displayHeight = 64  // oled display height in pixels
)

type Renderer struct {
	device             *drawer.SSD1306
	greetingsImgPath   string
	menuUpdateDuration time.Duration
	store              *storage.RepoStorageManager
	needUpdate         <-chan struct{}
}

func New(bus, greetingsImgPath string, menuUpdateDuration time.Duration,
	store *storage.RepoStorageManager, needUpdate <-chan struct{}) (*Renderer, error) {

	ssd1306, err := drawer.NewDrawerSSD1306(bus, _displayWidth, _displayHeight)
	if err != nil {
		return nil, fmt.Errorf("connect to oled: %w", err)
	}

	return &Renderer{
		device:             ssd1306,
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
	var err error
	for {
		select {
		case <-ctx.Done():
			return
		case <-r.needUpdate:
			err = r.update()
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
	case r.store.App.IsStationResult():
		return r.station(r.store.StationResults.Get())
	default:
		log.Println("WARNING: no one render rule found")
	}
	return nil
}
