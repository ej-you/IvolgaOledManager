package usecase

import (
	"context"
	"fmt"

	"IvolgaOledManager/internal/app/entity"
	"IvolgaOledManager/internal/app/repo"
)

// Ensure sensors' config usecase implements interface.
var _ SensconfUsecase = (*SensconfUC)(nil)

// SensconfUC represents a usecase for sensors' config.
type SensconfUC struct {
	sensorconfRepoFS repo.SensconfRepoFS
}

// NewSensconfUsecase returns a new instance of SensconfUC.
func NewSensconfUsecase(sensorconfRepoFS repo.SensconfRepoFS) *SensconfUC {
	return &SensconfUC{
		sensorconfRepoFS: sensorconfRepoFS,
	}
}

// Get parses sensors' config file and
// returns if as slice of sensor config lines.
func (s *SensconfUC) Get() (entity.Sensconf, error) {
	data, err := s.sensorconfRepoFS.ParseSensconf()
	if err != nil {
		return nil, fmt.Errorf("parse sensor conf: %w", err)
	}
	return data, nil
}

// Update rewrite old sensors' config with new data.
func (s *SensconfUC) Update(data entity.Sensconf) error {
	err := s.sensorconfRepoFS.UpdateSensconf(data)
	if err != nil {
		return fmt.Errorf("update sensor conf: %w", err)
	}
	return nil
}

// ToMenu translate sensors' config (parsed into slice) into menu for output.
// Each sensor config line will be a separate menu item.
func (s *SensconfUC) ToMenu(data entity.Sensconf) *entity.Menu {
	// context with full sensorconf slice
	ctxWithData := context.WithValue(context.Background(), entity.SensconfCtxKey, data)
	// collect menu items
	menuItems := make([]*entity.MenuItem, 0, len(data))
	for _, sensorconfItem := range data {
		menuItems = append(menuItems, entity.NewMenuItem(
			context.WithValue(ctxWithData, entity.SensconfItemCtxKey, sensorconfItem),
			sensorconfItem.Name,
		))
	}
	return &entity.Menu{
		Title: "Датчики",
		Items: menuItems,
	}
}
