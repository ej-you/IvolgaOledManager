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

// newMenuSensconfGetter returns menu getter func for sensorconf menu.
// Also it saves new sensorconf into storage with repo.SensorconfMenu key.
func newMenuSensconfGetter(storage pubsub.Storage,
	sensorconfUC usecase.SensconfUsecase) template.MenuGetter {

	return func() (*entity.Menu, error) {
		sensconfMenu, err := sensorconfUC.GetAsMenu()
		if err != nil {
			return nil, err
		}
		storage.Publish(repo.MenuSensconf, sensconfMenu)
		return sensconfMenu, nil
	}
}

// MenuSensconfScreen represents the main menu data screen.
type MenuSensconfScreen struct {
	name          Name
	activeChanMap ActiveChanMap
	btns          button.Buttons
	templ         *template.Menu
	storage       pubsub.Storage
	sensconfUC    usecase.SensconfUsecase
}

// NewMenuSensconfScreen returns a new instance of MenuMainScreen.
func NewMenuSensconfScreen(screenName Name, activeChanMap ActiveChanMap, btns button.Buttons,
	storage pubsub.Storage, sensconfUC usecase.SensconfUsecase) *MenuSensconfScreen {

	templ := template.NewMenu(
		string(screenName),
		activeChanMap[screenName],
		storage,
		newMenuSensconfGetter(storage, sensconfUC),
	)
	return &MenuSensconfScreen{
		name:          screenName,
		activeChanMap: activeChanMap,
		btns:          btns,
		templ:         templ,
		storage:       storage,
		sensconfUC:    sensconfUC,
	}
}

// Ready returns true if service was completely started and is ready-to-use now.
func (s *MenuSensconfScreen) Ready() <-chan struct{} {
	return s.templ.Ready()
}

// StartWithShutdown starts service and wait for context cancellation to shutdown it.
func (s *MenuSensconfScreen) StartWithShutdown(ctx context.Context) error {
	return s.templ.StartWithShutdown(ctx)
}

// prepareBtnHandlers creates button handlers to apply them after the screen is active
func (s *MenuSensconfScreen) prepareBtnHandlers() {
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
func (s *MenuSensconfScreen) btnEsc() error {
	s.activeChanMap[s.name] <- false
	s.activeChanMap[MenuSens] <- true
	return nil
}

// btnEnt represents an enter button handler for screen.
func (s *MenuSensconfScreen) btnEnt() error {
	s.activeChanMap[s.name] <- false
	s.activeChanMap[MenuSensconfItem] <- true
	return nil
}
