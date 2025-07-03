package storage

import (
	"errors"

	"IvolgaOledManager/internal/app/entity"
	"IvolgaOledManager/internal/pkg/errlog"
	"IvolgaOledManager/internal/pkg/storage"
)

var _ SensorRepoStorage = (*sensorRepoStorage)(nil)

// SensorRepoStorage implementation.
type sensorRepoStorage struct {
	store storage.Storage
}

func NewSensorStorage(store storage.Storage) SensorRepoStorage {
	return &sensorRepoStorage{
		store: store,
	}
}

// GetAll gets station sensors slice from storage.
func (s *sensorRepoStorage) GetAll() entity.StationSensors {
	sensors, ok := s.store.Get(_valueSensors).(entity.StationSensors)
	if !ok {
		errlog.Print(errors.New("sensor value is not entity.StationSensors"))
		return nil
	}
	return sensors
}

// SetAll sets station sensors slice to storage.
func (s *sensorRepoStorage) SetAll(value entity.StationSensors) {
	s.store.Set(_valueSensors, value)
}

// Get gets station sensor from storage.
func (s *sensorRepoStorage) Get() *entity.StationSensor {
	sensor, ok := s.store.Get(_valueSensor).(*entity.StationSensor)
	if !ok {
		errlog.Print(errors.New("sensor value is not *entity.StationSensor"))
		return &entity.StationSensor{}
	}
	return sensor
}

// Set sets station sensor to storage.
func (s *sensorRepoStorage) Set(value *entity.StationSensor) {
	s.store.Set(_valueSensor, value)
}
