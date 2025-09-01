package screen

import (
	"context"

	"IvolgaOledManager/internal/app/entity"
	"IvolgaOledManager/internal/app/entity/menuitem"
	"IvolgaOledManager/internal/app/repo"
	"IvolgaOledManager/internal/app/service/button"
	"IvolgaOledManager/internal/app/service/screen/template"
	"IvolgaOledManager/internal/pkg/pubsub"

	"github.com/sirupsen/logrus"
)

// Static main menu instance.
var _mainMenuInst = &entity.Menu{
	Title: "Главное меню",
	Items: []*menuitem.MenuItem{
		menuitem.New("Данные датчиков"),
		menuitem.New("Что-то ещё"),
		menuitem.New("И ещё"),
		menuitem.New("И для проверки прокрутки - ещё!"),
		menuitem.New("And the last for test"),
	},
}

// MenuMainScreen represents the main menu data screen.
type MenuMainScreen struct {
	activeChanMap ActiveChanMap
	btns          button.Buttons
	storage       pubsub.Storage
	templ         *template.Menu
}

// NewMenuMainScreen returns a new instance of MenuMainScreen.
func NewMenuMainScreen(activeChanMap ActiveChanMap, btns button.Buttons,
	storage pubsub.Storage) *MenuMainScreen {

	templ := template.NewMenu(
		string(MenuMain),
		activeChanMap[MenuMain],
		storage,
		_mainMenuInst,
	)
	return &MenuMainScreen{
		activeChanMap: activeChanMap,
		btns:          btns,
		storage:       storage,
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
		button.ButtonEsc:  s.btnEsc,
		button.ButtonUp:   s.btnUp,
		button.ButtonDown: s.btnDown,
		button.ButtonEnt:  s.btnEnt,
	}

	s.templ.SetBtnHandlersReg(func() {
		s.btns.SetRisingHandlers(btnHandlers)
	})
}

// btnEsc represents an escape button handler for screen.
func (s *MenuMainScreen) btnEsc() {
	s.activeChanMap[MenuMain] <- false
	s.activeChanMap[Greetings] <- true
}

// btnUp represents an up button handler for screen.
func (s *MenuMainScreen) btnUp() {
	// get object from storage and assert it to menu instance
	storageData := s.storage.Get(repo.RendererKey)
	menuInst, ok := storageData.(*entity.Menu)
	if !ok {
		logrus.Error("screen:menu:main: btn up: storage value is not menu object")
	}
	// update menu and publish into storage as renderer
	menuInst.SelectPrevious()
	s.storage.Publish(repo.RendererKey, menuInst)
}

// btnDown represents an down button handler for screen.
func (s *MenuMainScreen) btnDown() {
	// get object from storage and assert it to menu instance
	storageData := s.storage.Get(repo.RendererKey)
	menuInst, ok := storageData.(*entity.Menu)
	if !ok {
		logrus.Error("screen:menu:main: btn down: storage value is not menu object")
	}
	// update menu and publish into storage as renderer
	menuInst.SelectNext()
	s.storage.Publish(repo.RendererKey, menuInst)
}

// btnEnt represents an enter button handler for screen.
func (s *MenuMainScreen) btnEnt() {
	// get object from storage and assert it to menu instance
	storageData := s.storage.Get(repo.RendererKey)
	menuInst, ok := storageData.(*entity.Menu)
	if !ok {
		logrus.Error("screen:menu:main: btn ent: storage value is not menu object")
	}

	s.activeChanMap[MenuMain] <- false
	// set active screen according to selected menu item
	switch menuInst.Items[menuInst.SelectedItem] {
	case _mainMenuInst.Items[0]:
		s.activeChanMap[SensTemp] <- true
	case _mainMenuInst.Items[1]:
		s.activeChanMap[SensHumid] <- true
	case _mainMenuInst.Items[2]:
		s.activeChanMap[SensPress] <- true
	case _mainMenuInst.Items[3]:
		s.activeChanMap[SensWindSpeed] <- true
	case _mainMenuInst.Items[4]:
		s.activeChanMap[SensWindDir] <- true
	}
}
