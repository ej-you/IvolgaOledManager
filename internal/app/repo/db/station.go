package db

import (
	"fmt"

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
func (r *repoDB) GetTemperature() (float64, error) {
	var result float64
	err := r.dbStorage.
		Table("data").
		Select("param_value").
		Where(`table_name="temperatures"`).
		Order("data_id DESC").
		Limit(1).
		Scan(&result).Error
	if err != nil {
		return 0.0, fmt.Errorf("from db: %w", err)
	}
	return result, nil
}
