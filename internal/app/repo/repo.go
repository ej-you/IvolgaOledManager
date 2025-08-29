// Package repo contains interfaces of repositories for any sources and entities.
package repo

// SensorsRepo describes a repo for sensor data.
type SensorRepoDB interface {
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

// Precompiled string keys for pub/sub storage.
const (
	RendererKey      = "renderer"
	SensTempKey      = "sensordata:temperature"
	SensHumidKey     = "sensordata:humidity"
	SensPressKey     = "sensordata:pressure"
	SensWindSpeedKey = "sensordata:wind:speed"
	SensWindDirKey   = "sensordata:wind:direction"
)
