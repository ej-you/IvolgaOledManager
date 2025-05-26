package storage

import (
	"errors"
	"sschmc/internal/pkg/errlog"
	"sschmc/internal/pkg/storage"
)

const (
	_keyAppStatus            = "app-status"
	_valueAppStatusNone      = ""
	_valueAppStatusGreetings = "greetings"
	_valueAppStatusMenuMain  = "menu-main"
)

var _ AppStorage = (*appStorage)(nil)

// AppStorage implementation.
type appStorage struct {
	store storage.Storage
}

func NewAppStorage(store storage.Storage) AppStorage {
	return &appStorage{
		store: store,
	}
}

func (s *appStorage) GetStatus() string {
	status, ok := s.store.Get(_keyAppStatus).(string)
	if !ok {
		errlog.Print(errors.New("app-status value is not string"))
	}
	return status
}

func (s *appStorage) SetStatus(status string) {
	s.store.Set(_keyAppStatus, status)
}

// IsMenu checks if the current app-status is any menu.
func (s *appStorage) StatusIsMenu() bool {
	status := s.GetStatus()
	switch status {
	case _valueAppStatusMenuMain:
		return true
	default:
		return false
	}
}
