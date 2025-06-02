package entity

const (
	DefaultPrefix  = "   " // menu item prefix for default item
	SelectedPrefix = "> "  // menu item prefix for selected item

	// _maxLines = 4 // max lines amount for display (ssc-hmc/internal/pkg/ssd1306/text_drawer.go)
)

// Menu model.
type Menu struct {
	Title        string      // menu title
	Items        []*MenuItem // menu items
	FirstItem    int         // idx of first displayed item on the device (default: 0)
	SelectedItem int         // idx of selected item on the device (default: 0)
}

// SelectPrevious updates menu state for scrolling up imitation.
func (m *Menu) SelectPrevious() {
	// extreme up position
	if m.SelectedItem == 0 {
		m.SelectedItem = len(m.Items) - 1
		m.FirstItem = max(len(m.Items)-MaxDisplayedItems, 0)
		return
	}

	m.SelectedItem--
	// if selected item after update will not be visible
	if m.SelectedItem < m.FirstItem {
		m.FirstItem--
	}
}

// SelectNext updates menu state for scrolling down imitation.
func (m *Menu) SelectNext() {
	// extreme down position
	if m.SelectedItem == len(m.Items)-1 {
		m.SelectedItem = 0
		m.FirstItem = 0
		return
	}

	m.SelectedItem++
	// if selected item after update will not be visible
	if m.SelectedItem >= m.FirstItem+MaxDisplayedItems {
		m.FirstItem++
	}
}
