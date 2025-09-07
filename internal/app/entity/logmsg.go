// Package entity contains all app entities.
package entity

import (
	"fmt"
	"time"

	"IvolgaOledManager/internal/pkg/display"
	"IvolgaOledManager/internal/pkg/text"
)

const (
	LogLvlCountCtxKey = "logLvlCount" // key for log level count value in context
	LogMsgCtxKey      = "logMsg"      // key for log level count value in context

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

// Render implements display.Renderer. It renders log message on display.
func (l *LogMsg) Render(device display.Display) error {
	// create slice of text lines with menu title line
	textLines := make([]display.TextLine, 0, display.DefaultLinesAmount)
	textLines = append(textLines, display.NewDefaultTextLine(l.Datetime()))

	content := l.Formated()
	// iterate menu items
	for _, line := range content {
		if len(textLines) == display.DefaultLinesAmount {
			break
		}
		// add new line
		textLines = append(textLines, display.NewDefaultTextLine(line))
	}

	// output text screen
	err := device.DisplayTextLines(textLines...)
	if err != nil {
		return fmt.Errorf("display text screen: %w", err)
	}
	return nil
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
