// Package storage provides interfaces with key-value storage for each entity.
// It contains storage manager with all interfaces.
package storage

import (
	"sschmc/internal/app/entity"
	"sschmc/internal/pkg/storage"
)

// AppStorage contains general storage methods.
type AppStorage interface {
	GetStatus() string
	SetStatus(string)
	StatusIsMenu() bool
}

// MenuStorage contains menu entity methods.
type MenuStorage interface {
	Get(key string) *entity.Menu
	Set(key string, value *entity.Menu)
}

type StorageManager struct {
	App  AppStorage
	Menu MenuStorage
}

func NewStorageManager(store storage.Storage) *StorageManager {
	return &StorageManager{
		App:  NewAppStorage(store),
		Menu: NewMenuStorage(store),
	}
}
