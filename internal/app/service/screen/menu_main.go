package screen

import (
	"context"

	"IvolgaOledManager/internal/app/entity"
	"IvolgaOledManager/internal/app/service/button"
	"IvolgaOledManager/internal/app/service/screen/template"
	"IvolgaOledManager/internal/pkg/pubsub"
)

// newMenuMainGetter returns menu getter func for main menu.
func newMenuMainGetter() template.MenuGetter {
	menu := &entity.Menu{
		Title: "Главное меню",
		Items: []*entity.MenuItem{
			entity.NewMenuItem(context.Background(), "Датчики"),
			entity.NewMenuItem(context.Background(), "Назад"),
			entity.NewMenuItem(context.Background(), "И ещё"),
			entity.NewMenuItem(context.Background(), "И для проверки прокрутки - ещё!"),
			entity.NewMenuItem(context.Background(), "And the last for test"),
		},
	}
	return func() (*entity.Menu, error) {
		menu.FirstItem = 0
		menu.SelectedItem = 0
		return menu, nil
	}
}

// MenuMainScreen represents the main menu screen.
type MenuMainScreen struct {
	name          Name
	activeChanMap ActiveChanMap
	btns          button.Buttons
	templ         *template.Menu
}

// NewMenuMainScreen returns a new instance of MenuMainScreen.
func NewMenuMainScreen(screenName Name, activeChanMap ActiveChanMap, btns button.Buttons,
	storage pubsub.Storage) *MenuMainScreen {

	templ := template.NewMenu(
		string(screenName),
		activeChanMap[screenName],
		storage,
		newMenuMainGetter(),
	)
	return &MenuMainScreen{
		name:          screenName,
		activeChanMap: activeChanMap,
		btns:          btns,
		templ:         templ,
	}
}

// Ready returns true if service was completely started and is ready-to-use now.
func (s *MenuMainScreen) Ready() <-chan struct{} {
	return s.templ.Ready()
}

// StartWithShutdown starts service and wait for context cancellation to shutdown it.
func (s *MenuMainScreen) StartWithShutdown(ctx context.Context) error {
	return s.templ.StartWithShutdown(ctx)
}

// prepareBtnHandlers creates button handlers to apply them after the screen is active
func (s *MenuMainScreen) prepareBtnHandlers() {
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
func (s *MenuMainScreen) btnEsc() error {
	s.activeChanMap[s.name] <- false
	s.activeChanMap[ImgGreetings] <- true
	return nil
}

// btnEnt represents an enter button handler for screen.
func (s *MenuMainScreen) btnEnt() error {
	menu, err := s.templ.GetFromStorage()
	if err != nil {
		return err
	}

	s.activeChanMap[s.name] <- false
	// set active screen according to selected menu item
	switch menu.SelectedItem {
	case 0:
		s.activeChanMap[MenuSens] <- true
	case 1:
		s.activeChanMap[ImgGreetings] <- true
	}
	return nil
}
