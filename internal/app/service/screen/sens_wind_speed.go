package screen

import (
	"context"

	"IvolgaOledManager/internal/app/repo"
	"IvolgaOledManager/internal/app/service/button"
	"IvolgaOledManager/internal/app/service/screen/template"
	"IvolgaOledManager/internal/pkg/pubsub"
)

// SensTempScreen represents a wind speed sensor data screen.
type SensWindSpeedScreen struct {
	name          Name
	activeChanMap ActiveChanMap
	btns          button.Buttons
	templ         *template.Sensdata
}

// NewSensWindSpeedScreen returns a new instance of SensWindSpeedScreen.
func NewSensWindSpeedScreen(screenName Name, activeChanMap ActiveChanMap, btns button.Buttons,
	storage pubsub.Storage) *SensWindSpeedScreen {

	templ := template.NewSensdata(
		string(screenName),
		activeChanMap[screenName],
		storage,
		repo.SensWindSpeedKey,
	)
	return &SensWindSpeedScreen{
		name:          screenName,
		activeChanMap: activeChanMap,
		btns:          btns,
		templ:         templ,
	}
}

// Ready returns true if service was completely started and is ready-to-use now.
func (s *SensWindSpeedScreen) Ready() <-chan struct{} {
	return s.templ.Ready()
}

// StartWithShutdown starts service and wait for context cancellation to shutdown it.
func (s *SensWindSpeedScreen) StartWithShutdown(ctx context.Context) error {
	return s.templ.StartWithShutdown(ctx)
}

// prepareBtnHandlers creates button handlers to apply them after the screen is active
func (s *SensWindSpeedScreen) prepareBtnHandlers() {
	btnHandlers := button.Handlers{
		button.Esc:  s.btnEsc,
		button.Up:   s.btnUp,
		button.Down: s.btnDown,
	}

	s.templ.SetBtnHandlersReg(func() {
		s.btns.SetRisingHandlers(btnHandlers)
	})
}

// btnEsc represents an escape button handler for screen.
func (s *SensWindSpeedScreen) btnEsc() error {
	s.activeChanMap[s.name] <- false
	s.activeChanMap[MenuSens] <- true
	return nil
}

// btnUp represents an up button handler for screen.
func (s *SensWindSpeedScreen) btnUp() error {
	s.activeChanMap[s.name] <- false
	s.activeChanMap[SensWindDir] <- true
	return nil
}

// btnDown represents an down button handler for screen.
func (s *SensWindSpeedScreen) btnDown() error {
	s.activeChanMap[s.name] <- false
	s.activeChanMap[SensPress] <- true
	return nil
}
