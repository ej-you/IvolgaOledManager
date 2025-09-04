package usecase

import (
	"context"
	"errors"
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

// GetAsMenu parses sensors' config file and
// returns if as slice of sensor config lines.
func (s *SensconfUC) GetAsMenu() (*entity.Menu, error) {
	data, err := s.sensorconfRepoFS.ParseSensconf()
	if err != nil {
		return nil, fmt.Errorf("parse sensor conf: %w", err)
	}
	return sensconfToMenu(data), nil
}

// UpdateAsMenu rewrite old sensors' config with new data.
func (s *SensconfUC) UpdateAsMenu(menu *entity.Menu) error {
	sensconf, err := menuToSensconf(menu)
	if err != nil {
		return fmt.Errorf("convert menu to sensorconf: %w", err)
	}

	if err := s.sensorconfRepoFS.UpdateSensconf(sensconf); err != nil {
		return fmt.Errorf("update sensor conf: %w", err)
	}
	return nil
}

// sensconfToMenu translates sensconf object to menu representation.
func sensconfToMenu(sensconf entity.Sensconf) *entity.Menu {
	// collect menu items
	menuItems := make([]*entity.MenuItem, 0, len(sensconf))
	for _, sensorconfItem := range sensconf {
		menuItems = append(menuItems, entity.NewMenuItem(
			context.WithValue(context.Background(), entity.SensconfItemCtxKey, sensorconfItem),
			sensorconfItem.Name,
		))
	}
	return &entity.Menu{
		Title: "Датчики",
		Items: menuItems,
	}
}

// menuToSensconf translates menu representation to sensconf object.
func menuToSensconf(menu *entity.Menu) (entity.Sensconf, error) {
	sensconf := make(entity.Sensconf, 0, len(menu.Items))

	var ok bool
	var sensconfItem *entity.SensconfItem
	// extract sensconf items from every menu item
	for _, menuItem := range menu.Items {
		sensconfItem, ok = menuItem.Ctx.Value(entity.SensconfItemCtxKey).(*entity.SensconfItem)
		if !ok {
			return nil, errors.New("menu item has not sensconf item value")
		}
		sensconf = append(sensconf, sensconfItem)
	}
	return sensconf, nil
}
