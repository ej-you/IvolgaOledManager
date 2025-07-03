package storage

import (
	"errors"

	"IvolgaOledManager/internal/app/entity"
	"IvolgaOledManager/internal/pkg/errlog"
	"IvolgaOledManager/internal/pkg/storage"
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

// GetLevel gets menu-level struct from storage.
func (s *menuRepoStorage) GetLevel() *entity.Menu {
	return s.get(_valueMenuLevel)
}

// SetLevel sets menu-level struct to storage.
func (s *menuRepoStorage) SetLevel(value *entity.Menu) {
	s.set(_valueMenuLevel, value)
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
