package template

import (
	"context"

	"github.com/sirupsen/logrus"

	"IvolgaOledManager/internal/app/repo"
	"IvolgaOledManager/internal/pkg/pubsub"
)

// Sensdata represents any sensor data screen.
type Sensdata struct {
	serviceName string
	// will be closed if the service was completely started and is ready-to-use now
	ready chan struct{}

	active         <-chan bool
	storage        pubsub.Storage
	storageKey     string
	btnHandlersReg func()
}

// NewSensdata returns a new instance of SensorData.
func NewSensdata(serviceName string, active <-chan bool,
	storage pubsub.Storage, storageKey string) *Sensdata {

	return &Sensdata{
		serviceName:    serviceName,
		ready:          make(chan struct{}),
		active:         active,
		storage:        storage,
		storageKey:     storageKey,
		btnHandlersReg: func() {},
	}
}

// SetBtnHandlersReg sets btn handlers reg func
// used to update btn handlers when screen is active.
func (s *Sensdata) SetBtnHandlersReg(btnHandlersReg func()) {
	s.btnHandlersReg = btnHandlersReg
}

// Ready signals that the service is ready-to-use.
func (s *Sensdata) Ready() <-chan struct{} {
	return s.ready
}

// StartWithShutdown starts screen service.
// It can be stopped by cancellaiton the given context.
func (s *Sensdata) StartWithShutdown(ctx context.Context) error {
	logrus.Infof("start %s...", s.serviceName)
	defer logrus.Infof("stop %s: ok", s.serviceName)
	// notify that service is ready-to-use
	close(s.ready)

	var (
		activeCtx context.Context
		cancel    context.CancelFunc
	)
	for {
		select {
		case <-ctx.Done():
			if cancel != nil {
				cancel()
			}
			return nil
		case isActive, ok := <-s.active:
			// if chan is closed
			if !ok {
				if cancel != nil {
					cancel()
				}
				return nil
			}
			if !isActive {
				if cancel != nil {
					cancel()
					cancel = nil
				}
				continue
			}
			// if "run" method already started
			if cancel != nil {
				continue
			}
			s.btnHandlersReg()
			activeCtx, cancel = context.WithCancel(ctx)
			go s.run(activeCtx)
		}
	}
}

// run subscribes on sensor data updates and publishes gotten data for render service.
func (s *Sensdata) run(ctx context.Context) {
	s.sendRenderTask()

	newSensorData := s.storage.Subscribe(ctx, s.storageKey)
	for {
		select {
		case <-ctx.Done():
			return
		case <-newSensorData:
			s.sendRenderTask()
		}
	}
}

// sendRenderTask publishes sensor data to storage as renderer for render service.
func (s *Sensdata) sendRenderTask() {
	newData := s.storage.Get(s.storageKey)
	s.storage.Publish(repo.RendererKey, newData)
}
