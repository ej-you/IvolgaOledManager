package screen

import (
	"context"

	"IvolgaOledManager/internal/app/entity"
	"IvolgaOledManager/internal/app/repo"
	"IvolgaOledManager/internal/app/service/button"
	"IvolgaOledManager/internal/app/service/screen/template"
	"IvolgaOledManager/internal/app/usecase"
	"IvolgaOledManager/internal/pkg/pubsub"
)

// newMenuSensorconfGetter returns menu getter func for sensorconf menu.
// Also it saves new sensorconf into storage with repo.SensorconfMenu key.
func newMenuSensorconfGetter(storage pubsub.Storage,
	sensorconfUC usecase.SensconfUsecase) template.MenuGetter {

	return func() (*entity.Menu, error) {
		data, err := sensorconfUC.Get()
		if err != nil {
			return nil, err
		}
		menu := sensorconfUC.ToMenu(data)
		storage.Publish(repo.MenuSensconf, menu)
		return menu, nil
	}
}

// MenuSensorconfScreen represents the main menu data screen.
type MenuSensorconfScreen struct {
	name          Name
	activeChanMap ActiveChanMap
	btns          button.Buttons
	templ         *template.Menu
	storage       pubsub.Storage
	sensorconfUC  usecase.SensconfUsecase
}

// NewMenuSensorconfScreen returns a new instance of MenuMainScreen.
func NewMenuSensorconfScreen(screenName Name, activeChanMap ActiveChanMap, btns button.Buttons,
	storage pubsub.Storage, sensorconfUC usecase.SensconfUsecase) *MenuSensorconfScreen {

	templ := template.NewMenu(
		string(screenName),
		activeChanMap[screenName],
		storage,
		newMenuSensorconfGetter(storage, sensorconfUC),
	)
	return &MenuSensorconfScreen{
		name:          screenName,
		activeChanMap: activeChanMap,
		btns:          btns,
		templ:         templ,
		storage:       storage,
		sensorconfUC:  sensorconfUC,
	}
}

// Ready returns true if service was completely started and is ready-to-use now.
func (s *MenuSensorconfScreen) Ready() <-chan struct{} {
	return s.templ.Ready()
}

// StartWithShutdown starts service and wait for context cancellation to shutdown it.
func (s *MenuSensorconfScreen) StartWithShutdown(ctx context.Context) error {
	return s.templ.StartWithShutdown(ctx)
}

// prepareBtnHandlers creates button handlers to apply them after the screen is active
func (s *MenuSensorconfScreen) prepareBtnHandlers() {
	btnHandlers := button.Handlers{
		button.Esc:  s.btnEsc,
		button.Up:   s.templ.BtnUpDefault,
		button.Down: s.templ.BtnDownDefault,
		button.Ent:  s.btnEnt,
	}

	s.templ.SetBtnHandlersReg(func() {
		s.btns.SetRisingHandlers(btnHandlers)
	})
}

// btnEsc represents an escape button handler for screen.
func (s *MenuSensorconfScreen) btnEsc() error {
	s.activeChanMap[s.name] <- false
	s.activeChanMap[MenuSens] <- true
	return nil
}

// btnEnt represents an enter button handler for screen.
func (s *MenuSensorconfScreen) btnEnt() error {
	// update sensorconf menu in storage (for new selected item record)
	data, err := s.templ.GetFromStorage()
	if err != nil {
		return err
	}
	s.storage.Publish(repo.MenuSensconf, data)

	s.activeChanMap[s.name] <- false
	s.activeChanMap[MenuSensorconfItem] <- true
	return nil
}
