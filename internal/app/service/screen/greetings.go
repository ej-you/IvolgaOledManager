package screen

import (
	"IvolgaOledManager/internal/app/entity"
	"IvolgaOledManager/internal/app/repo"
	"IvolgaOledManager/internal/pkg/pubsub"
	"context"

	"github.com/sirupsen/logrus"
)

// TODO: create base struct Image as template for all image screens.
type Greetings struct {
	// will be closed if the service was completely started and is ready-to-use now
	ready chan struct{}

	active         chan bool
	storage        pubsub.Storage
	imagePath      string
	btnHandlersReg BtnHandlersRegFunc
}

func NewGreetings(active chan bool, btnHandlersReg BtnHandlersRegFunc, storage pubsub.Storage, imagePath string) *Greetings {
	return &Greetings{
		ready:          make(chan struct{}),
		active:         active,
		storage:        storage,
		imagePath:      imagePath,
		btnHandlersReg: btnHandlersReg,
	}
}

// Ready signals that the service is ready-to-use.
func (g *Greetings) Ready() <-chan struct{} {
	return g.ready
}

func (g *Greetings) StartWithShutdown(ctx context.Context) error {
	logrus.Info("start screen:greetings service...")
	// notify that service is ready-to-use
	close(g.ready)

	for {
		select {
		case <-ctx.Done():
			return nil
		case isActive, ok := <-g.active:
			// if chan is closed
			if !ok {
				return nil
			}
			if !isActive {
				continue
			}
			g.btnHandlersReg()
			g.run()
		}
	}
}

func (g *Greetings) run() {
	image := &entity.Image{
		ImagePath: g.imagePath,
	}

	g.storage.Publish(repo.RendererKey, image)
}
