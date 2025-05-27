package storage

import (
	"errors"
	"sschmc/internal/pkg/errlog"
	"sschmc/internal/pkg/storage"
)

var _ AppRepoStorage = (*appRepoStorage)(nil)

// AppStorage implementation.
type appRepoStorage struct {
	store storage.Storage
}

func NewAppStorage(store storage.Storage) AppRepoStorage {
	return &appRepoStorage{
		store: store,
	}
}

// SetNone sets current app-status to none.
func (s *appRepoStorage) SetNone() {
	s.setStatus(_valueNone)
}

// StatusIsNone checks if the current app-status is none.
func (s *appRepoStorage) IsNone() bool {
	return s.getStatus() == _valueNone
}

// SetGreetings sets current app-status to greetings.
func (s *appRepoStorage) SetGreetings() {
	s.setStatus(_valueGreetings)
}

// StatusIsGreetings checks if the current app-status is greetings.
func (s *appRepoStorage) IsGreetings() bool {
	return s.getStatus() == _valueGreetings
}

// SetMenuMain sets current app-status to menu-main.
func (s *appRepoStorage) SetMenuMain() {
	s.setStatus(_valueMenuMain)
}

// StatusIsMenuMain checks if the current app-status is menu-main.
func (s *appRepoStorage) IsMenuMain() bool {
	return s.getStatus() == _valueMenuMain
}

// StatusIsMenu checks if the current app-status is any menu.
func (s *appRepoStorage) IsMenuAny() bool {
	status := s.getStatus()
	switch status {
	case _valueMenuMain:
		return true
	default:
		return false
	}
}

// setStatus gets current app-status.
func (s *appRepoStorage) getStatus() string {
	status, ok := s.store.Get(_keyAppStatus).(string)
	if !ok {
		errlog.Print(errors.New("app-status value is not string"))
		return ""
	}
	return status
}

// setStatus sets new app-status.
func (s *appRepoStorage) setStatus(status string) {
	s.store.Set(_keyAppStatus, status)
}
