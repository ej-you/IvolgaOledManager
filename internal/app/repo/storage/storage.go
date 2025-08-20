// Package storage provides interfaces with key-value storage for each entity.
// It contains storage manager with all interfaces.
package storage

import (
	"IvolgaOledManager/internal/app/entity"
	"IvolgaOledManager/internal/pkg/storage"
)

const (
	_keyAppStatus       = "app-status"     // key for app status
	_valueNone          = ""               // value for app status
	_valueGreetings     = "greetings"      // value for app status
	_valueStationResult = "station-result" // value for app status and key for station result struct
)

// AppRepoStorage contains general storage methods for app status.
type AppRepoStorage interface {
	SetNone()
	IsNone() bool

	SetGreetings()
	IsGreetings() bool

	SetStationResult()
	IsStationResult() bool
}

// StationResultsRepoStorage contains station result entity methods.
type StationResultsRepoStorage interface {
	Get() *entity.StationResults
	Set(value *entity.StationResults)
}

// RepoStorageManager contains all storage repos.
type RepoStorageManager struct {
	App            AppRepoStorage
	StationResults StationResultsRepoStorage
}

func NewRepoStorageManager(store storage.Storage) *RepoStorageManager {
	return &RepoStorageManager{
		App:            NewAppStorage(store),
		StationResults: NewStationResultsStorage(store),
	}
}
