// Package repo contains interfaces of repositories for any sources and entities.
package repo

type SensorRepoDB interface {
	GetTemperature() (float64, error)
	GetHumidity() (float64, error)
	GetPressure() (float64, error)
	GetWindSpeed() (float64, error)
	GetWindDirection() (float64, error)
}

// Precompiled string keys for pub/sub storage.
const (
	RendererKey      = "renderer"
	SensTempKey      = "sensor:temperature"
	SensHumidKey     = "sensor:humidity"
	SensPressKey     = "sensor:pressure"
	SensWindSpeedKey = "sensor:wind:speed"
	SensWindDirKey   = "sensor:wind:direction"
)
