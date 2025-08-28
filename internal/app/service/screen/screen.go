package screen

import (
	"context"
	"fmt"
	"sync"

	"IvolgaOledManager/config"
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
	screens []Screen
}

// NewManager returns a new instance of ScreenManager.
func NewManager(cfg *config.Config, storage pubsub.Storage) *Manager {
	// init screen active chans
	greetCh := make(chan bool, 1)
	sensorTempCh := make(chan bool, 1)
	// active greetings screen by default
	greetCh <- true

	greetings := NewGreetings(greetCh, storage, cfg.App.GreetingsImgPath)
	sensorTemp := NewTemperature(sensorTempCh, storage)
	return &Manager{
		screens: []Screen{greetings, sensorTemp},
	}
}

// StartWithShutdown starts all screens. This method is blocking.
// Cancellation of given context may be used to shutdown service.
func (m *Manager) StartWithShutdown(ctx context.Context) error {
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
