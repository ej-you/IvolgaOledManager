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
	activeChanMap ActiveChanMap
	btns          button.Buttons
	templ         *template.SensorData
}

// NewSensWindSpeedScreen returns a new instance of SensWindSpeedScreen.
func NewSensWindSpeedScreen(activeChanMap ActiveChanMap, btns button.Buttons,
	storage pubsub.Storage) *SensWindSpeedScreen {

	templ := template.NewSensorData(
		string(SensWindSpeed),
		activeChanMap[SensWindSpeed],
		storage,
		repo.SensWindSpeedKey,
	)
	return &SensWindSpeedScreen{
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
		button.ButtonEsc:  s.btnEsc,
		button.ButtonUp:   s.btnUp,
		button.ButtonDown: s.btnDown,
	}

	s.templ.SetBtnHandlersReg(func() {
		s.btns.SetRisingHandlers(btnHandlers)
	})
}

// btnEsc represents an escape button handler for screen.
func (s *SensWindSpeedScreen) btnEsc() {
	s.activeChanMap[SensWindSpeed] <- false
	s.activeChanMap[MenuMain] <- true
}

// btnUp represents an up button handler for screen.
func (s *SensWindSpeedScreen) btnUp() {
	s.activeChanMap[SensWindSpeed] <- false
	s.activeChanMap[SensWindDir] <- true
}

// btnDown represents an down button handler for screen.
func (s *SensWindSpeedScreen) btnDown() {
	s.activeChanMap[SensWindSpeed] <- false
	s.activeChanMap[SensPress] <- true
}
