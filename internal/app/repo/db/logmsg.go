package db

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"IvolgaOledManager/internal/app/entity"
	"IvolgaOledManager/internal/app/repo"
)

// Ensure message repo implementats interface.
var _ repo.LogMsgRepoDB = (*LogMsgRepo)(nil)

// LogMsgRepo represents a DB repo for logs messages.
type LogMsgRepo struct {
	dbStorage *gorm.DB
}

// NewLogMsgRepoDB returns a new instance of LogMsgRepo.
func NewLogMsgRepoDB(dbStorage *gorm.DB) *LogMsgRepo {
	return &LogMsgRepo{
		dbStorage: dbStorage,
	}
}

// GetLevelsCount returns map with level numbers and amount of messages with this level.
func (r *LogMsgRepo) GetLevelsCount() ([]entity.LogMsgLevelCount, error) {
	var results []entity.LogMsgLevelCount
	err := r.dbStorage.
		Model(&entity.LogMsg{}).
		Select("level, count(1) as count").
		Group("level").
		Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("get levels count: %w", err)
	}
	return results, nil
}

// GetWithLevel returns slice of messages with given level ordered by created datetime.
func (r *LogMsgRepo) GetWithLevel(level string) ([]entity.LogMsgWithLevel, error) {
	var results []entity.LogMsgWithLevel
	err := r.dbStorage.
		Model(&entity.LogMsg{}).
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
func (r *LogMsgRepo) GetByID(msg *entity.LogMsg) error {
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
func (r *LogMsgRepo) DeleteByID(id string) error {
	err := r.dbStorage.Delete(&entity.LogMsg{}, id).Error
	if err != nil {
		return fmt.Errorf("delete by id: %w", err)
	}
	return nil
}

// DeleteAllWithLevel deletes all message records with given level.
func (r *LogMsgRepo) DeleteAllWithLevel(level string) error {
	err := r.dbStorage.Delete(&entity.LogMsg{}, "level = ?", level).Error
	if err != nil {
		return fmt.Errorf("delete with level: %w", err)
	}
	return nil
}
