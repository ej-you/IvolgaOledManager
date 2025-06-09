package renderer

import (
	"fmt"

	"sschmc/internal/app/entity"
)

// sensor renders sensor info.
func (r *Renderer) sensor(stationSensor *entity.StationSensor) error {
	drawer, err := r.device.NewTextDrawer()
	if err != nil {
		return fmt.Errorf("create text drawer: %w", err)
	}

	drawer.AddLine("", "Sensor info")
	drawer.AddLine("", stationSensor.Name)

	if stationSensor.Active {
		drawer.AddLine("", "включён (on)")
		drawer.AddLine(">  ", "отключить")
	} else {
		drawer.AddLine("", "выключен (off)")
		drawer.AddLine(">  ", "включить")
	}
	drawer.FillEmpty()

	if err := drawer.Draw(); err != nil {
		return fmt.Errorf("display text lines: %w", err)
	}
	return nil
}
