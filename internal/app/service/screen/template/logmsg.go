package template

import (
	"context"
	"fmt"

	"IvolgaOledManager/internal/app/entity"
	"IvolgaOledManager/internal/app/repo"
	"IvolgaOledManager/internal/pkg/pubsub"

	"github.com/sirupsen/logrus"
)

// LogMsgGetter represents a getter func for any log message.
// It is used to update log message with screen activation.
type LogMsgGetter func() (*entity.LogMsg, error)

// LogMsg represents any log message screen.
type LogMsg struct {
	serviceName string
	// will be closed if the service was completely started and is ready-to-use now
	ready chan struct{}

	active         <-chan bool
	storage        pubsub.Storage
	logMsgGetter   LogMsgGetter
	btnHandlersReg func()
}

// NewLogMsg returns a new instance of LogMsg.
func NewLogMsg(serviceName string, active <-chan bool,
	storage pubsub.Storage, logMsgGetter LogMsgGetter) *LogMsg {

	return &LogMsg{
		serviceName:    serviceName,
		ready:          make(chan struct{}),
		active:         active,
		storage:        storage,
		logMsgGetter:   logMsgGetter,
		btnHandlersReg: func() {},
	}
}

// SetBtnHandlersReg sets btn handlers reg func
// used to update btn handlers when screen is active.
func (l *LogMsg) SetBtnHandlersReg(btnHandlersReg func()) {
	l.btnHandlersReg = btnHandlersReg
}

// Ready signals that the service is ready-to-use.
func (l *LogMsg) Ready() <-chan struct{} {
	return l.ready
}

// StartWithShutdown starts screen service.
// It can be stopped by cancellaiton the given context.
func (l *LogMsg) StartWithShutdown(ctx context.Context) error {
	logrus.Infof("start %s...", l.serviceName)
	defer logrus.Infof("stop %s: ok", l.serviceName)
	// notify that service is ready-to-use
	close(l.ready)

	for {
		select {
		case <-ctx.Done():
			return nil
		case isActive, ok := <-l.active:
			// if chan is closed
			if !ok {
				return nil
			}
			// screen is inactive
			if !isActive {
				continue
			}
			l.btnHandlersReg()
			if err := l.sendRenderTask(); err != nil {
				logrus.Error(err)
			}
		}
	}
}

// GetFromStorage returns data from storage asserted to log message object.
func (l *LogMsg) GetFromStorage() (*entity.LogMsg, error) {
	storageData := l.storage.Get(repo.RendererKey)
	msg, ok := storageData.(*entity.LogMsg)
	if !ok {
		return nil, fmt.Errorf("%s: storage value is not log message object", l.serviceName)
	}
	return msg, nil
}

// BtnUpDefault represents a default up button handler for log message screen.
func (l *LogMsg) BtnUpDefault() error {
	menuInst, err := l.GetFromStorage()
	if err != nil {
		return err
	}
	// update log message and publish into storage as renderer
	menuInst.ScrollUp()
	l.storage.Publish(repo.RendererKey, menuInst)
	return nil
}

// BtnDownDefault represents a default down button handler for log message screen.
func (l *LogMsg) BtnDownDefault() error {
	msg, err := l.GetFromStorage()
	if err != nil {
		return err
	}
	// update log message and publish into storage as renderer
	msg.ScrollDown()
	l.storage.Publish(repo.RendererKey, msg)
	return nil
}

// sendRenderTask publishes log message to storage as renderer for render service.
func (l *LogMsg) sendRenderTask() error {
	// update msg after screen activation
	msg, err := l.logMsgGetter()
	if err != nil {
		return fmt.Errorf("%s: update log message: %w", l.serviceName, err)
	}
	// publish log message for render service
	l.storage.Publish(repo.RendererKey, msg)
	return nil
}
