// Package usecase contains usecase interfaces and
// its implementations for entities.
package usecase

import "IvolgaOledManager/internal/app/entity"

// SensorDataUC describes a usecase for sensor data.
type SensorDataUsecase interface {
	// GetTemperature returns last temperature value from DB.
	// If value was not found it returns nil.
	GetTemperature() (*entity.SensorData, error)
	// GetHumidity returns last humidity value from DB.
	// If value was not found it returns nil.
	GetHumidity() (*entity.SensorData, error)
	// GetPressure returns last pressure value from DB.
	// If value was not found it returns nil.
	GetPressure() (*entity.SensorData, error)
	// GetWindSpeed returns last wind speed value from DB.
	// If value was not found it returns nil.
	GetWindSpeed() (*entity.SensorData, error)
	// GetWindDirection returns last wind direction value from DB.
	// If value was not found it returns nil.
	GetWindDirection() (*entity.SensorData, error)
}
