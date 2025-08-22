// Package pubsub contains key-value storage with pub/sub funcs.
package pubsub

import (
	"context"
	"sync"
)

var _ Storage = (*KeyValueStorage)(nil)

// Storage is aa pub/sub storage interface.
type Storage interface {
	Get(key string) any
	Publish(key string, val any)
	Subscribe(ctx context.Context, key string) <-chan struct{}
}

// KeyValueStorage is a key-value storage with pub/sub supporting.
// KeyValueStorage provides subscription on key.
// Subscriber receives notify message then the value of key
// (to which subscriber is subscribed) is updated.
type KeyValueStorage struct {
	// map with subscribers' chans for each of published key
	notifyMap map[string][]chan struct{}
	// key-value storage for published data
	storage map[string]any
	mu      sync.RWMutex
}

// NewKeyValueStorage returns new pub/sub storage.
func NewKeyValueStorage() *KeyValueStorage {
	return &KeyValueStorage{
		notifyMap: make(map[string][]chan struct{}),
		storage:   make(map[string]any),
	}
}

// Get returns value from storage by given key.
func (s *KeyValueStorage) Get(key string) any {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.storage[key]
}

// Publish saves key-value pair into storage and notify all key subscribers.
func (s *KeyValueStorage) Publish(key string, val any) {
	// set new value
	s.mu.Lock()
	s.storage[key] = val
	s.mu.Unlock()

	// send notify to subscribers
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, channel := range s.notifyMap[key] {
		channel <- struct{}{}
	}
}

// Subscribe creates and returns notify chan for given key.
func (s *KeyValueStorage) Subscribe(ctx context.Context, key string) <-chan struct{} {
	notify := make(chan struct{}, 1)
	s.addNotifyChanToMap(key, notify)

	// wait for context canceling and close notify chan
	go func() {
		defer close(notify)
		defer s.removeNotifyChanFromMap(key, notify)
		<-ctx.Done()
	}()
	return notify
}

// addNotifyChanToMap adds chan into slice of chans (in map) if
// slice is already init for given key.
// Else this method creates new slice of chans for given key and adds it into map.
func (s *KeyValueStorage) addNotifyChanToMap(key string, notify chan struct{}) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.notifyMap[key] = append(s.notifyMap[key], notify)
}

// removeNotifyChanFromMap removes chan from slice of chans (in map).
func (s *KeyValueStorage) removeNotifyChanFromMap(key string, notify chan struct{}) {
	s.mu.Lock()
	defer s.mu.Unlock()

	chanSlice, ok := s.notifyMap[key]
	if !ok {
		return
	}
	for idx := range chanSlice {
		if chanSlice[idx] == notify {
			s.notifyMap[key] = append(chanSlice[:idx], chanSlice[idx+1:]...)
			break
		}
	}
}
