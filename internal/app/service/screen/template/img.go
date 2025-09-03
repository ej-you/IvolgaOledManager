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
	serviceName string
	// will be closed if the service was completely started and is ready-to-use now
	ready chan struct{}

	active         <-chan bool
	storage        pubsub.Storage
	image          *entity.Image
	btnHandlersReg func()
}

// NewImage returns a new instance of Image.
func NewImage(serviceName string, active <-chan bool,
	storage pubsub.Storage, image *entity.Image) *Image {

	return &Image{
		serviceName:    serviceName,
		ready:          make(chan struct{}),
		active:         active,
		storage:        storage,
		image:          image,
		btnHandlersReg: func() {},
	}
}

// SetBtnHandlersReg sets btn handlers reg func
// used to update btn handlers when screen is active.
func (m *Image) SetBtnHandlersReg(btnHandlersReg func()) {
	m.btnHandlersReg = btnHandlersReg
}

// Ready signals that the service is ready-to-use.
func (m *Image) Ready() <-chan struct{} {
	return m.ready
}

// StartWithShutdown starts screen service.
// It can be stopped by cancellaiton the given context.
func (m *Image) StartWithShutdown(ctx context.Context) error {
	logrus.Infof("start %s...", m.serviceName)
	defer logrus.Infof("stop %s: ok", m.serviceName)
	// notify that service is ready-to-use
	close(m.ready)

	for {
		select {
		case <-ctx.Done():
			return nil
		case isActive, ok := <-m.active:
			// if chan is closed
			if !ok {
				return nil
			}
			// screen is inactive
			if !isActive {
				continue
			}
			m.btnHandlersReg()
			// publish image data for render service
			m.storage.Publish(repo.RendererKey, m.image)
		}
	}
}
