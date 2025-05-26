// Package entity contains all app entities.
package entity

// Log message model.
type Message struct {
	ID     string `gorm:"primaryKey;autoIncrement;type:INT"`
	Level  string `gorm:"not null;size:1"`
	Header string `gorm:"not null;size:50"`
	Text   string `gorm:"not null;size:255"`
}
