package storage

import (
	"fmt"

	"IvolgaOledManager/internal/app/entity"
	"IvolgaOledManager/internal/pkg/errlog"
	"IvolgaOledManager/internal/pkg/storage"
)

var _ StationResultRepoStorage = (*stationResultRepoStorage)(nil)

// StationResultRepoStorage implementation.
type stationResultRepoStorage struct {
	store storage.Storage
}

func NewStationResultStorage(store storage.Storage) StationResultRepoStorage {
	return &stationResultRepoStorage{
		store: store,
	}
}

// GetMain gets station-result struct from storage.
func (s *stationResultRepoStorage) Get() *entity.StationResult {
	return s.get(_valueStationResult)
}

// SetMain sets station-result struct to storage.
func (s *stationResultRepoStorage) Set(value *entity.StationResult) {
	s.set(_valueStationResult, value)
}

// get gets station-result struct from storage.
func (s *stationResultRepoStorage) get(key string) *entity.StationResult {
	menu, ok := s.store.Get(key).(*entity.StationResult)
	if !ok {
		err := fmt.Errorf("station result value type is %T (%#v), not *entity.StationResult", menu, menu)
		errlog.Print(err)
		return &entity.StationResult{}
	}
	return menu
}

// set sets new station-result struct to storage.
func (s *stationResultRepoStorage) set(key string, value *entity.StationResult) {
	s.store.Set(key, value)
}
