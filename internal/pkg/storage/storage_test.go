package storage

import (
	"testing"
)

var store = NewMap()

func TestSetStorageValue(t *testing.T) {
	t.Log("Set new value")

	store.Set("hello", "world")
}

func TestGetStorageValue(t *testing.T) {
	t.Log("Get existing value")

	v := store.Get("hello")
	t.Log("Value:", v)
}

func TestUpdateStorageValue(t *testing.T) {
	t.Log("Set existing value and get it back")

	store.Set("hello", "fred")

	v := store.Get("hello")
	t.Log("Value:", v)
}
