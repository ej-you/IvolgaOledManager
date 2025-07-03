// Package file contains interfaces of File repositories for all entities and its implementations.
package file

import (
	"IvolgaOledManager/internal/app/entity"
)

type StationRepoFile interface {
	ParseSensors() (entity.StationSensors, error)
	UpdateSensors(sensors entity.StationSensors) error
}
