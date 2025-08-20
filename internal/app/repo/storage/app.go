package storage

import (
	"errors"

	"IvolgaOledManager/internal/pkg/errlog"
	"IvolgaOledManager/internal/pkg/storage"
)

var _ AppRepoStorage = (*appRepoStorage)(nil)

// AppStorage implementation.
type appRepoStorage struct {
	store storage.Storage
}

func NewAppStorage(store storage.Storage) AppRepoStorage {
	// set none status to avoid error
	store.Set(_keyAppStatus, "")

	return &appRepoStorage{
		store: store,
	}
}

// SetNone sets current app-status to none.
func (s *appRepoStorage) SetNone() {
	s.setStatus(_valueNone)
}

// IsNone checks if the current app-status is none.
func (s *appRepoStorage) IsNone() bool {
	return s.getStatus() == _valueNone
}

// SetGreetings sets current app-status to greetings.
func (s *appRepoStorage) SetGreetings() {
	s.setStatus(_valueGreetings)
}

// IsGreetings checks if the current app-status is greetings.
func (s *appRepoStorage) IsGreetings() bool {
	return s.getStatus() == _valueGreetings
}

// SetStationResult sets current app-status to station-result.
func (s *appRepoStorage) SetStationResult() {
	s.setStatus(_valueStationResult)
}

// IsStationResult checks if the current app-status is station-result.
func (s *appRepoStorage) IsStationResult() bool {
	return s.getStatus() == _valueStationResult
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
