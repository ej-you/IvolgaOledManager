package buttons

// btnAllGreetings renders greetings and updates app-status in storage.
func (b *Buttons) btnAllGreetings() {
	b.store.App.SetGreetings()
	// update render according to new app-status
	b.render <- struct{}{}
}
