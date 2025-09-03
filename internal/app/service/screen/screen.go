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
	"IvolgaOledManager/internal/app/usecase"
	"IvolgaOledManager/internal/pkg/pubsub"
)

// Ensure specific screens implement interfaces.
var _ Screen = (*GreetingsScreen)(nil)
var _ Screen = (*MenuMainScreen)(nil)
var _ Screen = (*MenuSensScreen)(nil)
var _ Screen = (*MenuSensorconfScreen)(nil)
var _ Screen = (*SensTempScreen)(nil)
var _ Screen = (*SensHumidScreen)(nil)
var _ Screen = (*SensPressScreen)(nil)
var _ Screen = (*SensWindSpeedScreen)(nil)
var _ Screen = (*SensWindDirScreen)(nil)

// Name represents a screen name.
type Name string

var (
	ImgGreetings   Name = "screen:greetings"
	MenuMain       Name = "screen:menu:main"
	MenuSens       Name = "screen:menu:sensor"
	MenuSensorconf Name = "screen:menu:sensorconf"
	SensTemp       Name = "screen:sensordata:temperature"
	SensHumid      Name = "screen:sensordata:humidity"
	SensPress      Name = "screen:sensordata:pressure"
	SensWindSpeed  Name = "screen:sensordata:wind:speed"
	SensWindDir    Name = "screen:sensordata:wind:direction"
)

// ActiveChan represents an active chan for screen.
type ActiveChan chan bool

// ActiveChanMap is a map of screen active chans.
type ActiveChanMap map[Name]ActiveChan

// newActiveChanMap returns a new instance of ActiveChMap.
func newActiveChanMap() ActiveChanMap {
	return ActiveChanMap{
		ImgGreetings:   make(ActiveChan, 1),
		MenuMain:       make(ActiveChan, 1),
		MenuSens:       make(ActiveChan, 1),
		MenuSensorconf: make(ActiveChan, 1),
		SensTemp:       make(ActiveChan, 1),
		SensHumid:      make(ActiveChan, 1),
		SensPress:      make(ActiveChan, 1),
		SensWindSpeed:  make(ActiveChan, 1),
		SensWindDir:    make(ActiveChan, 1),
	}
}

// Screen describes a screen service.
type Screen interface {
	// StartWithShutdown starts service and wait for context cancellation to shutdown it.
	StartWithShutdown(ctx context.Context) error
	// Ready returns true if service was completely started and is ready-to-use now.
	Ready() <-chan struct{}

	// prepareBtnHandlers creates button handlers to apply them after the screen is active
	prepareBtnHandlers()
}

// Manager represents screen manager.
// It prepares all screen services and connects them with chans.
type Manager struct {
	screens       []Screen
	renderService *render.Render
}

// NewManager returns a new instance of ScreenManager.
func NewManager(cfg *config.Config, btns button.Buttons,
	renderService *render.Render, storage pubsub.Storage,
	sensorconfUC usecase.SensorconfUsecase) *Manager {

	// init screen active chans
	activeCh := newActiveChanMap()
	// active greetings screen by default
	activeCh[ImgGreetings] <- true

	// greetings screen
	greetings := NewGreetingsScreen(ImgGreetings, activeCh, btns, storage, cfg.App.GreetingsImgPath)
	// menus screen
	menuMain := NewMenuMainScreen(MenuMain, activeCh, btns, storage)
	menuSens := NewMenuSensScreen(MenuSens, activeCh, btns, storage)
	menuSensorconf := NewMenuSensorconfScreen(MenuSensorconf, activeCh, btns, storage, sensorconfUC)
	// sensor data screens
	sensTemp := NewSensTempScreen(SensTemp, activeCh, btns, storage)
	sensHumid := NewSensHumidScreen(SensHumid, activeCh, btns, storage)
	sensPress := NewSensPressScreen(SensPress, activeCh, btns, storage)
	sensWindSpeed := NewSensWindSpeedScreen(SensWindSpeed, activeCh, btns, storage)
	sensWindDir := NewSensWindDirScreen(SensWindDir, activeCh, btns, storage)

	screens := []Screen{
		greetings,
		menuMain, menuSens, menuSensorconf,
		sensTemp, sensHumid, sensPress, sensWindSpeed, sensWindDir,
	}
	// prepare button handlers for all screens
	for _, screen := range screens {
		screen.prepareBtnHandlers()
	}
	return &Manager{
		screens:       screens,
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
