package screen

import (
	"context"

	"IvolgaOledManager/internal/app/entity"
	"IvolgaOledManager/internal/app/repo"
	"IvolgaOledManager/internal/app/service/button"
	"IvolgaOledManager/internal/app/service/screen/template"
	"IvolgaOledManager/internal/pkg/pubsub"
)

// newLogMsgGetter returns getter func for log message.
func newLogMsgGetter(storage pubsub.Storage) template.LogMsgGetter {
	return func() (*entity.LogMsg, error) {
		// get menu log message from storage
		menuLogMsg := storage.Get(repo.MenuLogMsg).(*entity.Menu)
		menuLogMsgCtx := menuLogMsg.Items[menuLogMsg.SelectedItem].Ctx
		// get log message from selected menu item
		logMsg := menuLogMsgCtx.Value(entity.LogMsgCtxKey).(*entity.LogMsg)

		return logMsg, nil
	}
}

// LogMsgScreen represents the screen with log message.
type LogMsgScreen struct {
	name          Name
	activeChanMap ActiveChanMap
	btns          button.Buttons
	templ         *template.LogMsg
	storage       pubsub.Storage
}

// NewLogMsgScreen returns a new instance of LogMsgScreen.
func NewLogMsgScreen(screenName Name, activeChanMap ActiveChanMap, btns button.Buttons,
	storage pubsub.Storage) *LogMsgScreen {

	templ := template.NewLogMsg(
		string(screenName),
		activeChanMap[screenName],
		storage,
		newLogMsgGetter(storage),
	)
	return &LogMsgScreen{
		name:          screenName,
		activeChanMap: activeChanMap,
		btns:          btns,
		templ:         templ,
		storage:       storage,
	}
}

// Ready returns true if service was completely started and is ready-to-use now.
func (s *LogMsgScreen) Ready() <-chan struct{} {
	return s.templ.Ready()
}

// StartWithShutdown starts service and wait for context cancellation to shutdown it.
func (s *LogMsgScreen) StartWithShutdown(ctx context.Context) error {
	return s.templ.StartWithShutdown(ctx)
}

// prepareBtnHandlers creates button handlers to apply them after the screen is active
func (s *LogMsgScreen) prepareBtnHandlers() {
	btnHandlers := button.Handlers{
		button.Esc: s.btnEsc,
	}

	s.templ.SetBtnHandlersReg(func() {
		s.btns.SetRisingHandlers(btnHandlers)
	})
}

// btnEsc represents an escape button handler for screen.
func (s *LogMsgScreen) btnEsc() error {
	s.activeChanMap[s.name] <- false
	s.activeChanMap[MenuLogMsgAction] <- true
	return nil
}
