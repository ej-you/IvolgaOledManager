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
	name          Name
	activeChanMap ActiveChanMap
	btns          button.Buttons
	templ         *template.SensorData
}

// NewSensPressScreen returns a new instance of SensPressScreen.
func NewSensPressScreen(screenName Name, activeChanMap ActiveChanMap, btns button.Buttons,
	storage pubsub.Storage) *SensPressScreen {

	templ := template.NewSensorData(
		string(screenName),
		activeChanMap[screenName],
		storage,
		repo.SensPressKey,
	)
	return &SensPressScreen{
		name:          screenName,
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
		button.Esc:  s.btnEsc,
		button.Up:   s.btnUp,
		button.Down: s.btnDown,
	}

	s.templ.SetBtnHandlersReg(func() {
		s.btns.SetRisingHandlers(btnHandlers)
	})
}

// btnEsc represents an escape button handler for screen.
func (s *SensPressScreen) btnEsc() error {
	s.activeChanMap[s.name] <- false
	s.activeChanMap[MenuSens] <- true
	return nil
}

// btnUp represents an up button handler for screen.
func (s *SensPressScreen) btnUp() error {
	s.activeChanMap[s.name] <- false
	s.activeChanMap[SensWindSpeed] <- true
	return nil
}

// btnDown represents an down button handler for screen.
func (s *SensPressScreen) btnDown() error {
	s.activeChanMap[s.name] <- false
	s.activeChanMap[SensHumid] <- true
	return nil
}
