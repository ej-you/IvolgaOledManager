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

// newMenuSensorconfItemGetter returns menu getter func for sensorconf item menu.
func newMenuSensorconfItemGetter() template.MenuGetter {
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

// MenuSensorconfItemScreen represents the sensorconf item menu screen.
type MenuSensorconfItemScreen struct {
	name          Name
	activeChanMap ActiveChanMap
	btns          button.Buttons
	templ         *template.Menu
	storage       pubsub.Storage
	sensorconfUC  usecase.SensconfUsecase
}

// NewMenuSensorconfItemScreen returns a new instance of MenuSensorconfItemScreen.
func NewMenuSensorconfItemScreen(screenName Name, activeChanMap ActiveChanMap,
	btns button.Buttons, storage pubsub.Storage,
	sensorconfUC usecase.SensconfUsecase) *MenuSensorconfItemScreen {

	templ := template.NewMenu(
		string(screenName),
		activeChanMap[screenName],
		storage,
		newMenuSensorconfItemGetter(),
	)
	return &MenuSensorconfItemScreen{
		name:          screenName,
		activeChanMap: activeChanMap,
		btns:          btns,
		templ:         templ,
		storage:       storage,
		sensorconfUC:  sensorconfUC,
	}
}

// Ready returns true if service was completely started and is ready-to-use now.
func (s *MenuSensorconfItemScreen) Ready() <-chan struct{} {
	return s.templ.Ready()
}

// StartWithShutdown starts service and wait for context cancellation to shutdown it.
func (s *MenuSensorconfItemScreen) StartWithShutdown(ctx context.Context) error {
	return s.templ.StartWithShutdown(ctx)
}

// prepareBtnHandlers creates button handlers to apply them after the screen is active
func (s *MenuSensorconfItemScreen) prepareBtnHandlers() {
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
func (s *MenuSensorconfItemScreen) btnEsc() error {
	s.activeChanMap[s.name] <- false
	s.activeChanMap[MenuSensorconf] <- true
	return nil
}

// btnEnt represents an enter button handler for screen.
func (s *MenuSensorconfItemScreen) btnEnt() error {
	curMenu, err := s.templ.GetFromStorage()
	if err != nil {
		return err
	}

	menuSensconf := s.storage.Get(repo.MenuSensconf).(*entity.Menu)

	selectedItemCtx := menuSensconf.Items[menuSensconf.SelectedItem].Ctx

	sensconfItemFromStorage := selectedItemCtx.Value(entity.SensconfItemCtxKey)
	switch curMenu.SelectedItem {
	// show info
	case 0:
		s.storage.Publish(repo.RendererKey, sensconfItemFromStorage)
		s.activeChanMap[s.name] <- false
		s.activeChanMap[ImgGreetings] <- true
	// change active status
	case 1:
		sensconf := selectedItemCtx.Value(entity.SensconfCtxKey).(entity.Sensconf)
		sensconfItem := sensconfItemFromStorage.(*entity.SensconfItem)
		sensconfItem.ChangeActive()
		if err := s.sensorconfUC.Update(sensconf); err != nil {
			return err
		}
		return s.btnEsc()
	}
	return nil
}
