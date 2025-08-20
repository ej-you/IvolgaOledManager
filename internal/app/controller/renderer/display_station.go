package renderer

import (
	"IvolgaOledManager/internal/app/entity"
	"IvolgaOledManager/internal/pkg/drawer"
	"fmt"
)

// station renders station result.
func (r *Renderer) station(stationResults *entity.StationResults) error {
	curRes := stationResults.Results[stationResults.CurrentResult]
	// collect text screen
	textScreenBuilder := drawer.NewTextScreenBuilder(r.device.ScreenWidth, r.device.ScreenHeight)
	textScreen, err := textScreenBuilder.Build(
		drawer.TextLine{Content: curRes.Title, RelativeSize: 2},
		drawer.TextLine{Content: "", RelativeSize: 1},
		drawer.TextLine{Content: curRes.ResultText, RelativeSize: 4},
		drawer.TextLine{Content: "", RelativeSize: 1},
	)
	if err != nil {
		return fmt.Errorf("create text screen: %w", err)
	}
	// output text screen
	if err := r.device.DrawTextScreen(textScreen); err != nil {
		return fmt.Errorf("display text screen: %w", err)
	}
	return nil
}
