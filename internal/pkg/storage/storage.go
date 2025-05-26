// Package storage provides Storage interface for key-value storage interaction.
// It contains golang map implementation of this interface.
package storage

type Storage interface {
	Get(key string) any
	Set(key string, value any)
	Delete(key string)
}
