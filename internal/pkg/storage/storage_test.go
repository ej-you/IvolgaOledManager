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

func TestDeleteStorageValue(t *testing.T) {
	t.Log("Delete existing value and try to get it")

	store.Delete("hello")

	v := store.Get("hello")
	if v != nil {
		t.Error("Unexpected value gotten. Value:", v)
	}
	t.Log("No value")
}
