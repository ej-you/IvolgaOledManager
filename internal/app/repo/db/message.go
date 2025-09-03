package db

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"IvolgaOledManager/internal/app/entity"
	"IvolgaOledManager/internal/app/repo"
)

// Ensure message repo implementats interface.
var _ repo.MessageRepoDB = (*MessageRepo)(nil)

// MessageRepo represents a DB repo for logs messages.
type MessageRepo struct {
	dbStorage *gorm.DB
}

// NewMessageRepoDB returns a new instance of MessageRepo.
func NewMessageRepoDB(dbStorage *gorm.DB) *MessageRepo {
	return &MessageRepo{
		dbStorage: dbStorage,
	}
}

// GetLevelsCount returns map with level numbers and amount of messages with this level.
func (r *MessageRepo) GetLevelsCount() ([]entity.MessageLevelCount, error) {
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
func (r *MessageRepo) GetWithLevel(level string) ([]entity.MessageWithLevel, error) {
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
func (r *MessageRepo) GetByID(msg *entity.Message) error {
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
func (r *MessageRepo) DeleteByID(id string) error {
	err := r.dbStorage.Delete(&entity.Message{}, id).Error
	if err != nil {
		return fmt.Errorf("delete by id: %w", err)
	}
	return nil
}

// DeleteAllWithLevel deletes all message records with given level.
func (r *MessageRepo) DeleteAllWithLevel(level string) error {
	err := r.dbStorage.Delete(&entity.Message{}, "level = ?", level).Error
	if err != nil {
		return fmt.Errorf("delete with level: %w", err)
	}
	return nil
}
