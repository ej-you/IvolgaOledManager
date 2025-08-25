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

type Updater struct {
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

// NewTemperatureUpdater returns new temperature updater service.
func NewTemperatureUpdater(cfg *config.Config, store pubsub.Storage,
	sensorUC usecase.SensorDataUsecase) *Updater {

	return newUpdater(cfg, store, sensorUC.GetTemperature, repo.SensTempKey)
}

// NewHumidityUpdater returns new humidity updater service.
func NewHumidityUpdater(cfg *config.Config, store pubsub.Storage,
	sensorUC usecase.SensorDataUsecase) *Updater {

	return newUpdater(cfg, store, sensorUC.GetHumidity, repo.SensHumidKey)
}

// NewPressureUpdater returns new pressure updater service.
func NewPressureUpdater(cfg *config.Config, store pubsub.Storage,
	sensorUC usecase.SensorDataUsecase) *Updater {

	return newUpdater(cfg, store, sensorUC.GetPressure, repo.SensPressKey)
}

// NewWindSpeedUpdater returns new wind speed updater service.
func NewWindSpeedUpdater(cfg *config.Config, store pubsub.Storage,
	sensorUC usecase.SensorDataUsecase) *Updater {

	return newUpdater(cfg, store, sensorUC.GetWindSpeed, repo.SensWindSpeedKey)
}

// NewWindDirUpdater returns new wind direction updater service.
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
		updatesDuration: cfg.Other.Sensors.DataUpdatesDuration,
		getData:         getData,
		store:           store,
		storageKey:      sensorKey,
	}
}
