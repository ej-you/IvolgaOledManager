package screen

import (
	"context"

	"IvolgaOledManager/internal/app/entity"
	"IvolgaOledManager/internal/app/service/button"
	"IvolgaOledManager/internal/app/service/screen/template"
	"IvolgaOledManager/internal/app/usecase"
	"IvolgaOledManager/internal/pkg/pubsub"
)

// MenuSensorconfScreen represents the main menu data screen.
type MenuSensorconfScreen struct {
	activeChanMap ActiveChanMap
	btns          button.Buttons
	storage       pubsub.Storage
	sensorconfUC  usecase.SensorconfUsecase
	templ         *template.Menu
}

// NewMenuSensorconfScreen returns a new instance of MenuMainScreen.
func NewMenuSensorconfScreen(screenName Name, activeChanMap ActiveChanMap, btns button.Buttons,
	storage pubsub.Storage, sensorconfUC usecase.SensorconfUsecase) *MenuSensorconfScreen {

	templ := template.NewMenu(
		string(screenName),
		activeChanMap[screenName],
		storage,
		newMenuSensconfGetter(sensorconfUC),
	)
	return &MenuSensorconfScreen{
		activeChanMap: activeChanMap,
		btns:          btns,
		storage:       storage,
		sensorconfUC:  sensorconfUC,
		templ:         templ,
	}
}

// Ready returns true if service was completely started and is ready-to-use now.
func (s *MenuSensorconfScreen) Ready() <-chan struct{} {
	return s.templ.Ready()
}

// StartWithShutdown starts service and wait for context cancellation to shutdown it.
func (s *MenuSensorconfScreen) StartWithShutdown(ctx context.Context) error {
	return s.templ.StartWithShutdown(ctx)
}

// prepareBtnHandlers creates button handlers to apply them after the screen is active
func (s *MenuSensorconfScreen) prepareBtnHandlers() {
	btnHandlers := button.Handlers{
		button.Esc:  s.btnEsc,
		button.Up:   s.templ.BtnUpDefault,
		button.Down: s.templ.BtnDownDefault,
		// button.ButtonEnt:  s.btnEnt,
	}

	s.templ.SetBtnHandlersReg(func() {
		s.btns.SetRisingHandlers(btnHandlers)
	})
}

// btnEsc represents an escape button handler for screen.
func (s *MenuSensorconfScreen) btnEsc() error {
	s.activeChanMap[MenuSensorconf] <- false
	s.activeChanMap[MenuSens] <- true
	return nil
}

// // btnEnt represents an enter button handler for screen.
// func (s *MenuSensorconfScreen) btnEnt() error {
// 	// get object from storage and assert it to menu instance
// 	storageData := s.storage.Get(repo.RendererKey)
// 	menuInst, ok := storageData.(*entity.Menu)
// 	if !ok {
// 		logrus.Error("screen:menu:main: btn ent: storage value is not menu object")
// 	}

// 	s.activeChanMap[MenuMain] <- false
// 	// set active screen according to selected menu item
// 	switch menuInst.Items[menuInst.SelectedItem] {
// 	case _mainMenuInst.Items[0]:
// 		s.activeChanMap[SensTemp] <- true
// 	case _mainMenuInst.Items[1]:
// 		s.activeChanMap[SensHumid] <- true
// 	case _mainMenuInst.Items[2]:
// 		s.activeChanMap[SensPress] <- true
// 	case _mainMenuInst.Items[3]:
// 		s.activeChanMap[SensWindSpeed] <- true
// 	case _mainMenuInst.Items[4]:
// 		s.activeChanMap[SensWindDir] <- true
// 	}
// return nil
// }

// newMenuSensconfGetter returns menu getter func for sensorconf menu.
func newMenuSensconfGetter(sensorconfUC usecase.SensorconfUsecase) template.MenuGetter {
	return func() (*entity.Menu, error) {
		data, err := sensorconfUC.GetSensorconf()
		if err != nil {
			return nil, err
		}
		return sensorconfUC.ToMenu(data), nil
	}
}
