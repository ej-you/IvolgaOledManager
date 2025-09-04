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

// newMenuSensconfItemGetter returns menu getter func for sensorconf item menu.
func newMenuSensconfItemGetter() template.MenuGetter {
	menu := &entity.Menu{
		Title: "Действия",
		Items: []*entity.MenuItem{
			entity.NewMenuItem(context.Background(), "Просмотреть"),
			entity.NewMenuItem(context.Background(), "Изменить статус"),
		},
	}
	return func() (*entity.Menu, error) {
		menu.FirstItem = 0
		menu.SelectedItem = 0
		return menu, nil
	}
}

// MenuSensconfItemScreen represents the sensorconf item menu screen.
type MenuSensconfItemScreen struct {
	name          Name
	activeChanMap ActiveChanMap
	btns          button.Buttons
	templ         *template.Menu
	storage       pubsub.Storage
	sensorconfUC  usecase.SensconfUsecase
}

// NewMenuSensconfItemScreen returns a new instance of MenuSensorconfItemScreen.
func NewMenuSensconfItemScreen(screenName Name, activeChanMap ActiveChanMap,
	btns button.Buttons, storage pubsub.Storage,
	sensorconfUC usecase.SensconfUsecase) *MenuSensconfItemScreen {

	templ := template.NewMenu(
		string(screenName),
		activeChanMap[screenName],
		storage,
		newMenuSensconfItemGetter(),
	)
	return &MenuSensconfItemScreen{
		name:          screenName,
		activeChanMap: activeChanMap,
		btns:          btns,
		templ:         templ,
		storage:       storage,
		sensorconfUC:  sensorconfUC,
	}
}

// Ready returns true if service was completely started and is ready-to-use now.
func (s *MenuSensconfItemScreen) Ready() <-chan struct{} {
	return s.templ.Ready()
}

// StartWithShutdown starts service and wait for context cancellation to shutdown it.
func (s *MenuSensconfItemScreen) StartWithShutdown(ctx context.Context) error {
	return s.templ.StartWithShutdown(ctx)
}

// prepareBtnHandlers creates button handlers to apply them after the screen is active
func (s *MenuSensconfItemScreen) prepareBtnHandlers() {
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
func (s *MenuSensconfItemScreen) btnEsc() error {
	s.activeChanMap[s.name] <- false
	s.activeChanMap[MenuSensorconf] <- true
	return nil
}

// btnEnt represents an enter button handler for screen.
func (s *MenuSensconfItemScreen) btnEnt() error {
	menu, err := s.templ.GetFromStorage()
	if err != nil {
		return err
	}

	s.activeChanMap[s.name] <- false
	switch menu.SelectedItem {
	// show info
	case 0:
		s.activeChanMap[Sensconf] <- true
	// change active status
	case 1:
		menuSensconf := s.storage.Get(repo.MenuSensconf).(*entity.Menu)
		menuSensconfCtx := menuSensconf.Items[menuSensconf.SelectedItem].Ctx
		// get sensconf item and change active status
		sensconfItem := menuSensconfCtx.Value(entity.SensconfItemCtxKey).(*entity.SensconfItem)
		sensconfItem.ChangeActive()
		// update sensors' config file
		if err := s.sensorconfUC.UpdateAsMenu(menuSensconf); err != nil {
			return err
		}
		logrus.Infof("sensor %s: new status: %s", sensconfItem.Name, sensconfItem.Status())
	}
	return nil
}
