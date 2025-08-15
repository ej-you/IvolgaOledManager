// Package entity contains all app entities.
package entity

import (
	"unicode/utf8"
)

const (
	_defaultResultFontHeight = 16.0 // default height of result text font
	_maxLineLen              = 16   // max len of line

	MaxDisplayedItems = 3 // max menu items (message lines) amount that can be displayed
)

type StationResult struct {
	title            string
	resultText       string
	resultFontHeight float64
}

// NewStationResult creates new station result instance.
func NewStationResult(title, res string) *StationResult {
	inst := &StationResult{
		resultFontHeight: _defaultResultFontHeight,
	}

	// set title
	if utf8.RuneCountInString(title) > _maxLineLen {
		inst.title = title[:_maxLineLen-3] + "..."
	} else {
		inst.title = title
	}
	// set result text
	if utf8.RuneCountInString(res) > _maxLineLen {
		inst.resultText = res[:_maxLineLen-3] + "..."
	} else {
		inst.resultText = res
	}
	return inst
}

// Title is a title getter.
func (s StationResult) Title() string {
	return s.title
}

// ResultText is a result text getter.
func (s StationResult) ResultText() string {
	return s.resultText
}
