package screen

import (
	"context"

	"IvolgaOledManager/internal/app/repo"
	"IvolgaOledManager/internal/app/service/button"
	"IvolgaOledManager/internal/app/service/screen/template"
	"IvolgaOledManager/internal/pkg/pubsub"
)

// SensTempScreen represents a wind direction sensor data screen.
type SensWindDirScreen struct {
	activeChanMap ActiveChanMap
	btns          button.Buttons
	templ         *template.SensorData
}

// NewSensWindDirScreen returns a new instance of SensWindDirScreen.
func NewSensWindDirScreen(activeChanMap ActiveChanMap, btns button.Buttons,
	storage pubsub.Storage) *SensWindDirScreen {

	templ := template.NewSensorData(
		string(SensWindDir),
		activeChanMap[SensWindDir],
		storage,
		repo.SensWindDirKey,
	)
	return &SensWindDirScreen{
		activeChanMap: activeChanMap,
		btns:          btns,
		templ:         templ,
	}
}

// Ready returns true if service was completely started and is ready-to-use now.
func (s *SensWindDirScreen) Ready() <-chan struct{} {
	return s.templ.Ready()
}

// StartWithShutdown starts service and wait for context cancellation to shutdown it.
func (s *SensWindDirScreen) StartWithShutdown(ctx context.Context) error {
	return s.templ.StartWithShutdown(ctx)
}

// prepareBtnHandlers creates button handlers to apply them after the screen is active
func (s *SensWindDirScreen) prepareBtnHandlers() {
	btnHandlers := button.Handlers{
		button.ButtonEsc:  s.btnEsc,
		button.ButtonUp:   s.btnUp,
		button.ButtonDown: s.btnDown,
	}

	s.templ.SetBtnHandlersReg(func() {
		s.btns.SetRisingHandlers(btnHandlers)
	})
}

// btnEsc represents an escape button handler for screen.
func (s *SensWindDirScreen) btnEsc() {
	s.activeChanMap[SensWindDir] <- false
	s.activeChanMap[MenuMain] <- true
}

// btnUp represents an up button handler for screen.
func (s *SensWindDirScreen) btnUp() {
	s.activeChanMap[SensWindDir] <- false
	s.activeChanMap[SensTemp] <- true
}

// btnDown represents an down button handler for screen.
func (s *SensWindDirScreen) btnDown() {
	s.activeChanMap[SensWindDir] <- false
	s.activeChanMap[SensWindSpeed] <- true
}
