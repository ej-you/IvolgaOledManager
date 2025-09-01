package template

import (
	"context"

	"github.com/sirupsen/logrus"

	"IvolgaOledManager/internal/app/entity"
	"IvolgaOledManager/internal/app/repo"
	"IvolgaOledManager/internal/pkg/pubsub"
)

// Menu represents any menu screen.
type Menu struct {
	screenName string
	// will be closed if the service was completely started and is ready-to-use now
	ready chan struct{}

	active         <-chan bool
	storage        pubsub.Storage
	menu           *entity.Menu
	btnHandlersReg func()
}

// NewMenu returns a new instance of Menu.
func NewMenu(screenName string, active <-chan bool,
	storage pubsub.Storage, menu *entity.Menu) *Menu {

	return &Menu{
		screenName:     screenName,
		ready:          make(chan struct{}),
		active:         active,
		storage:        storage,
		menu:           menu,
		btnHandlersReg: func() {},
	}
}

// SetBtnHandlersReg sets btn handlers reg func
// used to update btn handlers when screen is active.
func (m *Menu) SetBtnHandlersReg(btnHandlersReg func()) {
	m.btnHandlersReg = btnHandlersReg
}

// Ready signals that the service is ready-to-use.
func (m *Menu) Ready() <-chan struct{} {
	return m.ready
}

// StartWithShutdown starts screen service.
// It can be stopped by cancellaiton the given context.
func (m *Menu) StartWithShutdown(ctx context.Context) error {
	logrus.Infof("start screen:%s service...", m.screenName)
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
			m.storage.Publish(repo.RendererKey, m.menu)
		}
	}
}
