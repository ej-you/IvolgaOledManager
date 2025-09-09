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

	"github.com/sirupsen/logrus"
)

// newMenuLogLvlActionGetter returns menu getter func for log level actions menu.
func newMenuLogLvlActionGetter() template.MenuGetter {
	menu := &entity.Menu{
		Title: "Действия",
		Items: []*entity.MenuItem{
			entity.NewMenuItem(context.Background(), "Список логов"),
			entity.NewMenuItem(context.Background(), "Очистить"),
		},
	}
	return func() (*entity.Menu, error) {
		menu.FirstItem = 0
		menu.SelectedItem = 0
		return menu, nil
	}
}

// MenuLogLvlActionScreen represents the screen with menu of log level actions.
type MenuLogLvlActionScreen struct {
	name          Name
	activeChanMap ActiveChanMap
	btns          button.Buttons
	templ         *template.Menu
	storage       pubsub.Storage
	logMsgUC      usecase.LogMsgUsecase
}

// NewMenuLogLvlActionScreen returns a new instance of MenuLogLvlActionScreen.
func NewMenuLogLvlActionScreen(screenName Name, activeChanMap ActiveChanMap,
	btns button.Buttons, storage pubsub.Storage,
	logMsgUC usecase.LogMsgUsecase) *MenuLogLvlActionScreen {

	templ := template.NewMenu(
		string(screenName),
		activeChanMap[screenName],
		storage,
		newMenuLogLvlActionGetter(),
	)
	return &MenuLogLvlActionScreen{
		name:          screenName,
		activeChanMap: activeChanMap,
		btns:          btns,
		templ:         templ,
		storage:       storage,
		logMsgUC:      logMsgUC,
	}
}

// Ready returns true if service was completely started and is ready-to-use now.
func (s *MenuLogLvlActionScreen) Ready() <-chan struct{} {
	return s.templ.Ready()
}

// StartWithShutdown starts service and wait for context cancellation to shutdown it.
func (s *MenuLogLvlActionScreen) StartWithShutdown(ctx context.Context) error {
	return s.templ.StartWithShutdown(ctx)
}

// prepareBtnHandlers creates button handlers to apply them after the screen is active
func (s *MenuLogLvlActionScreen) prepareBtnHandlers() {
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
func (s *MenuLogLvlActionScreen) btnEsc() error {
	s.activeChanMap[s.name] <- false
	s.activeChanMap[MenuLogLvl] <- true
	return nil
}

// btnEnt represents an enter button handler for screen.
func (s *MenuLogLvlActionScreen) btnEnt() error {
	menu, err := s.templ.GetFromStorage()
	if err != nil {
		return err
	}

	s.activeChanMap[s.name] <- false
	switch menu.SelectedItem {
	// show messages
	case 0:
		s.activeChanMap[MenuLogMsg] <- true
	// delete all log messages with selected level
	case 1:
		menuLogLvl := s.storage.Get(repo.MenuLogLvl).(*entity.Menu)
		menuLogLvlCtx := menuLogLvl.Items[menuLogLvl.SelectedItem].Ctx
		// get log level count item and delete all messages with it
		logLvlItem := menuLogLvlCtx.Value(entity.LogLvlCountCtxKey).(entity.LogMsgLvlCount)
		if err := s.logMsgUC.DeleteAllWithLevel(strconv.Itoa(logLvlItem.Level)); err != nil {
			return err
		}
		logrus.Infof("log level %d: delete all messages", logLvlItem.Level)
		s.activeChanMap[MenuLogLvl] <- true
	}
	return nil
}
