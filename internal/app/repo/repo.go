// Package repo contains interfaces of repositories for any sources and entities.
package repo

import "context"

type SensorRepoDB interface {
	GetTemperature() (float64, error)
	GetHumidity() (float64, error)
	GetPressure() (float64, error)
	GetWindSpeed() (float64, error)
	GetWindDirection() (float64, error)
}

// PubSubKey is a type of string key for pub/sub storage.
type PubSubKey string

// Precompiled pub/sub keys.
var (
	RendererKey        PubSubKey = "renderer"
	SensorTempKey      PubSubKey = "sensor:temperature"
	SensorHumidKey     PubSubKey = "sensor:humidity"
	SensorPressKey     PubSubKey = "sensor:pressure"
	SensorWindSpeedKey PubSubKey = "sensor:wind:speed"
	SensorWindDirKey   PubSubKey = "sensor:wind:direction"
)

type PubSubStorage interface {
	Get(key PubSubKey) any
	Publish(key PubSubKey, val any)
	Subscribe(ctx context.Context, key PubSubKey) <-chan struct{}
}
