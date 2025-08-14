package db

import (
	"gorm.io/gorm"
)

var _ StationResultRepoDB = (*repoDB)(nil)

// StationResultRepoDB implementation.
type repoDB struct {
	dbStorage *gorm.DB
}

func NewStationResultRepoDB(dbStorage *gorm.DB) StationResultRepoDB {
	return &repoDB{
		dbStorage: dbStorage,
	}
}

// GetTemperature returns temperature value from DB.
// TODO: find out DB conn info and rewrite method
func (r *repoDB) GetTemperature() (float64, error) {
	panic("not defined")
	// var results []entity.MessageLevelCount
	// err := r.dbStorage.
	// 	Model(&entity.Message{}).
	// 	Select("level, count(1) as count").
	// 	Group("level").
	// 	Find(&results).Error
	// if err != nil {
	// 	return nil, fmt.Errorf("get levels count: %w", err)
	// }
	// return results, nil
}
