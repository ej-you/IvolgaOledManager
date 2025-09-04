package screen

import (
	"context"

	"IvolgaOledManager/internal/app/entity"
	"IvolgaOledManager/internal/app/service/button"
	"IvolgaOledManager/internal/app/service/screen/template"
	"IvolgaOledManager/internal/pkg/pubsub"
)

// newMenuSensGetter returns menu getter func for sensor menu.
func newMenuSensGetter() template.MenuGetter {
	menu := &entity.Menu{
		Title: "Датчики",
		Items: []*entity.MenuItem{
			entity.NewMenuItem(context.Background(), "Данные"),
			entity.NewMenuItem(context.Background(), "Настройка"),
		},
	}
	return func() (*entity.Menu, error) {
		menu.FirstItem = 0
		menu.SelectedItem = 0
		return menu, nil
	}
}

// MenuSensScreen represents the sensor menu screen.
type MenuSensScreen struct {
	name          Name
	activeChanMap ActiveChanMap
	btns          button.Buttons
	templ         *template.Menu
}

// NewMenuSensScreen returns a new instance of MenuSensScreen.
func NewMenuSensScreen(screenName Name, activeChanMap ActiveChanMap, btns button.Buttons,
	storage pubsub.Storage) *MenuSensScreen {

	templ := template.NewMenu(
		string(screenName),
		activeChanMap[screenName],
		storage,
		newMenuSensGetter(),
	)
	return &MenuSensScreen{
		name:          screenName,
		activeChanMap: activeChanMap,
		btns:          btns,
		templ:         templ,
	}
}

// Ready returns true if service was completely started and is ready-to-use now.
func (s *MenuSensScreen) Ready() <-chan struct{} {
	return s.templ.Ready()
}

// StartWithShutdown starts service and wait for context cancellation to shutdown it.
func (s *MenuSensScreen) StartWithShutdown(ctx context.Context) error {
	return s.templ.StartWithShutdown(ctx)
}

// prepareBtnHandlers creates button handlers to apply them after the screen is active
func (s *MenuSensScreen) prepareBtnHandlers() {
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
func (s *MenuSensScreen) btnEsc() error {
	s.activeChanMap[s.name] <- false
	s.activeChanMap[MenuMain] <- true
	return nil
}

// btnEnt represents an enter button handler for screen.
func (s *MenuSensScreen) btnEnt() error {
	menu, err := s.templ.GetFromStorage()
	if err != nil {
		return err
	}

	s.activeChanMap[s.name] <- false
	// set active screen according to selected menu item
	switch menu.SelectedItem {
	case 0:
		s.activeChanMap[SensTemp] <- true
	case 1:
		s.activeChanMap[MenuSensorconf] <- true
	}
	return nil
}
