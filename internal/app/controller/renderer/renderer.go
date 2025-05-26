// Package renderer provides functions to set up renderer
// for periodically update SSD1306 OLED display.
package renderer

import (
	"context"
	"fmt"
	"log"
	"time"

	"sschmc/internal/app/constants"
	"sschmc/internal/pkg/errlog"
	"sschmc/internal/pkg/ssd1306"
	"sschmc/internal/pkg/storage"
)

type Renderer struct {
	device           *ssd1306.SSD1306
	greetingsImgPath string
	needUpdate       <-chan struct{}
	updateDuration   time.Duration
	store            storage.Storage
}

func New(bus, greetingsImgPath string,
	needUpdate <-chan struct{}, store storage.Storage) (*Renderer, error) {

	oled, err := ssd1306.NewSSD1306(bus)
	if err != nil {
		return nil, fmt.Errorf("connect to oled: %w", err)
	}

	return &Renderer{
		device:           oled,
		greetingsImgPath: greetingsImgPath,
		needUpdate:       needUpdate,
		updateDuration:   time.Second,
		store:            store,
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
	var (
		err       error
		appStatus string
	)
	ticker := time.NewTicker(r.updateDuration)

	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			return
		case <-ticker.C:
			fmt.Println("Ticker")
			appStatus = r.store.Get(constants.KeyAppStatus)
			err = r.update(appStatus)
		case <-r.needUpdate:
			fmt.Println("Update")
			appStatus = r.store.Get(constants.KeyAppStatus)
			err = r.update(appStatus)
			ticker.Reset(r.updateDuration)
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
func (r *Renderer) update(appStatus string) error {
	fmt.Println("update render")
	switch appStatus {
	case constants.ValueAppStatusNone:
		log.Println("*** clear rendered ***")
		return r.clear()
	case constants.ValueAppStatusGreetings:
		log.Println("*** render greetings ***")
		return r.greetings()
	case constants.ValueAppStatusMenuMain:
		log.Println("*** render main menu ***")
		return r.menu()
	default:
		log.Println("*** no one render rule found ***")
	}

	return nil
}
