package screen

import (
	"context"

	"IvolgaOledManager/internal/app/repo"
	"IvolgaOledManager/internal/app/service/button"
	"IvolgaOledManager/internal/app/service/screen/template"
	"IvolgaOledManager/internal/pkg/pubsub"
)

// SensTempScreen represents a temperature sensor data screen.
type SensTempScreen struct {
	activeChanMap ActiveChanMap
	btns          button.Buttons
	templ         *template.SensorData
}

// NewSensTempScreen returns a new instance of SensTempScreen.
func NewSensTempScreen(activeChanMap ActiveChanMap, btns button.Buttons,
	storage pubsub.Storage) *SensTempScreen {

	templ := template.NewSensorData(
		string(SensTemp),
		activeChanMap[SensTemp],
		storage,
		repo.SensTempKey,
	)
	return &SensTempScreen{
		activeChanMap: activeChanMap,
		btns:          btns,
		templ:         templ,
	}
}

// Ready returns true if service was completely started and is ready-to-use now.
func (s *SensTempScreen) Ready() <-chan struct{} {
	return s.templ.Ready()
}

// StartWithShutdown starts service and wait for context cancellation to shutdown it.
func (s *SensTempScreen) StartWithShutdown(ctx context.Context) error {
	return s.templ.StartWithShutdown(ctx)
}

// prepareBtnHandlers creates button handlers to apply them after the screen is active
func (s *SensTempScreen) prepareBtnHandlers() {
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
func (s *SensTempScreen) btnEsc() {
	s.activeChanMap[SensTemp] <- false
	s.activeChanMap[MenuMain] <- true
}

// btnUp represents an up button handler for screen.
func (s *SensTempScreen) btnUp() {
	s.activeChanMap[SensTemp] <- false
	s.activeChanMap[SensHumid] <- true
}

// btnDown represents an down button handler for screen.
func (s *SensTempScreen) btnDown() {
	s.activeChanMap[SensTemp] <- false
	s.activeChanMap[SensWindDir] <- true
}
