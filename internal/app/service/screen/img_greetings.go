package screen

import (
	"context"

	"IvolgaOledManager/internal/app/entity"
	"IvolgaOledManager/internal/app/service/button"
	"IvolgaOledManager/internal/app/service/screen/template"
	"IvolgaOledManager/internal/pkg/pubsub"
)

// GreetingsScreen represents a greetings screen.
type GreetingsScreen struct {
	name          Name
	activeChanMap ActiveChanMap
	btns          button.Buttons
	templ         *template.Image
}

// NewGreetingsScreen returns a new instance of GreetingsScreen.
func NewGreetingsScreen(screenName Name, activeChanMap ActiveChanMap, btns button.Buttons,
	storage pubsub.Storage, imagePath string) *GreetingsScreen {

	templ := template.NewImage(
		string(screenName),
		activeChanMap[screenName],
		storage,
		&entity.Image{ImagePath: imagePath},
	)
	return &GreetingsScreen{
		name:          screenName,
		activeChanMap: activeChanMap,
		btns:          btns,
		templ:         templ,
	}
}

// Ready returns true if service was completely started and is ready-to-use now.
func (s *GreetingsScreen) Ready() <-chan struct{} {
	return s.templ.Ready()
}

// StartWithShutdown starts service and wait for context cancellation to shutdown it.
func (s *GreetingsScreen) StartWithShutdown(ctx context.Context) error {
	return s.templ.StartWithShutdown(ctx)
}

// prepareBtnHandlers creates button handlers to apply them after the screen is active
func (s *GreetingsScreen) prepareBtnHandlers() {
	btnHandlers := button.Handlers{
		button.Ent: s.btnEnt,
	}

	s.templ.SetBtnHandlersReg(func() {
		s.btns.SetRisingHandlers(btnHandlers)
	})
}

// btnEnt represents an enter button handler for screen.
func (s *GreetingsScreen) btnEnt() error {
	s.activeChanMap[s.name] <- false
	s.activeChanMap[MenuMain] <- true
	return nil
}
