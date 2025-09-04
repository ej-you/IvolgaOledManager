package menu

import "IvolgaOledManager/internal/pkg/display"

type Menu interface {
	// renderer
	display.Renderer

	// getters
	GetTitle() string
	GetItems() []Item
	GetFstItem() int
	GetSelItem() int

	// behavior
	SelectPrevious()
	SelectNext()
}

type Item interface {
	// getters
	GetTitle() string
	GetValue(key any) any
	ToOutput() string
	ToOutputSelected() string

	// behavior
	Scroll()
}
