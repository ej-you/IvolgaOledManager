// Package sensordata provides services to update sensors data.
package sensordata

import (
	"context"
	"time"

	"IvolgaOledManager/config"
	"IvolgaOledManager/internal/app/entity"
	"IvolgaOledManager/internal/app/repo"
	"IvolgaOledManager/internal/app/usecase"
	"IvolgaOledManager/internal/pkg/pubsub"

	"github.com/sirupsen/logrus"
)

// getDataFunc is a function to get new sensor data value from usecase.
type getDataFunc func() (*entity.SensorData, error)

// Updater represents a sensor data updater.
type Updater struct {
	// will be closed if the service was completely started and is ready-to-use now
	ready chan struct{}

	// time between data update requests
	updatesDuration time.Duration
	// function from sensor usecase to get data from specific sensor
	getData getDataFunc

	// storage instance
	store pubsub.Storage
	// storage key for value of specific sensor
	storageKey string
}

// StartWithShutdown starts updater loop. It may be stopped by context cancellaiton.
func (s *Updater) StartWithShutdown(ctx context.Context) error {
	logrus.Infof("start %s service...", s.storageKey)

	// init ticker for sensor data periodically updates
	ticker := time.NewTicker(s.updatesDuration)
	defer ticker.Stop()
	// notify that service is ready-to-use
	close(s.ready)

	var newData *entity.SensorData
	var err error
	for {
		select {
		case <-ctx.Done():
			return nil

		case <-ticker.C:
			newData, err = s.getData()
			if err != nil {
				logrus.Errorf("update %s: %v", s.storageKey, err)
				continue
			}
			// if data is new
			if newData != nil {
				logrus.Infof("new data on %s", s.storageKey)
				s.store.Publish(s.storageKey, newData)
			}
		}
	}
}

// Ready signals that the service is ready-to-use.
func (s *Updater) Ready() <-chan struct{} {
	return s.ready
}

// NewTemperatureUpdater returns a new instance of temperature Updater.
func NewTemperatureUpdater(cfg *config.Config, store pubsub.Storage,
	sensorUC usecase.SensorDataUsecase) *Updater {

	return newUpdater(cfg, store, sensorUC.GetTemperature, repo.SensTempKey)
}

// NewHumidityUpdater returns a new instance of humidity Updater.
func NewHumidityUpdater(cfg *config.Config, store pubsub.Storage,
	sensorUC usecase.SensorDataUsecase) *Updater {

	return newUpdater(cfg, store, sensorUC.GetHumidity, repo.SensHumidKey)
}

// NewPressureUpdater returns a new instance of pressure Updater.
func NewPressureUpdater(cfg *config.Config, store pubsub.Storage,
	sensorUC usecase.SensorDataUsecase) *Updater {

	return newUpdater(cfg, store, sensorUC.GetPressure, repo.SensPressKey)
}

// NewWindSpeedUpdater returns a new instance of wind speed Updater.
func NewWindSpeedUpdater(cfg *config.Config, store pubsub.Storage,
	sensorUC usecase.SensorDataUsecase) *Updater {

	return newUpdater(cfg, store, sensorUC.GetWindSpeed, repo.SensWindSpeedKey)
}

// NewWindDirUpdater returns a new instance of wind direction Updater.
func NewWindDirUpdater(cfg *config.Config, store pubsub.Storage,
	sensorUC usecase.SensorDataUsecase) *Updater {

	return newUpdater(cfg, store, sensorUC.GetWindDirection, repo.SensWindDirKey)
}

// newUpdater returns new sensor data updater service.
func newUpdater(cfg *config.Config, store pubsub.Storage,
	getData getDataFunc, sensorKey string) *Updater {

	// publish empty sensor data into storage
	emptyData := &entity.SensorData{
		Title: usecase.SensorTitleMap[sensorKey],
		Data:  "------",
	}
	store.Publish(sensorKey, emptyData)
	// return update service
	return &Updater{
		ready:           make(chan struct{}),
		updatesDuration: cfg.Other.Sensors.DataUpdatesDuration,
		getData:         getData,
		store:           store,
		storageKey:      sensorKey,
	}
}
