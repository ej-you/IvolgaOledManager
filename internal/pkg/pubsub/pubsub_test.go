package pubsub

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPubSubStorage(t *testing.T) {
	storage := NewPubSubStorage()

	// publish data (no subscribers)
	data1Key := "test"
	data1 := "test string"
	storage.Publish(data1Key, data1)

	// get published data
	gottenData := storage.Get(data1Key)
	t.Logf("gottenData: %+v", gottenData)
	require.Equal(t, data1, gottenData)

	// subscribe on test key
	ctx, cancel := context.WithCancel(context.Background())
	notify := storage.Subscribe(ctx, data1Key)

	data2 := "test string notify"
	// wait for notify and print out updated data
	go func() {
		defer cancel()
		<-notify
		t.Log("Notified about data update")

		gottenData := storage.Get(data1Key)
		t.Logf("gottenData after notify: %+v", gottenData)
		require.Equal(t, data2, gottenData)
	}()

	storage.Publish(data1Key, data2)
	// waiting for context is done
	<-ctx.Done()
}
