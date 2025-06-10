package buttons

// screenNone sets "none" app-status and update render.
func (b *Buttons) screenNone() {
	b.store.App.SetNone()
	b.render <- struct{}{}
}

// screenGreetings sets "greetings" app-status and update render.
func (b *Buttons) screenGreetings() {
	b.store.App.SetGreetings()
	b.render <- struct{}{}
}
