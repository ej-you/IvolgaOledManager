package usecase

import (
	"fmt"

	"IvolgaOledManager/internal/app/entity"
	"IvolgaOledManager/internal/app/repo"
)

var _ SensorDataUsecase = (*SensorDataUC)(nil)

var (
	SensorTitleTemp = "Температура"
	_sensorFmtTemp  = "%.2f°C"

	SensorTitleHumid = "Влажность"
	_sensorFmtHumid  = "%.2f%%"

	SensorTitlePress = "Давление"
	_sensorFmtPress  = "%.0f гПа"

	SensorTitleWindSpeed = "Скорость ветра"
	_sensorFmtWindSpeed  = "%.2f м/с"

	SensorTitleWindDir = "Направление ветра"
	_sensorFmtWindDir  = "%.2f°"
)

// SensorDataUC is a SensorDataUsecase implementation.
type SensorDataUC struct {
	sensorRepoDB repo.SensorRepoDB
	storage      repo.PubSubStorage
}

func NewSensorDataUsecase(sensorsRepoDB repo.SensorRepoDB,
	storage repo.PubSubStorage) SensorDataUsecase {

	return &SensorDataUC{
		sensorRepoDB: sensorsRepoDB,
		storage:      storage,
	}
}

// GetTemperature returns last temperature value from DB.
// If value was not found it returns nil.
func (u *SensorDataUC) GetTemperature() (*entity.SensorData, error) {
	// get new value from DB
	newData, err := u.sensorRepoDB.GetTemperature()
	if err != nil {
		return nil, fmt.Errorf("get temperature: %w", err)
	}
	// return nil if value is not found
	if newData == 0.0 {
		return nil, nil
	}
	return &entity.SensorData{
		Title:   SensorTitleTemp,
		Data:    fmt.Sprintf(_sensorFmtTemp, newData),
		RawData: newData,
	}, nil
}

// GetHumidity returns last humidity value from DB.
// If value was not found it returns nil.
func (u *SensorDataUC) GetHumidity() (*entity.SensorData, error) {
	// get new value from DB
	newData, err := u.sensorRepoDB.GetHumidity()
	if err != nil {
		return nil, fmt.Errorf("get humidity: %w", err)
	}
	// return nil if value is not found
	if newData == 0.0 {
		return nil, nil
	}
	return &entity.SensorData{
		Title:   SensorTitleHumid,
		Data:    fmt.Sprintf(_sensorFmtHumid, newData),
		RawData: newData,
	}, nil
}

// GetPressure returns last pressure value from DB.
// If value was not found it returns nil.
func (u *SensorDataUC) GetPressure() (*entity.SensorData, error) {
	// get new value from DB
	newData, err := u.sensorRepoDB.GetPressure()
	if err != nil {
		return nil, fmt.Errorf("get pressure: %w", err)
	}
	// return nil if value is not found
	if newData == 0.0 {
		return nil, nil
	}
	return &entity.SensorData{
		Title:   SensorTitlePress,
		Data:    fmt.Sprintf(_sensorFmtPress, newData),
		RawData: newData,
	}, nil
}

// GetWindSpeed returns last wind speed value from DB.
// If value was not found it returns nil.
func (u *SensorDataUC) GetWindSpeed() (*entity.SensorData, error) {
	// get new value from DB
	newData, err := u.sensorRepoDB.GetWindSpeed()
	if err != nil {
		return nil, fmt.Errorf("get wind speed: %w", err)
	}
	// return nil if value is not found
	if newData == 0.0 {
		return nil, nil
	}
	return &entity.SensorData{
		Title:   SensorTitleWindSpeed,
		Data:    fmt.Sprintf(_sensorFmtWindSpeed, newData),
		RawData: newData,
	}, nil
}

// GetWindDirection returns last wind direction value from DB.
// If value was not found it returns nil.
func (u *SensorDataUC) GetWindDirection() (*entity.SensorData, error) {
	// get new value from DB
	newData, err := u.sensorRepoDB.GetWindDirection()
	if err != nil {
		return nil, fmt.Errorf("get wind direction: %w", err)
	}
	// return nil if value is not found
	if newData == 0.0 {
		return nil, nil
	}
	return &entity.SensorData{
		Title:   SensorTitleWindDir,
		Data:    fmt.Sprintf(_sensorFmtWindDir, newData),
		RawData: newData,
	}, nil
}
