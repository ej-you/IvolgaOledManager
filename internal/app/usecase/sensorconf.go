package usecase

import (
	"context"
	"fmt"

	"IvolgaOledManager/internal/app/entity"
	"IvolgaOledManager/internal/app/repo"
)

// Ensure sensor data usecase implementats interface.
var _ SensorconfUsecase = (*SensorconfUC)(nil)

// SensorconfUC represents a usecase for sensors' config.
type SensorconfUC struct {
	sensorconfRepoFS repo.SensorconfRepoFS
}

// NewSensorconfUsecase returns a new instance of SensorconfUC.
func NewSensorconfUsecase(sensorconfRepoFS repo.SensorconfRepoFS) *SensorconfUC {
	return &SensorconfUC{
		sensorconfRepoFS: sensorconfRepoFS,
	}
}

// GetSensorconf parses sensors' config file and
// returns if as slice of sensor config lines.
func (s *SensorconfUC) GetSensorconf() (entity.Sensorconf, error) {
	data, err := s.sensorconfRepoFS.ParseSensors()
	if err != nil {
		return nil, fmt.Errorf("parse sensor conf: %w", err)
	}
	return data, nil
}

// UpdateSensorconf rewrite old sensors' config with new data.
func (s *SensorconfUC) UpdateSensorconf(data entity.Sensorconf) error {
	err := s.sensorconfRepoFS.UpdateSensors(data)
	if err != nil {
		return fmt.Errorf("update sensor conf: %w", err)
	}
	return nil
}

// ToMenu translate sensors' config (parsed into slice) into menu for output.
// Each sensor config line will be a separate menu item.
func (s *SensorconfUC) ToMenu(data entity.Sensorconf) *entity.Menu {
	// context with full sensorconf slice
	ctxWithData := context.WithValue(context.Background(), entity.SensorconfCtxKey, data)
	// collect menu items
	menuItems := make([]*entity.MenuItem, 0, len(data))
	for _, sensorconfItem := range data {
		menuItems = append(menuItems, entity.NewMenuItem(
			context.WithValue(ctxWithData, entity.SensorconfItemCtxKey, sensorconfItem),
			sensorconfItem.Name,
		))
	}
	return &entity.Menu{
		Title: "Датчики",
		Items: menuItems,
	}
}
