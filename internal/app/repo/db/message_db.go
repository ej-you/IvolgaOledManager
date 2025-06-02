package db

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"sschmc/internal/app/entity"
)

var _ MessageRepoDB = (*repoDB)(nil)

// MessageRepoDB implementation.
type repoDB struct {
	dbStorage *gorm.DB
}

func NewMessageRepoDB(dbStorage *gorm.DB) MessageRepoDB {
	return &repoDB{
		dbStorage: dbStorage,
	}
}

// GetLevelsCount returns map with level numbers and amount of messages with this level.
func (r *repoDB) GetLevelsCount() ([]entity.MessageLevelCount, error) {
	var results []entity.MessageLevelCount
	err := r.dbStorage.
		Model(&entity.Message{}).
		Select("level, count(1) as count").
		Group("level").
		Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("get levels count: %w", err)
	}
	return results, nil
}

// GetWithLevel returns slice of messages with given level ordered by created datetime.
func (r *repoDB) GetWithLevel(level string) ([]entity.MessageWithLevel, error) {
	var results []entity.MessageWithLevel
	err := r.dbStorage.
		Model(&entity.Message{}).
		Where("level = ?", level).
		Order("created_at DESC").
		Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("get with level: %w", err)
	}
	return results, nil
}

// GetByID returns message with given ID.
// Field ID must be presented.
func (r *repoDB) GetByID(msg *entity.Message) error {
	err := r.dbStorage.Where("id = ?", msg.ID).First(&msg).Error
	// not found error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		msg.Header = "db error"
		msg.Content = "log message not found"
	}
	if err != nil {
		return fmt.Errorf("get by id: %w", err)
	}
	return nil
}

// DeleteByID deletes message record by its ID.
// ID field must be presented.
func (r *repoDB) DeleteByID(msg *entity.Message) error {
	err := r.dbStorage.Delete(&entity.Message{}, msg.ID).Error
	if err != nil {
		return fmt.Errorf("delete by id: %w", err)
	}
	return nil
}

// DeleteAllWithLevel deletes all message records with given level.
func (r *repoDB) DeleteAllWithLevel(level string) error {
	err := r.dbStorage.Delete(&entity.Message{}, "level = ?", level).Error
	if err != nil {
		return fmt.Errorf("delete with level: %w", err)
	}
	return nil
}
