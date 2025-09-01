package screen

import (
	"context"

	"IvolgaOledManager/internal/app/repo"
	"IvolgaOledManager/internal/app/service/button"
	"IvolgaOledManager/internal/app/service/screen/template"
	"IvolgaOledManager/internal/pkg/pubsub"
)

// SensTempScreen represents a pressure sensor data screen.
type SensPressScreen struct {
	activeChanMap ActiveChanMap
	btns          button.Buttons
	templ         *template.SensorData
}

// NewSensPressScreen returns a new instance of SensPressScreen.
func NewSensPressScreen(activeChanMap ActiveChanMap, btns button.Buttons,
	storage pubsub.Storage) *SensPressScreen {

	templ := template.NewSensorData(
		string(SensPress),
		activeChanMap[SensPress],
		storage,
		repo.SensPressKey,
	)
	return &SensPressScreen{
		activeChanMap: activeChanMap,
		btns:          btns,
		templ:         templ,
	}
}

// Ready returns true if service was completely started and is ready-to-use now.
func (s *SensPressScreen) Ready() <-chan struct{} {
	return s.templ.Ready()
}

// StartWithShutdown starts service and wait for context cancellation to shutdown it.
func (s *SensPressScreen) StartWithShutdown(ctx context.Context) error {
	return s.templ.StartWithShutdown(ctx)
}

// prepareBtnHandlers creates button handlers to apply them after the screen is active
func (s *SensPressScreen) prepareBtnHandlers() {
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
func (s *SensPressScreen) btnEsc() {
	s.activeChanMap[SensPress] <- false
	s.activeChanMap[MenuMain] <- true
}

// btnUp represents an up button handler for screen.
func (s *SensPressScreen) btnUp() {
	s.activeChanMap[SensPress] <- false
	s.activeChanMap[SensWindSpeed] <- true
}

// btnDown represents an down button handler for screen.
func (s *SensPressScreen) btnDown() {
	s.activeChanMap[SensPress] <- false
	s.activeChanMap[SensHumid] <- true
}
