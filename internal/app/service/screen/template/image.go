package template

import (
	"context"

	"github.com/sirupsen/logrus"

	"IvolgaOledManager/internal/app/entity"
	"IvolgaOledManager/internal/app/repo"
	"IvolgaOledManager/internal/pkg/pubsub"
)

// Image represents any image screen.
type Image struct {
	screenName string
	// will be closed if the service was completely started and is ready-to-use now
	ready chan struct{}

	active         <-chan bool
	storage        pubsub.Storage
	image          *entity.Image
	btnHandlersReg func()
}

// NewImage returns a new instance of Image.
func NewImage(screenName string, active <-chan bool,
	storage pubsub.Storage, image *entity.Image) *Image {

	return &Image{
		screenName:     screenName,
		ready:          make(chan struct{}),
		active:         active,
		storage:        storage,
		image:          image,
		btnHandlersReg: func() {},
	}
}

// SetBtnHandlersReg sets btn handlers reg func
// used to update btn handlers when screen is active.
func (i *Image) SetBtnHandlersReg(btnHandlersReg func()) {
	i.btnHandlersReg = btnHandlersReg
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
