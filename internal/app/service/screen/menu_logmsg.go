package screen

import (
	"context"
	"strconv"

	"IvolgaOledManager/internal/app/entity"
	"IvolgaOledManager/internal/app/repo"
	"IvolgaOledManager/internal/app/service/button"
	"IvolgaOledManager/internal/app/service/screen/template"
	"IvolgaOledManager/internal/app/usecase"
	"IvolgaOledManager/internal/pkg/pubsub"
)

// newMenuLogMsgGetter returns menu getter func for log messages menu.
// Also it saves new log messages into storage with repo.MenuLogMsg key.
func newMenuLogMsgGetter(storage pubsub.Storage,
	logMsgUC usecase.LogMsgUsecase) template.MenuGetter {

	return func() (*entity.Menu, error) {
		// get menu log level from storage
		menuLogLvl := storage.Get(repo.MenuLogLvl).(*entity.Menu)
		menuLogLvlCtx := menuLogLvl.Items[menuLogLvl.SelectedItem].Ctx
		// get log level count item and delete all messages with it
		logLvlItem := menuLogLvlCtx.Value(entity.LogLvlCountCtxKey).(entity.LogMsgLvlCount)

		// get log messages as menu
		menu, err := logMsgUC.GetWithLevelAsMenu(strconv.Itoa(logLvlItem.Level))
		if err != nil {
			return nil, err
		}
		storage.Publish(repo.MenuLogMsg, menu)
		return menu, nil
	}
}

// MenuLogMsgScreen represents the screen with menu of log level messages.
type MenuLogMsgScreen struct {
	name          Name
	activeChanMap ActiveChanMap
	btns          button.Buttons
	templ         *template.Menu
	storage       pubsub.Storage
	logMsgUC      usecase.LogMsgUsecase
}

// NewMenuLogMsgScreen returns a new instance of MenuLogMsgScreen.
func NewMenuLogMsgScreen(screenName Name, activeChanMap ActiveChanMap, btns button.Buttons,
	storage pubsub.Storage, logMsgUC usecase.LogMsgUsecase) *MenuLogMsgScreen {

	templ := template.NewMenu(
		string(screenName),
		activeChanMap[screenName],
		storage,
		newMenuLogMsgGetter(storage, logMsgUC),
	)
	return &MenuLogMsgScreen{
		name:          screenName,
		activeChanMap: activeChanMap,
		btns:          btns,
		templ:         templ,
		storage:       storage,
		logMsgUC:      logMsgUC,
	}
}

// Ready returns true if service was completely started and is ready-to-use now.
func (s *MenuLogMsgScreen) Ready() <-chan struct{} {
	return s.templ.Ready()
}

// StartWithShutdown starts service and wait for context cancellation to shutdown it.
func (s *MenuLogMsgScreen) StartWithShutdown(ctx context.Context) error {
	return s.templ.StartWithShutdown(ctx)
}

// prepareBtnHandlers creates button handlers to apply them after the screen is active
func (s *MenuLogMsgScreen) prepareBtnHandlers() {
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
func (s *MenuLogMsgScreen) btnEsc() error {
	s.activeChanMap[s.name] <- false
	s.activeChanMap[MenuLogLvlAction] <- true
	return nil
}

// btnEnt represents an enter button handler for screen.
func (s *MenuLogMsgScreen) btnEnt() error {
	s.activeChanMap[s.name] <- false
	s.activeChanMap[MenuLogMsgAction] <- true
	return nil
}
