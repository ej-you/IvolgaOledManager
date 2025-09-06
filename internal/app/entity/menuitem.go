package entity

import (
	"context"
	"unicode/utf8"
)

const (
	_defaultPrefix  = "   " // menu item prefix for default item
	_selectedPrefix = "> "  // menu item prefix for selected item

	_menuItemSymbols = 15 // max amount of symbols on one line (without prefix)
)

// Menu item
type MenuItem struct {
	Title string          // item title
	Ctx   context.Context // context for additional values

	firstSymbol  int  // idx of first displayed symbol of the item (default: 0)
	lastSymbol   int  // idx of last displayed symbol of the item (use NewMenuItem to set up)
	scrollToLeft bool // true if item text must be scrolled back to left (default: false)
	skipScroll   bool // true if length of item title is less than _maxItemLen (default: false)
}

// NewMenuItem returns new prepared menu item. It's not recommended to init menu item directly.
func NewMenuItem(ctx context.Context, title string) *MenuItem {
	menuItem := &MenuItem{
		Title: title,
		Ctx:   ctx,
	}
	menuItem.setLastSymbol()
	return menuItem
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

// ToOutput returns substring of the item title to print it out
// according to the current item scroll status with default prefix (unselected item).
func (i *MenuItem) ToOutput() string {
	return _defaultPrefix + i.Title[i.firstSymbol:i.lastSymbol]
}

// ToOutputSelected returns substring of the item title to print it out
// according to the current item scroll status with selected prefix (selected item).
func (i *MenuItem) ToOutputSelected() string {
	return _selectedPrefix + i.Title[i.firstSymbol:i.lastSymbol]
}

// setLastSymbol computes and sets lastSymbol value for menu item.
func (i *MenuItem) setLastSymbol() {
	// if length of item title is less than _maxItemLen
	if utf8.RuneCountInString(i.Title) <= _menuItemSymbols {
		i.lastSymbol = len(i.Title)
		i.skipScroll = true
		return
	}
	// loop from 0 to _maxItemLen rune of item title
	var runeIdx int
	for runeIdx = range i.Title {
		if utf8.RuneCountInString(i.Title[:runeIdx]) == _menuItemSymbols {
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
