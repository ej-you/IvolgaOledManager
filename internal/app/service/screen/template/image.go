package template

import (
	"context"

	"github.com/sirupsen/logrus"

	"IvolgaOledManager/internal/app/entity"
	"IvolgaOledManager/internal/app/repo"
	"IvolgaOledManager/internal/pkg/pubsub"
)

// SensorData represents a template for image screen.
type Image struct {
	// will be closed if the service was completely started and is ready-to-use now
	ready chan struct{}

	active         <-chan bool
	btnHandlersReg func()
	storage        pubsub.Storage
	image          *entity.Image
	screenName     string
}

// NewImage returns a new instance of Image.
func NewImage(active <-chan bool, btnHandlersReg func(),
	storage pubsub.Storage, image *entity.Image, screenName string) *Image {

	return &Image{
		ready:          make(chan struct{}),
		active:         active,
		storage:        storage,
		image:          image,
		screenName:     screenName,
		btnHandlersReg: btnHandlersReg,
	}
}

// Ready signals that the service is ready-to-use.
func (i *Image) Ready() <-chan struct{} {
	return i.ready
}

// StartWithShutdown starts screen service.
// It can be stopped by cancellaiton the given context.
func (i *Image) StartWithShutdown(ctx context.Context) error {
	logrus.Infof("start screen:%s service...", i.screenName)
	// notify that service is ready-to-use
	close(i.ready)

	for {
		select {
		case <-ctx.Done():
			return nil
		case isActive, ok := <-i.active:
			// if chan is closed
			if !ok {
				return nil
			}
			// screen is inactive
			if !isActive {
				continue
			}
			i.btnHandlersReg()
			// publish image data for render service
			i.storage.Publish(repo.RendererKey, i.image)
		}
	}
}
