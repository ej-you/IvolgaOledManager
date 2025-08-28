package screen

import (
	"context"
	"fmt"
	"sync"

	"IvolgaOledManager/config"
	buttonservice "IvolgaOledManager/internal/app/service/button"
	"IvolgaOledManager/internal/app/service/render"
	"IvolgaOledManager/internal/pkg/button"
	"IvolgaOledManager/internal/pkg/pubsub"
)

// Ensure specific screens implements interface.
var _ Screen = (*Greetings)(nil)

type Screen interface {
	// StartWithShutdown starts service and wait for context cancellation to shutdown it.
	StartWithShutdown(ctx context.Context) error
	// Ready returns true if service was completely started and is ready-to-use now.
	Ready() <-chan struct{}
}

// Manager prepare all screen services and connect them with chans.
type Manager struct {
	screens       []Screen
	renderService *render.Render
}

// NewManager returns a new instance of ScreenManager.
func NewManager(cfg *config.Config, btns *buttonservice.Buttons,
	renderService *render.Render, storage pubsub.Storage) *Manager {

	// init screen active chans
	greetCh := make(chan bool, 1)
	sensorTempCh := make(chan bool, 1)
	// active greetings screen by default
	greetCh <- true

	// greetings screen
	greetingsReg := getBtnHandlersRegFunc(btns, buttonservice.Handlers{
		button.ButtonEnt: func(_ context.Context) { greetCh <- false; sensorTempCh <- true }})
	greetings := NewGreetings(greetCh, greetingsReg, storage, cfg.App.GreetingsImgPath)

	// temperature sensor data screen
	sensorTempReg := getBtnHandlersRegFunc(btns, buttonservice.Handlers{
		button.ButtonEsc: func(_ context.Context) { sensorTempCh <- false; greetCh <- true }})
	sensorTemp := NewTemperature(sensorTempCh, sensorTempReg, storage)

	return &Manager{
		screens:       []Screen{greetings, sensorTemp},
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
func getBtnHandlersRegFunc(btns *buttonservice.Buttons,
	handlers buttonservice.Handlers) BtnHandlersRegFunc {

	return func() {
		var handler button.HandlerFunc
		var found bool
		// iterate buttons
		for btnName, btn := range *btns {
			// use default handler if handler for button is not specified
			if handler, found = handlers[btnName]; !found {
				handler = button.DefaultHandlerFunc
			}
			// set handler
			btn.SetRisingHandler(context.Background(), handler)
		}
	}
}
