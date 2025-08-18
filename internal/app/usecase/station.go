package usecase

import (
	"fmt"

	"IvolgaOledManager/internal/app/entity"
	"IvolgaOledManager/internal/app/repo/db"
	"IvolgaOledManager/internal/app/repo/storage"
)

var _ StationResultUsecase = (*StationResultUC)(nil)

// StationResultUsecase implementation.
type StationResultUC struct {
	stationResultRepoDB      db.StationResultRepoDB
	stationResultRepoStorage storage.StationResultRepoStorage
}

func NewStationResultUsecase(stationResultRepoDB db.StationResultRepoDB,
	stationResultRepoStorage storage.StationResultRepoStorage) StationResultUsecase {

	return &StationResultUC{
		stationResultRepoDB:      stationResultRepoDB,
		stationResultRepoStorage: stationResultRepoStorage,
	}
}

// GetTemperature returns last temperature value from DB.
// If value was not found it returns last temperature value saved into storage.
func (u *StationResultUC) GetTemperature() (*entity.StationResult, error) {
	// get new value from DB
	newTemp, err := u.stationResultRepoDB.GetTemperature()
	if err != nil {
		return nil, fmt.Errorf("get temperature: %w", err)
	}
	// return value if exists
	if newTemp != 0.0 {
		return &entity.StationResult{
			Title:      "Температура",
			ResultText: fmt.Sprintf("%.2f°C", newTemp),
		}, err
	}

	// return value from storage
	return u.stationResultRepoStorage.Get(), nil
}
