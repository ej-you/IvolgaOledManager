package template

import (
	"context"
	"fmt"

	"IvolgaOledManager/internal/app/entity"
	"IvolgaOledManager/internal/app/repo"
	"IvolgaOledManager/internal/pkg/pubsub"

	"github.com/sirupsen/logrus"
)

// MenuGetter represents a getter func for any menu.
// It is used to update menu with screen activation.
type MenuGetter func() (*entity.Menu, error)

// Menu represents any menu screen.
type Menu struct {
	serviceName string
	// will be closed if the service was completely started and is ready-to-use now
	ready chan struct{}

	active         <-chan bool
	storage        pubsub.Storage
	menuGetter     MenuGetter
	btnHandlersReg func()
}

// NewMenu returns a new instance of Menu.
func NewMenu(serviceName string, active <-chan bool,
	storage pubsub.Storage, menuGetter MenuGetter) *Menu {

	return &Menu{
		serviceName:    serviceName,
		ready:          make(chan struct{}),
		active:         active,
		storage:        storage,
		menuGetter:     menuGetter,
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
	logrus.Infof("start %s...", m.serviceName)
	defer logrus.Infof("stop %s: ok", m.serviceName)
	// notify that service is ready-to-use
	close(m.ready)

	var menu *entity.Menu
	var err error
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
			// update menu after screen activation
			menu, err = m.menuGetter()
			if err != nil {
				logrus.Errorf("%s: update menu: %v", m.serviceName, err)
			}
			// publish menu for render service
			m.storage.Publish(repo.RendererKey, menu)
		}
	}
}

// GetFromStorage returns data from storage asserted to menu object.
func (m *Menu) GetFromStorage() (*entity.Menu, error) {
	storageData := m.storage.Get(repo.RendererKey)
	menuInst, ok := storageData.(*entity.Menu)
	if !ok {
		return nil, fmt.Errorf("%s: storage value is not menu object", m.serviceName)
	}
	return menuInst, nil
}

// BtnUpDefault represents a default up button handler for menu screen.
func (m *Menu) BtnUpDefault() error {
	menuInst, err := m.GetFromStorage()
	if err != nil {
		return err
	}
	// update menu and publish into storage as renderer
	menuInst.SelectPrevious()
	m.storage.Publish(repo.RendererKey, menuInst)
	return nil
}

// BtnDownDefault represents a default down button handler for menu screen.
func (m *Menu) BtnDownDefault() error {
	menuInst, err := m.GetFromStorage()
	if err != nil {
		return err
	}
	// update menu and publish into storage as renderer
	menuInst.SelectNext()
	m.storage.Publish(repo.RendererKey, menuInst)
	return nil
}
