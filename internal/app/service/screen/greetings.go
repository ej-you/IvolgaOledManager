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
	activeChanMap ActiveChanMap
	btns          button.Buttons
	templ         *template.Image
}

// NewGreetingsScreen returns a new instance of GreetingsScreen.
func NewGreetingsScreen(activeChanMap ActiveChanMap, btns button.Buttons,
	storage pubsub.Storage, imagePath string) *GreetingsScreen {

	image := &entity.Image{ImagePath: imagePath}
	return &GreetingsScreen{
		activeChanMap: activeChanMap,
		btns:          btns,
		templ:         template.NewImage(string(Greetings), activeChanMap[Greetings], storage, image),
	}
}

// Ready returns true if service was completely started and is ready-to-use now.
func (g *GreetingsScreen) Ready() <-chan struct{} {
	return g.templ.Ready()
}

// StartWithShutdown starts service and wait for context cancellation to shutdown it.
func (g *GreetingsScreen) StartWithShutdown(ctx context.Context) error {
	return g.templ.StartWithShutdown(ctx)
}

// prepareBtnHandlers creates button handlers to apply them after the screen is active
func (g *GreetingsScreen) prepareBtnHandlers() {
	btnHandlers := button.Handlers{
		button.ButtonEnt: g.btnEnt,
	}

	g.templ.SetBtnHandlersReg(func() {
		g.btns.SetRisingHandlers(btnHandlers)
	})
}

// btnEnt represents an enter button handler for screen.
func (g *GreetingsScreen) btnEnt() {
	g.activeChanMap[Greetings] <- false
	g.activeChanMap[MenuMain] <- true
}
