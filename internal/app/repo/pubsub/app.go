// Package pubsub contains implementation of pub/sub storage.
package pubsub

import (
	"IvolgaOledManager/internal/app/repo"
	"IvolgaOledManager/internal/pkg/pubsub"
	"context"
)

var _ repo.PubSubStorage = (*AppPubSub)(nil)

// AppPubSub is a repo.PubSubStorage adaptor with
// precompiled pub/sub storage keys.
type AppPubSub struct {
	storage pubsub.PubSubStorage
}

func NewAppPubSub() *AppPubSub {
	return &AppPubSub{
		storage: *pubsub.NewPubSubStorage(),
	}
}

func (p *AppPubSub) Get(key repo.PubSubKey) any {
	return p.storage.Get(string(key))
}

func (p *AppPubSub) Publish(key repo.PubSubKey, val any) {
	p.storage.Publish(string(key), val)
}

func (p *AppPubSub) Subscribe(ctx context.Context, key repo.PubSubKey) <-chan struct{} {
	return p.storage.Subscribe(ctx, string(key))
}
