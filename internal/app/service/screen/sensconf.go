package screen

import (
	"context"

	"IvolgaOledManager/internal/app/entity"
	"IvolgaOledManager/internal/app/repo"
	"IvolgaOledManager/internal/app/service/button"
	"IvolgaOledManager/internal/app/service/screen/template"
	"IvolgaOledManager/internal/pkg/pubsub"
)

// newSensconfItemGetter returns sensconf item getter func for sensconf item.
func newSensconfItemGetter(storage pubsub.Storage) template.SensconfItemGetter {
	return func() (*entity.SensconfItem, error) {
		// get menu sensconf from storage
		menuSensconf := storage.Get(repo.MenuSensconf).(*entity.Menu)
		menuSensconfCtx := menuSensconf.Items[menuSensconf.SelectedItem].Ctx
		// get sensconf item from selected menu item
		sensconfItemFromStorage := menuSensconfCtx.Value(entity.SensconfItemCtxKey)
		sensconfItem := sensconfItemFromStorage.(*entity.SensconfItem)

		return sensconfItem, nil
	}
}

// SensconfScreen represents the sensconf item screen.
type SensconfScreen struct {
	name          Name
	activeChanMap ActiveChanMap
	btns          button.Buttons
	templ         *template.Sensconf
	storage       pubsub.Storage
}

// NewSensconfScreen returns a new instance of MenuMainScreen.
func NewSensconfScreen(screenName Name, activeChanMap ActiveChanMap, btns button.Buttons,
	storage pubsub.Storage) *SensconfScreen {

	templ := template.NewSensconf(
		string(screenName),
		activeChanMap[screenName],
		storage,
		newSensconfItemGetter(storage),
	)
	return &SensconfScreen{
		name:          screenName,
		activeChanMap: activeChanMap,
		btns:          btns,
		templ:         templ,
		storage:       storage,
	}
}

// Ready returns true if service was completely started and is ready-to-use now.
func (s *SensconfScreen) Ready() <-chan struct{} {
	return s.templ.Ready()
}

// StartWithShutdown starts service and wait for context cancellation to shutdown it.
func (s *SensconfScreen) StartWithShutdown(ctx context.Context) error {
	return s.templ.StartWithShutdown(ctx)
}

// prepareBtnHandlers creates button handlers to apply them after the screen is active
func (s *SensconfScreen) prepareBtnHandlers() {
	btnHandlers := button.Handlers{
		button.Esc: s.btnEsc,
	}

	s.templ.SetBtnHandlersReg(func() {
		s.btns.SetRisingHandlers(btnHandlers)
	})
}

// btnEsc represents an escape button handler for screen.
func (s *SensconfScreen) btnEsc() error {
	s.activeChanMap[s.name] <- false
	s.activeChanMap[MenuSensconfItem] <- true
	return nil
}
