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

// SetMenuMain sets current app-status to menu-main.
func (s *appRepoStorage) SetMenuMain() {
	s.setStatus(_valueMenuMain)
}

// IsMenuMain checks if the current app-status is menu-main.
func (s *appRepoStorage) IsMenuMain() bool {
	return s.getStatus() == _valueMenuMain
}

// SetMenuStation sets current app-status to menu-station.
func (s *appRepoStorage) SetMenuStation() {
	s.setStatus(_valueMenuStation)
}

// IsMenuLogs checks if the current app-status is menu-station.
func (s *appRepoStorage) IsMenuStation() bool {
	return s.getStatus() == _valueMenuStation
}

// SetSensor sets current app-status to sensor.
func (s *appRepoStorage) SetSensor() {
	s.setStatus(_valueSensor)

}

// IsMenuLogs checks if the current app-status is sensor.
func (s *appRepoStorage) IsSensor() bool {
	return s.getStatus() == _valueSensor
}

// SetMenuLogs sets current app-status to menu-logs.
func (s *appRepoStorage) SetMenuLogs() {
	s.setStatus(_valueMenuLogs)
}

// IsMenuLogs checks if the current app-status is menu-logs.
func (s *appRepoStorage) IsMenuLogs() bool {
	return s.getStatus() == _valueMenuLogs
}

// SetMenuLevel sets current app-status to menu-level.
func (s *appRepoStorage) SetMenuLevel() {
	s.setStatus(_valueMenuLevel)
}

// IsMenuLevel checks if the current app-status is menu-level.
func (s *appRepoStorage) IsMenuLevel() bool {
	return s.getStatus() == _valueMenuLevel
}

// SetMessage sets current app-status to message.
func (s *appRepoStorage) SetMessage() {
	s.setStatus(_valueMessage)
}

// IsMessage checks if the current app-status is message.
func (s *appRepoStorage) IsMessage() bool {
	return s.getStatus() == _valueMessage
}

// IsMenuAny checks if the current app-status is any menu.
func (s *appRepoStorage) IsMenuAny() bool {
	status := s.getStatus()
	switch status {
	case _valueMenuMain, _valueMenuLevel:
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
