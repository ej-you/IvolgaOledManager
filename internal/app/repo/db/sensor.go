// Package db contains implementations of DB repository interfaces.
package db

import (
	"fmt"

	"gorm.io/gorm"

	"IvolgaOledManager/internal/app/repo"
)

// Ensure sensors repo implements interface.
var _ repo.SensorRepoDB = (*SensorsRepo)(nil)

// SensorsRepo is a repo for sensor data.
// It gets data values from DB.
type SensorsRepo struct {
	dbStorage *gorm.DB
}

// NewSensorsRepoDB returns a new instance of SensorsRepo.
func NewSensorsRepoDB(dbStorage *gorm.DB) *SensorsRepo {
	return &SensorsRepo{
		dbStorage: dbStorage,
	}
}

// GetTemperature returns temperature value from DB.
func (r *SensorsRepo) GetTemperature() (float64, error) {
	return r.getSensorResult(`table_name="temperatures"`)
}

// GetHumidity returns humidity value from DB.
func (r *SensorsRepo) GetHumidity() (float64, error) {
	return r.getSensorResult(`table_name="humiditys"`)
}

// GetPressure returns pressure value from DB.
func (r *SensorsRepo) GetPressure() (float64, error) {
	return r.getSensorResult(`table_name="abs_pressures"`)
}

// GetWindSpeed returns wind speed value from DB.
func (r *SensorsRepo) GetWindSpeed() (float64, error) {
	return r.getSensorResult(`table_name="winds" AND param_name="speed"`)
}

// GetWindDirection returns wind direction value from DB.
func (r *SensorsRepo) GetWindDirection() (float64, error) {
	return r.getSensorResult(`table_name="winds" AND param_name="direction"`)
}

// getSensorResult returns sensor result from DB with given condition.
func (r *SensorsRepo) getSensorResult(cond string) (float64, error) {
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
