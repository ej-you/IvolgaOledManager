package buttons

// screenStationResult sets "station-result" app-status and update render.
func (b *Buttons) screenStationResult() {
	b.store.App.SetStationResult()
	b.render <- struct{}{}
}
