package usecase

import (
	"fmt"

	"IvolgaOledManager/internal/app/entity"
	"IvolgaOledManager/internal/app/repo"
	"IvolgaOledManager/internal/pkg/pubsub"
)

// Ensure sensor data usecase implementats interface.
var _ SensdataUsecase = (*SensdataUC)(nil)

var (
	// map with user-friendly sensor titles
	SensTitleMap = map[string]string{
		repo.SensTempKey:      "Температура",
		repo.SensHumidKey:     "Влажность",
		repo.SensPressKey:     "Давление",
		repo.SensWindSpeedKey: "Скорость ветра",
		repo.SensWindDirKey:   "Направление ветра",
	}
	// map with formats for sensors data values
	_sensFmtMap = map[string]string{
		repo.SensTempKey:      "%.2f°C",
		repo.SensHumidKey:     "%.2f%%",
		repo.SensPressKey:     "%.0f гПа",
		repo.SensWindSpeedKey: "%.2f м/с",
		repo.SensWindDirKey:   "%.2f°",
	}
)

// SensdataUC represents a usecase for sensor data.
type SensdataUC struct {
	sensorRepoDB repo.SensdataRepoDB
	storage      pubsub.Storage
}

// NewSensdataUsecase returns a new instance of SensdataUC.
func NewSensdataUsecase(sensorsRepoDB repo.SensdataRepoDB,
	storage pubsub.Storage) *SensdataUC {

	return &SensdataUC{
		sensorRepoDB: sensorsRepoDB,
		storage:      storage,
	}
}

// GetTemperature returns last temperature value from DB.
// If value was not found it returns nil.
func (u *SensdataUC) GetTemperature() (*entity.Sensdata, error) {
	return getSensorData(repo.SensTempKey, u.sensorRepoDB.GetTemperature)
}

// GetHumidity returns last humidity value from DB.
// If value was not found it returns nil.
func (u *SensdataUC) GetHumidity() (*entity.Sensdata, error) {
	return getSensorData(repo.SensHumidKey, u.sensorRepoDB.GetHumidity)
}

// GetPressure returns last pressure value from DB.
// If value was not found it returns nil.
func (u *SensdataUC) GetPressure() (*entity.Sensdata, error) {
	return getSensorData(repo.SensPressKey, u.sensorRepoDB.GetPressure)
}

// GetWindSpeed returns last wind speed value from DB.
// If value was not found it returns nil.
func (u *SensdataUC) GetWindSpeed() (*entity.Sensdata, error) {
	return getSensorData(repo.SensWindSpeedKey, u.sensorRepoDB.GetWindSpeed)
}

// GetWindDirection returns last wind direction value from DB.
// If value was not found it returns nil.
func (u *SensdataUC) GetWindDirection() (*entity.Sensdata, error) {
	return getSensorData(repo.SensWindDirKey, u.sensorRepoDB.GetWindDirection)
}

// getDataFunc is a function to get new sensor data value from DB.
type getDataFunc func() (float64, error)

// getSensorData returns last sensor data value from DB.
// If value was not found it returns nil.
func getSensorData(sensorKey string, getData getDataFunc) (*entity.Sensdata, error) {
	// get new value from DB
	newData, err := getData()
	if err != nil {
		return nil, fmt.Errorf("get %s: %w", sensorKey, err)
	}
	// return nil if value is not found
	if newData == 0.0 {
		return nil, nil
	}
	return &entity.Sensdata{
		Title:   SensTitleMap[sensorKey],
		Data:    fmt.Sprintf(_sensFmtMap[sensorKey], newData),
		RawData: newData,
	}, nil
}
