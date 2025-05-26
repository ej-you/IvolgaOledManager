// Package storage provides Storage interface for key-value storage interaction.
// It contains golang map implementation of this interface.
package storage

type Storage interface {
	Get(key string) string
	Set(key string, value string)
	Delete(key string)
}
