package template

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"

	"IvolgaOledManager/internal/app/entity"
	"IvolgaOledManager/internal/app/repo"
	"IvolgaOledManager/internal/pkg/pubsub"
)

// SensconfItemGetter represents a getter func for any sensconf item.
// It is used to update sensconf item with screen activation.
type SensconfItemGetter func() (*entity.SensconfItem, error)

// Sensconf represents any sensconf item screen.
type Sensconf struct {
	serviceName string
	// will be closed if the service was completely started and is ready-to-use now
	ready chan struct{}

	active             <-chan bool
	storage            pubsub.Storage
	sensconfItemGetter SensconfItemGetter
	btnHandlersReg     func()
}

// NewSensconf returns a new instance of Sensconf.
func NewSensconf(serviceName string, active <-chan bool,
	storage pubsub.Storage, sensconfItemGetter SensconfItemGetter) *Sensconf {

	return &Sensconf{
		serviceName:        serviceName,
		ready:              make(chan struct{}),
		active:             active,
		storage:            storage,
		sensconfItemGetter: sensconfItemGetter,
		btnHandlersReg:     func() {},
	}
}

// SetBtnHandlersReg sets btn handlers reg func
// used to update btn handlers when screen is active.
func (s *Sensconf) SetBtnHandlersReg(btnHandlersReg func()) {
	s.btnHandlersReg = btnHandlersReg
}

// Ready signals that the service is ready-to-use.
func (s *Sensconf) Ready() <-chan struct{} {
	return s.ready
}

// StartWithShutdown starts screen service.
// It can be stopped by cancellaiton the given context.
func (s *Sensconf) StartWithShutdown(ctx context.Context) error {
	logrus.Infof("start %s...", s.serviceName)
	defer logrus.Infof("stop %s: ok", s.serviceName)
	// notify that service is ready-to-use
	close(s.ready)

	for {
		select {
		case <-ctx.Done():
			return nil
		case isActive, ok := <-s.active:
			// if chan is closed
			if !ok {
				return nil
			}
			if !isActive {
				continue
			}
			s.btnHandlersReg()
			if err := s.sendRenderTask(); err != nil {
				logrus.Error(err)
			}
		}
	}
}

// sendRenderTask publishes sensconf item to storage as renderer for render service.
func (s *Sensconf) sendRenderTask() error {
	// update sensconf item after screen activation
	sensconfItem, err := s.sensconfItemGetter()
	if err != nil {
		return fmt.Errorf("%s: update sensconf item: %w", s.serviceName, err)
	}
	// publish sensconf item for render service
	s.storage.Publish(repo.RendererKey, sensconfItem)
	return nil
}
