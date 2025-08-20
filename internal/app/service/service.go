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

type SensorsDataUpdate struct {
	updateDuration time.Duration
	store          *storage.RepoStorageManager
	stationUC      usecase.SensorsUsecase
	render         chan<- struct{}
}

func NewSensorsDataUpdate(store *storage.RepoStorageManager,
	stationUC usecase.SensorsUsecase, render chan<- struct{}) *SensorsDataUpdate {

	return &SensorsDataUpdate{
		updateDuration: _updateDuration,
		store:          store,
		stationUC:      stationUC,
		render:         render,
	}
}

func (t *SensorsDataUpdate) StartWithShutdown(ctx context.Context) {
	// set init empty temperature value
	t.store.StationResults.Set(&entity.StationResults{
		Results: []entity.StationResult{
			{Title: "Температура", ResultText: "------"},
			{Title: "Влажность", ResultText: "------"},
			{Title: "Скорость ветра", ResultText: "------"},
			{Title: "Направление ветра", ResultText: "------"},
			{Title: "Давление", ResultText: "------"},
		},
		CurrentResult: 0,
	})

	// init ticker for temperature periodically updates
	ticker := time.NewTicker(t.updateDuration)
	defer ticker.Stop()

	var newResults *entity.StationResults
	var err error
	for {
		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
			newResults, err = t.stationUC.GetAllResults()
			if err != nil {
				errlog.Print(err)
				continue
			}
			t.store.StationResults.Set(newResults)

			// update render if current screen is temperature result
			if t.store.App.IsStationResult() {
				t.render <- struct{}{}
			}
		}
	}
}
