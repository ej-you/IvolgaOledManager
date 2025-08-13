package buttons

import (
	"fmt"
	"log"

	"IvolgaOledManager/internal/app/entity"
	"IvolgaOledManager/internal/pkg/errlog"
)

// screenStationResult sets "station-result" app-status and update render.
func (b *Buttons) screenStationResult() {
	// TODO: get temperature from DB
	temper, err := b.stationResultRepoDB.GetTemperature()
	if err != nil {
		errlog.Print(err)
		return
	}
	log.Printf("Gotten temperature: %+v", temper)
	statRes := entity.NewStationResult("Температура", fmt.Sprintf("%.2f", temper))

	b.store.StationResult.Set(statRes)
	b.store.App.SetStationResult()
	b.render <- struct{}{}
}
