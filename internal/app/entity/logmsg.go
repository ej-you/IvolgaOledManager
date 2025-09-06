// Package entity contains all app entities.
package entity

import (
	"time"

	"IvolgaOledManager/internal/pkg/text"
)

const (
	_datetimeFormat = "02.01.06 15:04:05" // datetime format for createdAt message field
	_lineSymbols    = 16                  // max amount of symbols on one line
)

// LogMsg represents log message object.
type LogMsg struct {
	ID        string    `gorm:"primaryKey;autoIncrement;type:INT"`
	Level     string    `gorm:"not null;size:1"`
	Header    string    `gorm:"not null;size:30"`
	Content   string    `gorm:"not null;size:255"`
	CreatedAt time.Time `gorm:"not null;type:TIMESTAMP"`

	FirstLine int      `gorm:"-"` // idx of first displayed line on the device (default: 0)
	lines     []string `gorm:"-"` // a slice of lines formatted for display
}

func (LogMsg) TableName() string {
	return "storage"
}

// Datetime returns formated createdAt message field.
func (l *LogMsg) Datetime() string {
	return l.CreatedAt.Format(_datetimeFormat)
}

// Formated creates lines slice of message Text to display
// it on device as text lines and returns it.
func (l *LogMsg) Formated() []string {
	// if already formated
	if len(l.lines) != 0 {
		return l.lines
	}
	// join header, content and delete button text to a single string
	fullText := l.Header + "\n\n" + l.Content
	l.lines = text.Normalize(fullText, _lineSymbols)
	return l.lines
}

// ScrollUp updates message FirstLine for scrolling up imitation.
func (l *LogMsg) ScrollUp() {
	// extreme up position
	if l.FirstLine == 0 {
		return
	}
	l.FirstLine--
}

// ScrollDown updates message FirstLine for scrolling down imitation.
func (l *LogMsg) ScrollDown() {
	// extreme down position
	if l.FirstLine >= len(l.lines)-_displayLinesBlue {
		return
	}
	l.FirstLine++
}

// LogMsgLevelCount is a subset of LogMsg object fields with
// all log levels and messages amount for each of log level.
type LogMsgLevelCount struct {
	Level int
	Count int
}

// LogMsgWithLevel is a subset of LogMsg object fields with
// messages and their levels.
type LogMsgWithLevel struct {
	ID     string
	Header string
}
