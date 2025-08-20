package usecase

import "IvolgaOledManager/internal/app/entity"

type SensorsUsecase interface {
	GetAllResults() (*entity.StationResults, error)
}
