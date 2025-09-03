package entity

import (
	"fmt"

	"IvolgaOledManager/internal/pkg/display"
)

const _menuLinesAmount = 4 // amount of menu lines on screen (menu title and items)

// Menu represents any menu.
type Menu struct {
	Title        string      // menu title
	Items        []*MenuItem // menu items
	FirstItem    int         // idx of first displayed item on the device (default: 0)
	SelectedItem int         // idx of selected item on the device (default: 0)
}

// Render implements display.Renderer. It renders menu on display.
func (m *Menu) Render(device display.Display) error {
	// create slice of text lines with menu title line
	textLines := make([]display.TextLine, 0, _menuLinesAmount)
	textLines = append(textLines, display.TextLine{Content: m.Title, RelativeSize: 1})

	var lineContent string
	// iterate menu items
	for idx, menuItem := range m.Items {
		if len(textLines) == _menuLinesAmount {
			break
		}
		// skip items that are before first visible item
		if idx < m.FirstItem {
			continue
		}
		if idx == m.SelectedItem {
			lineContent = menuItem.ToOutputSelected()
			menuItem.Scroll() // scroll item if needed
		} else {
			lineContent = menuItem.ToOutput()
		}
		// add new item
		textLines = append(textLines, display.TextLine{
			Content:      lineContent,
			RelativeSize: 1,
		})
	}
	// append empty lines if items are not enough
	if len(textLines) < _menuLinesAmount {
		for range _menuLinesAmount - len(textLines) {
			textLines = append(textLines, display.TextLine{
				Content:      "",
				RelativeSize: 1,
			})
		}
	}
	// output text screen
	err := device.DisplayTextLines(textLines...)
	if err != nil {
		return fmt.Errorf("display text screen: %w", err)
	}
	return nil
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
