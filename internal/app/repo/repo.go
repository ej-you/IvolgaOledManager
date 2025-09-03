// Package repo contains interfaces of repositories for any sources and entities.
package repo

import "IvolgaOledManager/internal/app/entity"

// SensdataRepoDB describes a DB repo for sensor data.
type SensdataRepoDB interface {
	// GetTemperature returns temperature value from DB.
	GetTemperature() (float64, error)
	// GetHumidity returns humidity value from DB.
	GetHumidity() (float64, error)
	// GetPressure returns pressure value from DB.
	GetPressure() (float64, error)
	// GetWindSpeed returns wind speed value from DB.
	GetWindSpeed() (float64, error)
	// GetWindDirection returns wind direction value from DB.
	GetWindDirection() (float64, error)
}

// MessageRepoDB describes a DB repo for log messages.
type MessageRepoDB interface {
	// GetLevelsCount returns map with level numbers and amount of messages with this level.
	GetLevelsCount() ([]entity.MessageLevelCount, error)
	// GetWithLevel returns slice of messages with given level ordered by created datetime.
	GetWithLevel(level string) ([]entity.MessageWithLevel, error)
	// GetByID returns message with given ID. ID field must be presented.
	GetByID(msg *entity.Message) error
	// DeleteByID deletes message record by its ID. ID field must be presented.
	DeleteByID(id string) error
	// DeleteAllWithLevel deletes all message records with given level.
	DeleteAllWithLevel(level string) error
}

// SensconfRepoFS describes a file repo for sensors' config.
type SensconfRepoFS interface {
	// ParseSensconf returns slice of station sensors.
	// It parse station config file with the next layout:
	// `_sensorSectionPrefix \n _openBracket ...[config-lines]... \n _closeBracket`.
	ParseSensconf() (entity.Sensconf, error)
	// UpdateSensconf updates sensor section of config file according to given sensors data.
	UpdateSensconf(sensors entity.Sensconf) error
}

// Precompiled string keys for pub/sub storage.
const (
	RendererKey      = "renderer"
	MenuSensconf     = "menu:sensconf"
	SensTempKey      = "sensdata:temperature"
	SensHumidKey     = "sensdata:humidity"
	SensPressKey     = "sensdata:pressure"
	SensWindSpeedKey = "sensdata:wind:speed"
	SensWindDirKey   = "sensdata:wind:direction"
)
