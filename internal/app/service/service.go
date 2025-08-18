package service

import (
	"IvolgaOledManager/internal/app/entity"
	"IvolgaOledManager/internal/app/repo/storage"
	"IvolgaOledManager/internal/app/usecase"
	"IvolgaOledManager/internal/pkg/errlog"
	"context"
	"time"
)

const _updateDuration = 5 * time.Second // duration for temperature updates

type TemperatureUpdate struct {
	updateDuration time.Duration
	store          *storage.RepoStorageManager
	stationUC      usecase.StationResultUsecase
}

func NewTemperatureUpdate(store *storage.RepoStorageManager,
	stationUC usecase.StationResultUsecase) *TemperatureUpdate {

	return &TemperatureUpdate{
		updateDuration: _updateDuration,
		store:          store,
		stationUC:      stationUC,
	}
}

func (t *TemperatureUpdate) StartWithShutdown(ctx context.Context) {
	// set init empty temperature value
	t.store.StationResult.Set(&entity.StationResult{
		Title:      "Температура",
		ResultText: "------",
	})

	// init ticker for temperature periodically updates
	ticker := time.NewTicker(t.updateDuration)
	defer ticker.Stop()

	var newTemp *entity.StationResult
	var err error
	for {
		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
			newTemp, err = t.stationUC.GetTemperature()
			if err != nil {
				errlog.Print(err)
				continue
			}
			t.store.StationResult.Set(newTemp)
		}
	}
}
