package entity

const (
	DefaultPrefix  = "  "
	SelectedPrefix = "> "

	_maxItemLen = 16 // max item value len
	_maxLines   = 4  // max lines amount for display (ssc-hmc/internal/pkg/ssd1306/text_drawer.go)
)

// Menu item
type MenuItem struct {
	Title        string // item title
	FirstSymbol  int    // idx of first displayed symbol of the item (default: 0)
	ScrollToLeft bool   // true if item text must be scrolled back to left (default: false)
}

// Scroll updates item state for running line imitation
// if its value len is more than _maxItemLen.
func (i *MenuItem) Scroll() {
	// skip for items with short value len
	if len(i.Title) <= _maxItemLen {
		return
	}
	// extreme right position
	if !i.ScrollToLeft && (_maxItemLen+i.FirstSymbol == len(i.Title)) {
		i.ScrollToLeft = !i.ScrollToLeft
		i.FirstSymbol--
		return
	}
	// extreme left position
	if i.ScrollToLeft && (i.FirstSymbol == 0) {
		i.ScrollToLeft = !i.ScrollToLeft
		i.FirstSymbol++
		return
	}
	// scroll to right
	if i.ScrollToLeft {
		i.FirstSymbol--
	} else {
		i.FirstSymbol++
	}
}

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
		m.FirstItem = max(len(m.Items)-_maxLines, 0)
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
	if m.SelectedItem >= m.FirstItem+_maxLines {
		m.FirstItem++
	}
}
