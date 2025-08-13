// Package db contains interfaces of DB repositories for all entities and its implementations.
package db

type StationResultRepoDB interface {
	GetTemperature() (float64, error)
}
