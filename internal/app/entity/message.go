// Package entity contains all app entities.
package entity

import (
	"time"

	"sschmc/internal/pkg/text"
)

const (
	_maxDisplayedLines = 4  // max lines amount that can be displayed simultaneously
	_maxLineLen        = 18 // max len of line
)

// Log message model.
type Message struct {
	ID        string    `gorm:"primaryKey;autoIncrement;type:INT"`
	Level     string    `gorm:"not null;size:1"`
	Header    string    `gorm:"not null;size:50"`
	Text      string    `gorm:"not null;size:255"`
	CreatedAt time.Time `gorm:"not null;type:TIMESTAMP"`

	FirstLine int      `gorm:"-"` // idx of first displayed line on the device (default: 0)
	Lines     []string `gorm:"-"` // a slice of lines formatted for display
}

// Format create lines slice of message Text to display it on device as text lines.
func (m *Message) Format() {
	m.Lines = text.Normalize(m.Text, _maxLineLen)
}

// ScrollUp updates message FirstLine for scrolling up imitation.
func (m *Message) ScrollUp() {
	// extreme up position
	if m.FirstLine == 0 {
		return
	}
	m.FirstLine--
}

// ScrollDown updates message FirstLine for scrolling down imitation.
func (m *Message) ScrollDown() {
	// extreme down position
	if m.FirstLine >= len(m.Lines)-_maxDisplayedLines {
		return
	}
	m.FirstLine++
}
