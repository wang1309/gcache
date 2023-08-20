package gcache

import (
	"testing"
	"time"
)

func TestCache_Get(t *testing.T) {
	loader := func(key string) interface{} {
		return key
	}

	cache := NewCache(loader, 100)

	// Test getting a value that doesn't exist in the cache
	value := cache.Get("key1")
	if value != "key1" {
		t.Errorf("Expected value 'key1', but got '%v'", value)
	}

	// Test getting a value that exists in the cache
	value = cache.Get("key1")
	if value != "key1" {
		t.Errorf("Expected value 'key1', but got '%v'", value)
	}

	// Test getting a value that has expired in the cache
	time.Sleep(time.Second * 11)
	value = cache.Get("key1")
	if value != "key1" {
		t.Errorf("Expected value 'key1', but got '%v'", value)
	}
	t.Logf("value: %v", value)
}

func TestCache_Clear(t *testing.T) {
	loader := func(key string) interface{} {
		return key
	}

	cache := NewCache(loader, 100)
	// Add some values to the cache
	cache.Get("key1")
	cache.Get("key2")
	cache.Get("key3")

	// Clear the cache
	cache.Clear()

	// Test that the cache is empty
	if len(cache.Items()) != 0 {
		t.Errorf("Expected cache to be empty, but it has %v items", len(cache.Items()))
	}

}

func TestCache_Listeners(t *testing.T) {
	loader := func(key string) interface{} {
		return key
	}

	cache := NewCache(loader, 100)

	// Add a listener function
	var listenerValue interface{}
	listener := func(key string, value interface{}) {
		listenerValue = value
	}

	cache.AddListener(listener)
	// Get a value from the cache

	cache.Get("key1")
	// Test that the listener function was called with the correct value
	if listenerValue != "key1" {
		t.Errorf("Expected listener value 'key1', but got '%v'", listenerValue)
	}
}
