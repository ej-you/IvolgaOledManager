// Package db contains interfaces of DB repositories for all entities and its implementations.
package db

import (
	"sschmc/internal/app/entity"
)

type MessageRepoDB interface {
	GetLevelsCount() ([]entity.MessageLevelCount, error)
	GetWithLevel(level string) ([]entity.MessageWithLevel, error)
	GetByID(msg *entity.Message) error

	DeleteByID(msg *entity.Message) error
	DeleteAllWithLevel(level string) error
}
