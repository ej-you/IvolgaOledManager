// Package db contains interfaces of DB repositories for all entities and its implementations.
package db

type SensorsRepoDB interface {
	GetTemperature() (float64, error)
	GetHumidity() (float64, error)
	GetPressure() (int64, error)
	GetWindSpeed() (float64, error)
	GetWindDirection() (float64, error)
}
