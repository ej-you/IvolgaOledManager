// Package usecase contains usecase interfaces and
// its implementations for entities.
package usecase

import "IvolgaOledManager/internal/app/entity"

// SensdataUsecase describes a usecase for sensor data.
type SensdataUsecase interface {
	// GetTemperature returns last temperature value from DB.
	// If value was not found it returns nil.
	GetTemperature() (*entity.Sensdata, error)
	// GetHumidity returns last humidity value from DB.
	// If value was not found it returns nil.
	GetHumidity() (*entity.Sensdata, error)
	// GetPressure returns last pressure value from DB.
	// If value was not found it returns nil.
	GetPressure() (*entity.Sensdata, error)
	// GetWindSpeed returns last wind speed value from DB.
	// If value was not found it returns nil.
	GetWindSpeed() (*entity.Sensdata, error)
	// GetWindDirection returns last wind direction value from DB.
	// If value was not found it returns nil.
	GetWindDirection() (*entity.Sensdata, error)
}

// SensconfUsecase describes a usecase for sensors' config.
type SensconfUsecase interface {
	// GetAsMenu parses sensors' config file into a slice of sensor config lines
	// and returns it as menu.
	GetAsMenu() (*entity.Menu, error)
	// UpdateAsMenu gets sensor config data from given menu and
	// rewrite old sensors' config with new data.
	// After the config is updated, it restarts the station service.
	UpdateAsMenu(menu *entity.Menu) error
}
