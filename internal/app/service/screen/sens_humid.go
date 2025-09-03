package screen

import (
	"context"

	"IvolgaOledManager/internal/app/repo"
	"IvolgaOledManager/internal/app/service/button"
	"IvolgaOledManager/internal/app/service/screen/template"
	"IvolgaOledManager/internal/pkg/pubsub"
)

// SensTempScreen represents a humidity sensor data screen.
type SensHumidScreen struct {
	name          Name
	activeChanMap ActiveChanMap
	btns          button.Buttons
	templ         *template.SensorData
}

// NewSensHumidScreen returns a new instance of SensHumidScreen.
func NewSensHumidScreen(screenName Name, activeChanMap ActiveChanMap, btns button.Buttons,
	storage pubsub.Storage) *SensHumidScreen {

	templ := template.NewSensorData(
		string(screenName),
		activeChanMap[screenName],
		storage,
		repo.SensHumidKey,
	)
	return &SensHumidScreen{
		name:          screenName,
		activeChanMap: activeChanMap,
		btns:          btns,
		templ:         templ,
	}
}

// Ready returns true if service was completely started and is ready-to-use now.
func (s *SensHumidScreen) Ready() <-chan struct{} {
	return s.templ.Ready()
}

// StartWithShutdown starts service and wait for context cancellation to shutdown it.
func (s *SensHumidScreen) StartWithShutdown(ctx context.Context) error {
	return s.templ.StartWithShutdown(ctx)
}

// prepareBtnHandlers creates button handlers to apply them after the screen is active
func (s *SensHumidScreen) prepareBtnHandlers() {
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
func (s *SensHumidScreen) btnEsc() error {
	s.activeChanMap[s.name] <- false
	s.activeChanMap[MenuSens] <- true
	return nil
}

// btnUp represents an up button handler for screen.
func (s *SensHumidScreen) btnUp() error {
	s.activeChanMap[s.name] <- false
	s.activeChanMap[SensPress] <- true
	return nil
}

// btnDown represents an down button handler for screen.
func (s *SensHumidScreen) btnDown() error {
	s.activeChanMap[s.name] <- false
	s.activeChanMap[SensTemp] <- true
	return nil
}
