package storage

import (
	"errors"

	"sschmc/internal/app/entity"
	"sschmc/internal/pkg/errlog"
	"sschmc/internal/pkg/storage"
)

var _ MenuStorage = (*menuStorage)(nil)

// MenuStorage implementation.
type menuStorage struct {
	store storage.Storage
}

func NewMenuStorage(store storage.Storage) MenuStorage {
	return &menuStorage{
		store: store,
	}
}

func (s *menuStorage) Get(key string) *entity.Menu {
	menu, ok := s.store.Get(key).(*entity.Menu)
	if !ok {
		errlog.Print(errors.New("menu value is not *entity.Menu"))
	}
	return menu
}

func (s *menuStorage) Set(key string, value *entity.Menu) {
	s.store.Set(key, value)
}
