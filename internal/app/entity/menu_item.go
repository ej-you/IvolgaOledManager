package entity

import (
	"unicode/utf8"
)

const (
	_maxItemLen = 15 // max item value len

	_menuItemDeleteButtonText = "УДАЛИТЬ ВСЕ"
)

// Menu item
type MenuItem struct {
	Title string // item title
	Value any    // item value is used for DB

	firstSymbol  int  // idx of first displayed symbol of the item (default: 0)
	lastSymbol   int  // idx of last displayed symbol of the item (use NewMenuItem to set up)
	scrollToLeft bool // true if item text must be scrolled back to left (default: false)
	skipScroll   bool // true if length of item title is less than _maxItemLen (default: false)
	deleteButton bool // true if item is deleteWithLevel button (default: false)
}

// NewMenuItem returns new prepared menu item. It's not recommended to init menu item directly.
func NewMenuItem(title string, value any) *MenuItem {
	menuItem := &MenuItem{
		Title: title,
		Value: value,
	}
	menuItem.setLastSymbol()
	return menuItem
}

// setLastSymbol computes and sets lastSymbol value for menu item.
func (i *MenuItem) setLastSymbol() {
	// if length of item title is less than _maxItemLen
	if utf8.RuneCountInString(i.Title) <= _maxItemLen {
		i.lastSymbol = len(i.Title)
		i.skipScroll = true
		return
	}
	// loop from 0 to _maxItemLen rune of item title
	var runeIdx int
	for runeIdx = range i.Title {
		if utf8.RuneCountInString(i.Title[:runeIdx]) == _maxItemLen {
			break
		}
	}
	i.lastSymbol = runeIdx
}

// moveRight updates menu item values to show item title for one symbol to the right.
func (i *MenuItem) moveRight() {
	i.firstSymbol += utf8.RuneLen(rune(i.Title[i.firstSymbol]))
	i.lastSymbol += utf8.RuneLen(rune(i.Title[i.lastSymbol]))
}

// moveLeft updates menu item values to show item title for one symbol to the left.
func (i *MenuItem) moveLeft() {
	i.firstSymbol -= utf8.RuneLen(rune(i.Title[i.firstSymbol-1]))
	i.lastSymbol -= utf8.RuneLen(rune(i.Title[i.lastSymbol-1]))
}

// Scroll updates item state for running line imitation
// if its value len is more than _maxItemLen.
func (i *MenuItem) Scroll() {
	if i.skipScroll {
		return
	}
	// extreme right position then moving right
	if !i.scrollToLeft && (i.lastSymbol == len(i.Title)) {
		i.scrollToLeft = !i.scrollToLeft
		i.moveLeft()
		return
	}
	// extreme left position then moving left
	if i.scrollToLeft && (i.firstSymbol == 0) {
		i.scrollToLeft = !i.scrollToLeft
		i.moveRight()
		return
	}
	if i.scrollToLeft {
		i.moveLeft()
		return
	}
	i.moveRight()
}

// FormattedTitle returns substring of the item title to print it out
// according to the current item scroll status.
func (i *MenuItem) FormattedTitle() string {
	return i.Title[i.firstSymbol:i.lastSymbol]
}

// IsDeleteButton returns deleteButton value.
func (i *MenuItem) IsDeleteButton() bool {
	return i.deleteButton
}

// NewMenuItemDeleteButton returns menu item for button with any delete functionality meaning.
func NewMenuItemDeleteButton(value any) *MenuItem {
	menuItem := &MenuItem{
		Title:        _menuItemDeleteButtonText,
		Value:        value,
		deleteButton: true,
	}
	menuItem.setLastSymbol()
	return menuItem
}
