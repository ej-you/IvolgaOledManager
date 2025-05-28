package storage

import (
	"errors"

	"sschmc/internal/app/entity"
	"sschmc/internal/pkg/errlog"
	"sschmc/internal/pkg/storage"
)

var _ MessageRepoStorage = (*messageRepoStorage)(nil)

// MessageRepoStorage implementation.
type messageRepoStorage struct {
	store storage.Storage
}

func NewMessageStorage(store storage.Storage) MessageRepoStorage {
	return &messageRepoStorage{
		store: store,
	}
}

// Get gets message struct from storage.
func (s *messageRepoStorage) Get() *entity.Message {
	message, ok := s.store.Get(_valueMessage).(*entity.Message)
	if !ok {
		errlog.Print(errors.New("message value is not *entity.Message"))
		return &entity.Message{}
	}
	return message
}

// Set sets message struct to storage.
func (s *messageRepoStorage) Set(value *entity.Message) {
	s.store.Set(_valueMessage, value)
}
