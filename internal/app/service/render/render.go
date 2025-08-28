// Package render provides render service to create output tasks
// and print them out via display service.
package render

import (
	"context"
	"errors"

	"IvolgaOledManager/internal/app/repo"
	"IvolgaOledManager/internal/pkg/display"
	"IvolgaOledManager/internal/pkg/pubsub"

	"github.com/sirupsen/logrus"
)

// Render represents render service that sends
// signals to display service to update display screen.
type Render struct {
	// will be closed if the service was completely started and is ready-to-use now
	ready chan struct{}

	displayService *display.Service
	storage        pubsub.Storage
}

// New returns a new instance of Render.
func New(displayDev *display.Service, storage pubsub.Storage) *Render {
	return &Render{
		ready:          make(chan struct{}),
		displayService: displayDev,
		storage:        storage,
	}
}

// StartWithShutdown starts render service.
// It may be stopped by context cancellaiton.
func (r *Render) StartWithShutdown(ctx context.Context) error {
	logrus.Info("start render service...")
	// notify chan for display updates
	needRender := r.storage.Subscribe(ctx, repo.RendererKey)

	// notify that service is ready-to-use
	close(r.ready)

	var err error
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-needRender:
			if err = r.updateDisplay(); err != nil {
				logrus.Errorf("update display: %v", err)
			}
		}
	}
}

// Ready signals that the service is ready-to-use.
func (r *Render) Ready() <-chan struct{} {
	return r.ready
}

// updateDisplay sets new renderer for display.
func (r *Render) updateDisplay() error {
	// get object from storage and assert it to display renderer
	storageData := r.storage.Get(repo.RendererKey)
	renderer, ok := storageData.(display.Renderer)
	if !ok {
		return errors.New("invalid renderer in storage")
	}

	// set new renderer object for display
	r.displayService.SetRenderer(renderer)
	return nil
}
