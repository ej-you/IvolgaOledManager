package pubsub

import (
	"context"
	"sync"

	"github.com/sirupsen/logrus"
)

type PubSubStorage struct {
	// map with subscribers' chans for each of published key
	notifyMap map[string][]chan struct{}
	// key-value storage for published data
	storage map[string]any
	mu      sync.RWMutex
}

func NewPubSubStorage() *PubSubStorage {
	return &PubSubStorage{
		notifyMap: make(map[string][]chan struct{}),
		storage:   make(map[string]any),
	}
}

// Get returns value from storage by given key.
func (s *PubSubStorage) Get(key string) any {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.storage[key]
}

func (s *PubSubStorage) Publish(key string, val any) {
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
	logrus.Infof("Key %q subscribers: %d", key, len(s.notifyMap[key]))
}

// Subscribe creates and returns notify chan for given key.
func (s *PubSubStorage) Subscribe(ctx context.Context, key string) <-chan struct{} {
	notify := make(chan struct{}, 1)
	s.addNotifyChanToMap(key, notify)

	// wait for context canceling and close notify chan
	go func() {
		defer close(notify)
		defer s.removeNotifyChanFromMap(key, notify)
		<-ctx.Done()
		logrus.Info("AAAAAAAAA")
	}()
	logrus.Infof("DDDDDDDDD: %#v", s.notifyMap)
	return notify
}

// addNotifyChanToMap adds chan into slice of chans (in map) if
// slice is already init for given key.
// Else this method creates new slice of chans for given key and adds it into map.
func (s *PubSubStorage) addNotifyChanToMap(key string, notify chan struct{}) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.notifyMap[key] = append(s.notifyMap[key], notify)
	logrus.Infof("BBBBB: %#v", s.notifyMap)
}

// removeNotifyChanFromMap removes chan from slice of chans (in map).
func (s *PubSubStorage) removeNotifyChanFromMap(key string, notify chan struct{}) {
	s.mu.Lock()
	defer s.mu.Unlock()

	chanSlice, ok := s.notifyMap[key]
	if !ok {
		return
	}
	for idx := range chanSlice {
		logrus.Infof("chanSlice[idx]: %+v | notify: %+v", chanSlice[idx], notify)
		if chanSlice[idx] == notify {
			s.notifyMap[key] = append(chanSlice[:idx], chanSlice[idx+1:]...)
			break
		}
	}
	logrus.Infof("CCCCC: %#v", s.notifyMap)
}
