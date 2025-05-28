package storage

import (
	"errors"

	"sschmc/internal/app/entity"
	"sschmc/internal/pkg/errlog"
	"sschmc/internal/pkg/storage"
)

var _ MenuRepoStorage = (*menuRepoStorage)(nil)

// MenuStorage implementation.
type menuRepoStorage struct {
	store storage.Storage
}

func NewMenuStorage(store storage.Storage) MenuRepoStorage {
	return &menuRepoStorage{
		store: store,
	}
}

// GetMain gets menu-main struct from storage.
func (s *menuRepoStorage) GetMain() *entity.Menu {
	return s.get(_valueMenuMain)
}

// SetMain sets menu-main struct to storage.
func (s *menuRepoStorage) SetMain(value *entity.Menu) {
	s.set(_valueMenuMain, value)
}

// get gets menu struct from storage.
func (s *menuRepoStorage) get(key string) *entity.Menu {
	menu, ok := s.store.Get(key).(*entity.Menu)
	if !ok {
		errlog.Print(errors.New("menu value is not *entity.Menu"))
		return &entity.Menu{}
	}
	return menu
}

// set sets new menu struct to storage.
func (s *menuRepoStorage) set(key string, value *entity.Menu) {
	s.store.Set(key, value)
}
