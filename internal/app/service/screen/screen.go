// Package screen provides screen services for serving
// specific screens and navigating between them.
package screen

import (
	"context"
	"fmt"
	"sync"

	"IvolgaOledManager/config"
	"IvolgaOledManager/internal/app/service/button"
	"IvolgaOledManager/internal/app/service/render"
	"IvolgaOledManager/internal/app/service/screen/template"
	"IvolgaOledManager/internal/pkg/pubsub"
)

// Ensure specific screen templates implements interface.
var _ Screen = (*template.Image)(nil)
var _ Screen = (*template.SensorData)(nil)

// Screen describes a screen service.
type Screen interface {
	// StartWithShutdown starts service and wait for context cancellation to shutdown it.
	StartWithShutdown(ctx context.Context) error
	// Ready returns true if service was completely started and is ready-to-use now.
	Ready() <-chan struct{}
}

// Manager represents screen manager.
// It prepares all screen services and connects them with chans.
type Manager struct {
	screens       []Screen
	renderService *render.Render
}

// NewManager returns a new instance of ScreenManager.
func NewManager(cfg *config.Config, btns button.Buttons,
	renderService *render.Render, storage pubsub.Storage) *Manager {

	// init screen active chans
	greetCh := make(chan bool, 1)
	sensTempCh := make(chan bool, 1)
	sensHumidCh := make(chan bool, 1)
	sensPressCh := make(chan bool, 1)
	sensWindSpeedCh := make(chan bool, 1)
	sensWindDirCh := make(chan bool, 1)
	// active greetings screen by default
	greetCh <- true

	// greetings screen
	greetingsReg := getBtnHandlersRegFunc(btns, button.Handlers{
		button.ButtonEnt: func() { greetCh <- false; sensTempCh <- true }})
	greetings := NewGreetings(greetCh, greetingsReg, storage, cfg.App.GreetingsImgPath)
	// temperature sensor data screen
	sensTempReg := getBtnHandlersRegFunc(btns, button.Handlers{
		button.ButtonEsc:  func() { sensTempCh <- false; greetCh <- true },
		button.ButtonUp:   func() { sensTempCh <- false; sensHumidCh <- true },
		button.ButtonDown: func() { sensTempCh <- false; sensWindDirCh <- true }})
	sensTemp := NewSensTemp(sensTempCh, sensTempReg, storage)
	// humidity sensor data screen
	sensHumidReg := getBtnHandlersRegFunc(btns, button.Handlers{
		button.ButtonEsc:  func() { sensHumidCh <- false; greetCh <- true },
		button.ButtonUp:   func() { sensHumidCh <- false; sensPressCh <- true },
		button.ButtonDown: func() { sensHumidCh <- false; sensTempCh <- true }})
	sensHumid := NewSensHumid(sensHumidCh, sensHumidReg, storage)
	// pressure sensor data screen
	sensPressReg := getBtnHandlersRegFunc(btns, button.Handlers{
		button.ButtonEsc:  func() { sensPressCh <- false; greetCh <- true },
		button.ButtonUp:   func() { sensPressCh <- false; sensWindSpeedCh <- true },
		button.ButtonDown: func() { sensPressCh <- false; sensHumidCh <- true }})
	sensPress := NewSensPress(sensPressCh, sensPressReg, storage)
	// wind speed sensor data screen
	sensWindSpeedReg := getBtnHandlersRegFunc(btns, button.Handlers{
		button.ButtonEsc:  func() { sensWindSpeedCh <- false; greetCh <- true },
		button.ButtonUp:   func() { sensWindSpeedCh <- false; sensWindDirCh <- true },
		button.ButtonDown: func() { sensWindSpeedCh <- false; sensPressCh <- true }})
	sensWindSpeed := NewSensWindSpeed(sensWindSpeedCh, sensWindSpeedReg, storage)
	// wind direction sensor data screen
	sensWindDirReg := getBtnHandlersRegFunc(btns, button.Handlers{
		button.ButtonEsc:  func() { sensWindDirCh <- false; greetCh <- true },
		button.ButtonUp:   func() { sensWindDirCh <- false; sensTempCh <- true },
		button.ButtonDown: func() { sensWindDirCh <- false; sensWindSpeedCh <- true }})
	sensWindDir := NewSensWindDir(sensWindDirCh, sensWindDirReg, storage)

	return &Manager{
		screens: []Screen{greetings,
			sensTemp, sensHumid, sensPress, sensWindSpeed, sensWindDir},
		renderService: renderService,
	}
}

// StartWithShutdown starts all screens. This method is blocking.
// Cancellation of given context may be used to shutdown service.
func (m *Manager) StartWithShutdown(ctx context.Context) error {
	// wait for the start of the render service
	<-m.renderService.Ready()

	screensCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	errChan := make(chan error, 1)

	// start all screens
	for _, screen := range m.screens {
		wg.Add(1)
		go func(screen Screen) {
			defer wg.Done()
			if err := screen.StartWithShutdown(screensCtx); err != nil {
				// use default to prevent goroutine blocking if errChan is filled
				select {
				case errChan <- fmt.Errorf("start screen: %w", err):
					cancel()
				default:
				}
			}
		}(screen)
	}
	// wait for all screens
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
// ScreenManager is ready when all screens are ready.
func (m *Manager) Ready() <-chan struct{} {
	done := make(chan struct{})

	go func() {
		var wg sync.WaitGroup

		for _, val := range m.screens {
			wg.Add(1)
			go func(screen Screen) {
				defer wg.Done()
				// wait for screen
				<-screen.Ready()
			}(val)
		}
		wg.Wait()
		// all screens are ready
		close(done)
	}()

	return done
}

// BtnHandlersRegFunc is a func to register button handlers for
// connecting screen services and navigating between screens
type BtnHandlersRegFunc func()

// getBtnHandlersRegFunc returns a button handlers register func for given buttons' handlers.
func getBtnHandlersRegFunc(btns button.Buttons, handlers button.Handlers) BtnHandlersRegFunc {
	return func() {
		btns.SetRisingHandlers(handlers)
	}
}
