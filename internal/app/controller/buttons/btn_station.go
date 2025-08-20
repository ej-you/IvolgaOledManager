package buttons

// screenStationResult sets "station-result" app-status and update render.
func (b *Buttons) screenStationResult() {
	b.store.App.SetStationResult()
	b.render <- struct{}{}
}

// nextStationResult sets "station-result" app-status,
// set next result index and update render.
func (b *Buttons) nextStationResult() {
	// set next result index
	sensorsResults := b.store.StationResults.Get()
	sensorsResults.CurrentResult++
	if sensorsResults.CurrentResult == len(sensorsResults.Results) {
		sensorsResults.CurrentResult = 0
	}
	b.store.StationResults.Set(sensorsResults)

	b.store.App.SetStationResult()
	b.render <- struct{}{}
}

// previousStationResult sets "station-result" app-status,
// set previous result index and update render.
func (b *Buttons) previousStationResult() {
	// set previous
	sensorsResults := b.store.StationResults.Get()
	sensorsResults.CurrentResult--
	if sensorsResults.CurrentResult < 0 {
		sensorsResults.CurrentResult = len(sensorsResults.Results) - 1
	}
	b.store.StationResults.Set(sensorsResults)

	b.store.App.SetStationResult()
	b.render <- struct{}{}
}
