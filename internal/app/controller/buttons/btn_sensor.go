package buttons

import (
	"errors"

	"sschmc/internal/app/entity"
	"sschmc/internal/pkg/errlog"
	"sschmc/internal/pkg/system"
)

// screenSensor sets "sensor" app-status and update render.
func (b *Buttons) screenSensor() {
	stationMenu := b.store.Menu.GetStation()
	selectedItem := stationMenu.Items[stationMenu.SelectedItem]
	// get selected sensor from station menu
	selectedSensor, ok := selectedItem.Value.(*entity.StationSensor)
	if !ok {
		errlog.Print(errors.New("sensor value is not *entity.StationSensor"))
		return
	}

	b.store.Sensor.Set(selectedSensor)
	b.store.App.SetSensor()
	b.render <- struct{}{}
}

// updateStation change selected sensor status, updates
// station config file and restart station service.
func (b *Buttons) updateStation() {
	// get all sensors and selected sensor
	stationSensors := b.store.Sensor.GetAll()
	sensor := b.store.Sensor.Get()
	// update selected sensor status and station config
	stationSensors[sensor.Idx].ChangeActive()
	if err := b.stationRepoFile.UpdateSensors(stationSensors); err != nil {
		errlog.Print(err)
	}
	// restart station service
	if err := system.RestartService(b.stationService); err != nil {
		errlog.Print(err)
	}
	b.render <- struct{}{}
}
