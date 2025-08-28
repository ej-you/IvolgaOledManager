package screen

import (
	"IvolgaOledManager/internal/app/repo"
	"IvolgaOledManager/internal/pkg/pubsub"
	"context"

	"github.com/sirupsen/logrus"
)

// TODO: create base struct SensorData as template for all sensor data screens.
type Temperature struct {
	// will be closed if the service was completely started and is ready-to-use now
	ready chan struct{}

	active         chan bool
	storage        pubsub.Storage
	storageKey     string
	btnHandlersReg BtnHandlersRegFunc
}

func NewTemperature(active chan bool, btnHandlersReg BtnHandlersRegFunc,
	storage pubsub.Storage) *Temperature {

	return &Temperature{
		ready:          make(chan struct{}),
		active:         active,
		storage:        storage,
		storageKey:     repo.SensTempKey,
		btnHandlersReg: btnHandlersReg,
	}
}

// Ready signals that the service is ready-to-use.
func (t *Temperature) Ready() <-chan struct{} {
	return t.ready
}

func (t *Temperature) StartWithShutdown(ctx context.Context) error {
	logrus.Info("start screen:sensordata:temperature service...")
	// notify that service is ready-to-use
	close(t.ready)

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
		case isActive, ok := <-t.active:
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
			if cancel != nil {
				continue
			}
			t.btnHandlersReg()
			activeCtx, cancel = context.WithCancel(ctx)
			go t.run(activeCtx)
		}
	}
}

func (t *Temperature) run(ctx context.Context) {
	t.sendRenderTask()

	newSensorData := t.storage.Subscribe(ctx, t.storageKey)
	for {
		select {
		case <-ctx.Done():
			return
		case <-newSensorData:
			t.sendRenderTask()
		}
	}
}

func (t *Temperature) sendRenderTask() {
	newData := t.storage.Get(t.storageKey)
	t.storage.Publish(repo.RendererKey, newData)
}
