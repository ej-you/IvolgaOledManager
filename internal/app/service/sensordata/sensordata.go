package sensordata

import (
	"context"
	"time"

	"IvolgaOledManager/config"
	"IvolgaOledManager/internal/app/entity"
	"IvolgaOledManager/internal/app/repo"
	"IvolgaOledManager/internal/app/usecase"

	"github.com/sirupsen/logrus"
)

type SensorDataUpdate struct {
	// time between data update requests
	updatesDuration time.Duration
	// function from sensor usecase to get data from specific sensor
	getData func() (*entity.SensorData, error)

	// storage instance
	store repo.PubSubStorage
	// storage key for value of specific sensor
	storageKey repo.PubSubKey
}

func (s *SensorDataUpdate) StartWithShutdown(ctx context.Context) error {
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
				s.store.Publish(s.storageKey, newData)
			}
		}
	}
}

// NewTemperatureUpdate returns new temperature update service.
func NewTemperatureUpdate(cfg *config.Config, store repo.PubSubStorage,
	sensorUC usecase.SensorDataUsecase) *SensorDataUpdate {

	// publish empty sensor data into storage
	emptyData := &entity.SensorData{Title: usecase.SensorTitleTemp, Data: "------"}
	store.Publish(repo.SensorTempKey, emptyData)
	// return update service
	return &SensorDataUpdate{
		updatesDuration: cfg.Other.Sensors.DataUpdatesDuration,
		getData:         sensorUC.GetTemperature,
		store:           store,
		storageKey:      repo.SensorTempKey,
	}
}

// NewHumidityUpdate returns new humidity update service.
func NewHumidityUpdate(cfg *config.Config, store repo.PubSubStorage,
	sensorUC usecase.SensorDataUsecase) *SensorDataUpdate {

	// publish empty sensor data into storage
	emptyData := &entity.SensorData{Title: usecase.SensorTitleHumid, Data: "------"}
	store.Publish(repo.SensorHumidKey, emptyData)
	// return update service
	return &SensorDataUpdate{
		updatesDuration: cfg.Hardware.Oled.UpdatesDuration,
		getData:         sensorUC.GetHumidity,
		store:           store,
		storageKey:      repo.SensorHumidKey,
	}
}

// NewPressureUpdate returns new pressure update service.
func NewPressureUpdate(cfg *config.Config, store repo.PubSubStorage,
	sensorUC usecase.SensorDataUsecase) *SensorDataUpdate {

	// publish empty sensor data into storage
	emptyData := &entity.SensorData{Title: usecase.SensorTitlePress, Data: "------"}
	store.Publish(repo.SensorPressKey, emptyData)
	// return update service
	return &SensorDataUpdate{
		updatesDuration: cfg.Hardware.Oled.UpdatesDuration,
		getData:         sensorUC.GetPressure,
		store:           store,
		storageKey:      repo.SensorPressKey,
	}
}

// NewWindSpeedUpdate returns new wind speed update service.
func NewWindSpeedUpdate(cfg *config.Config, store repo.PubSubStorage,
	sensorUC usecase.SensorDataUsecase) *SensorDataUpdate {

	// publish empty sensor data into storage
	emptyData := &entity.SensorData{Title: usecase.SensorTitleWindSpeed, Data: "------"}
	store.Publish(repo.SensorWindSpeedKey, emptyData)
	// return update service
	return &SensorDataUpdate{
		updatesDuration: cfg.Hardware.Oled.UpdatesDuration,
		getData:         sensorUC.GetWindSpeed,
		store:           store,
		storageKey:      repo.SensorWindSpeedKey,
	}
}

// NewWindDirUpdate returns new wind direction update service.
func NewWindDirUpdate(cfg *config.Config, store repo.PubSubStorage,
	sensorUC usecase.SensorDataUsecase) *SensorDataUpdate {

	// publish empty sensor data into storage
	emptyData := &entity.SensorData{Title: usecase.SensorTitleWindDir, Data: "------"}
	store.Publish(repo.SensorWindDirKey, emptyData)
	// return update service
	return &SensorDataUpdate{
		updatesDuration: cfg.Hardware.Oled.UpdatesDuration,
		getData:         sensorUC.GetWindDirection,
		store:           store,
		storageKey:      repo.SensorWindDirKey,
	}
}
