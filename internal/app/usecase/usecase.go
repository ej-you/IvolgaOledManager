// Package usecase contains usecase interfaces and
// its implementations for entities.
package usecase

import "IvolgaOledManager/internal/app/entity"

// SensordataUsecase describes a usecase for sensor data.
type SensordataUsecase interface {
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

// SensorconfUsecase describes a usecase for sensors' config.
type SensorconfUsecase interface {
	// GetSensorconf parses sensors' config file and
	// returns if as slice of sensor config lines.
	GetSensorconf() (entity.Sensorconf, error)
	// UpdateSensorconf rewrite old sensors' config with new data.
	UpdateSensorconf(entity.Sensorconf) error
	// ToMenu translate sensors' config (parsed into slice) into menu for output.
	// Each sensor config line will be a separate menu item.
	ToMenu(data entity.Sensorconf) *entity.Menu
}
