package db

import (
	"fmt"

	"gorm.io/gorm"
)

var _ SensorsRepoDB = (*repoDB)(nil)

// SensorsRepoDB implementation.
type repoDB struct {
	dbStorage *gorm.DB
}

func NewSensorsRepoDB(dbStorage *gorm.DB) SensorsRepoDB {
	return &repoDB{
		dbStorage: dbStorage,
	}
}

// GetTemperature returns temperature value from DB.
func (r *repoDB) GetTemperature() (float64, error) {
	return r.getSensorResult(`table_name="temperatures"`)
}

// GetHumidity returns humidity value from DB.
func (r *repoDB) GetHumidity() (float64, error) {
	return r.getSensorResult(`table_name="humiditys"`)
}

// GetPressure returns pressure value from DB.
func (r *repoDB) GetPressure() (int64, error) {
	floatValue, err := r.getSensorResult(`table_name="abs_pressures"`)
	return int64(floatValue), err
}

// GetWindSpeed returns wind speed value from DB.
func (r *repoDB) GetWindSpeed() (float64, error) {
	return r.getSensorResult(`table_name="winds" AND param_name="speed"`)
}

// GetWindDirection returns wind direction value from DB.
func (r *repoDB) GetWindDirection() (float64, error) {
	return r.getSensorResult(`table_name="winds" AND param_name="direction"`)
}

// getSensorResult returns sensor result from DB with given condition.
func (r *repoDB) getSensorResult(cond string) (float64, error) {
	var result float64
	err := r.dbStorage.
		Table("data").
		Select("param_value").
		Where(cond).
		Order("data_id DESC").
		Limit(1).
		Scan(&result).Error
	if err != nil {
		return 0.0, fmt.Errorf("from db: %w", err)
	}
	return result, nil
}
