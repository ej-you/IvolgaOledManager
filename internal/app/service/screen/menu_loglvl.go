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

// newMenuLogLvlGetter returns menu getter func for log level menu.
// Also it saves new log level into storage with repo.MenuLogLvl key.
func newMenuLogLvlGetter(storage pubsub.Storage,
	logMsgUC usecase.LogMsgUsecase) template.MenuGetter {

	return func() (*entity.Menu, error) {
		menu, err := logMsgUC.GetLevelCountAsMenu()
		if err != nil {
			return nil, err
		}
		storage.Publish(repo.MenuLogLvl, menu)
		return menu, nil
	}
}

// MenuLogLvlScreen represents the screen with menu of log levels.
type MenuLogLvlScreen struct {
	name          Name
	activeChanMap ActiveChanMap
	btns          button.Buttons
	templ         *template.Menu
	storage       pubsub.Storage
	logMsgUC      usecase.LogMsgUsecase
}

// NewMenuLogLvlScreen returns a new instance of MenuLogLvlScreen.
func NewMenuLogLvlScreen(screenName Name, activeChanMap ActiveChanMap, btns button.Buttons,
	storage pubsub.Storage, logMsgUC usecase.LogMsgUsecase) *MenuLogLvlScreen {

	templ := template.NewMenu(
		string(screenName),
		activeChanMap[screenName],
		storage,
		newMenuLogLvlGetter(storage, logMsgUC),
	)
	return &MenuLogLvlScreen{
		name:          screenName,
		activeChanMap: activeChanMap,
		btns:          btns,
		templ:         templ,
		storage:       storage,
		logMsgUC:      logMsgUC,
	}
}

// Ready returns true if service was completely started and is ready-to-use now.
func (s *MenuLogLvlScreen) Ready() <-chan struct{} {
	return s.templ.Ready()
}

// StartWithShutdown starts service and wait for context cancellation to shutdown it.
func (s *MenuLogLvlScreen) StartWithShutdown(ctx context.Context) error {
	return s.templ.StartWithShutdown(ctx)
}

// prepareBtnHandlers creates button handlers to apply them after the screen is active
func (s *MenuLogLvlScreen) prepareBtnHandlers() {
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
func (s *MenuLogLvlScreen) btnEsc() error {
	s.activeChanMap[s.name] <- false
	s.activeChanMap[MenuMain] <- true
	return nil
}

// btnEnt represents an enter button handler for screen.
func (s *MenuLogLvlScreen) btnEnt() error {
	s.activeChanMap[s.name] <- false
	s.activeChanMap[MenuLogLvlAction] <- true
	return nil
}
