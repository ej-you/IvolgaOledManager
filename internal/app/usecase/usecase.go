package usecase

import "IvolgaOledManager/internal/app/entity"

type StationResultUsecase interface {
	GetTemperature() (*entity.StationResult, error)
}
