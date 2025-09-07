package screen

import (
	"context"

	"IvolgaOledManager/internal/app/entity"
	"IvolgaOledManager/internal/app/repo"
	"IvolgaOledManager/internal/app/service/button"
	"IvolgaOledManager/internal/app/service/screen/template"
	"IvolgaOledManager/internal/app/usecase"
	"IvolgaOledManager/internal/pkg/pubsub"

	"github.com/sirupsen/logrus"
)

// newMenuLogMsgActionGetter returns menu getter func for log message actions menu.
func newMenuLogMsgActionGetter() template.MenuGetter {
	menu := &entity.Menu{
		Title: "Действия",
		Items: []*entity.MenuItem{
			entity.NewMenuItem(context.Background(), "Просмотреть"),
			entity.NewMenuItem(context.Background(), "Удалить"),
		},
	}
	return func() (*entity.Menu, error) {
		menu.FirstItem = 0
		menu.SelectedItem = 0
		return menu, nil
	}
}

// MenuLogMsgActionScreen represents the screen with menu of log message actions.
type MenuLogMsgActionScreen struct {
	name          Name
	activeChanMap ActiveChanMap
	btns          button.Buttons
	templ         *template.Menu
	storage       pubsub.Storage
	logMsgUC      usecase.LogMsgUsecase
}

// NewMenuLogMsgActionScreen returns a new instance of MenuLogMsgActionScreen.
func NewMenuLogMsgActionScreen(screenName Name, activeChanMap ActiveChanMap,
	btns button.Buttons, storage pubsub.Storage,
	logMsgUC usecase.LogMsgUsecase) *MenuLogMsgActionScreen {

	templ := template.NewMenu(
		string(screenName),
		activeChanMap[screenName],
		storage,
		newMenuLogMsgActionGetter(),
	)
	return &MenuLogMsgActionScreen{
		name:          screenName,
		activeChanMap: activeChanMap,
		btns:          btns,
		templ:         templ,
		storage:       storage,
		logMsgUC:      logMsgUC,
	}
}

// Ready returns true if service was completely started and is ready-to-use now.
func (s *MenuLogMsgActionScreen) Ready() <-chan struct{} {
	return s.templ.Ready()
}

// StartWithShutdown starts service and wait for context cancellation to shutdown it.
func (s *MenuLogMsgActionScreen) StartWithShutdown(ctx context.Context) error {
	return s.templ.StartWithShutdown(ctx)
}

// prepareBtnHandlers creates button handlers to apply them after the screen is active
func (s *MenuLogMsgActionScreen) prepareBtnHandlers() {
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
func (s *MenuLogMsgActionScreen) btnEsc() error {
	s.activeChanMap[s.name] <- false
	s.activeChanMap[MenuLogMsg] <- true
	return nil
}

// btnEnt represents an enter button handler for screen.
func (s *MenuLogMsgActionScreen) btnEnt() error {
	menu, err := s.templ.GetFromStorage()
	if err != nil {
		return err
	}

	s.activeChanMap[s.name] <- false
	switch menu.SelectedItem {
	// show info
	case 0:
		s.activeChanMap[LogMsg] <- true
	// delete selected log messages
	case 1:
		menuLogMsg := s.storage.Get(repo.MenuLogLvl).(*entity.Menu)
		menuLogMsgCtx := menuLogMsg.Items[menuLogMsg.SelectedItem].Ctx
		// get log message item and delete it
		logMsgItem := menuLogMsgCtx.Value(entity.LogMsgCtxKey).(*entity.LogMsgWithLevel)
		if err := s.logMsgUC.DeleteByID(logMsgItem.ID); err != nil {
			return err
		}
		logrus.Infof("log message %s - %s: deleted", logMsgItem.ID, logMsgItem.Header)
	}
	return nil
}
