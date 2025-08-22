package usecase

import "IvolgaOledManager/internal/app/entity"

type SensorDataUsecase interface {
	GetTemperature() (*entity.SensorData, error)
	GetHumidity() (*entity.SensorData, error)
	GetPressure() (*entity.SensorData, error)
	GetWindSpeed() (*entity.SensorData, error)
	GetWindDirection() (*entity.SensorData, error)
}
