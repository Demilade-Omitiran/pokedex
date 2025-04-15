package main

import (
	"internal/pokeapi"
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	cache := pokeapi.NewCache(5 * time.Second)
	key := "testKey"
	value := []byte("testValue")

	cache.Add(key, value)

	cachedValue, found := cache.Get(key)
	if !found {
		t.Errorf("Expected to find key %s in cache", key)
	}

	if string(cachedValue) != string(value) {
		t.Errorf("Expected cached value %s, got %s", value, cachedValue)
	}

	time.Sleep(6 * time.Second)

	_, found = cache.Get(key)

	if found {
		t.Errorf("Expected key %s to be expired from cache", key)
	}
}
