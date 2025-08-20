package usecase

import (
	"fmt"
	"log"

	"IvolgaOledManager/internal/app/entity"
	"IvolgaOledManager/internal/app/repo/db"
	"IvolgaOledManager/internal/app/repo/storage"
)

var _ SensorsUsecase = (*SensorsUC)(nil)

// SensorsUsecase implementation.
type SensorsUC struct {
	sensorsRepoDB            db.SensorsRepoDB
	stationResultRepoStorage storage.StationResultsRepoStorage
}

func NewStationResultUsecase(sensorsRepoDB db.SensorsRepoDB,
	stationResultRepoStorage storage.StationResultsRepoStorage) SensorsUsecase {

	return &SensorsUC{
		sensorsRepoDB:            sensorsRepoDB,
		stationResultRepoStorage: stationResultRepoStorage,
	}
}

// GetAllResults returns last all sensors values from DB.
func (u *SensorsUC) GetAllResults() (*entity.StationResults, error) {
	// get new values from DB
	newTemp, err := u.sensorsRepoDB.GetTemperature()
	if err != nil {
		log.Printf("ERROR: get temperature: %v", err)
	}
	newHumid, err := u.sensorsRepoDB.GetHumidity()
	if err != nil {
		log.Printf("ERROR: get humidity: %v", err)
	}
	newPres, err := u.sensorsRepoDB.GetPressure()
	if err != nil {
		log.Printf("ERROR: get pressure: %v", err)
	}
	newWindSpeed, err := u.sensorsRepoDB.GetWindSpeed()
	if err != nil {
		log.Printf("ERROR: get wind speed: %v", err)
	}
	newWindDirection, err := u.sensorsRepoDB.GetWindDirection()
	if err != nil {
		log.Printf("ERROR: get wind direction: %v", err)
	}

	oldResults := u.stationResultRepoStorage.Get()

	values := []float64{newTemp, newHumid, newWindSpeed, newWindDirection}
	formats := []string{"%.2f°C", "%.2f%%", "%.2f м/с", "%.2f°", "%d гПа"}

	newResults := &entity.StationResults{
		Results: []entity.StationResult{
			{Title: "Температура"},
			{Title: "Влажность"},
			{Title: "Скорость ветра"},
			{Title: "Направление ветра"},
			{Title: "Давление", ResultText: oldResults.Results[4].ResultText},
		},
		CurrentResult: oldResults.CurrentResult,
	}
	var resultText string
	for idx, value := range values {
		resultText = oldResults.Results[idx].ResultText
		if value != 0.0 {
			resultText = fmt.Sprintf(formats[idx], value)
		}
		newResults.Results[idx].ResultText = resultText
	}
	if newPres != 0 {
		newResults.Results[4].ResultText = fmt.Sprintf(formats[4], newPres)
	}

	if !compareResults(oldResults, newResults) {
		log.Printf("Got new sensors data: %+v", newResults.Results)
	}
	return newResults, nil
}

// compareResults returns true if old and new results are equal and
// returns false if old and new results are not equal.
func compareResults(old, new *entity.StationResults) bool {
	if old.CurrentResult != new.CurrentResult || len(old.Results) != len(new.Results) {
		return false
	}
	for idx := range old.Results {
		if old.Results[idx] != new.Results[idx] {
			return false
		}
	}
	return true
}
