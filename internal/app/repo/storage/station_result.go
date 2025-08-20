package storage

import (
	"fmt"

	"IvolgaOledManager/internal/app/entity"
	"IvolgaOledManager/internal/pkg/errlog"
	"IvolgaOledManager/internal/pkg/storage"
)

var _ StationResultsRepoStorage = (*stationResultsRepoStorage)(nil)

// StationResultRepoStorage implementation.
type stationResultsRepoStorage struct {
	store storage.Storage
}

func NewStationResultsStorage(store storage.Storage) StationResultsRepoStorage {
	return &stationResultsRepoStorage{
		store: store,
	}
}

// GetMain gets station-result struct from storage.
func (s *stationResultsRepoStorage) Get() *entity.StationResults {
	return s.get(_valueStationResult)
}

// SetMain sets station-result struct to storage.
func (s *stationResultsRepoStorage) Set(value *entity.StationResults) {
	s.set(_valueStationResult, value)
}

// get gets station-result struct from storage.
func (s *stationResultsRepoStorage) get(key string) *entity.StationResults {
	menu, ok := s.store.Get(key).(*entity.StationResults)
	if !ok {
		err := fmt.Errorf("station result value type is %T (%#v), not *entity.StationResults", menu, menu)
		errlog.Print(err)
		return &entity.StationResults{}
	}
	return menu
}

// set sets new station-result struct to storage.
func (s *stationResultsRepoStorage) set(key string, value *entity.StationResults) {
	s.store.Set(key, value)
}
