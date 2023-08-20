package gcache

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func BenchmarkCache_Get(b *testing.B) {
	cache := NewCache(func(key string) interface{} {
		return key
	}, 1000)

	for i := 0; i < b.N; i++ {
		cache.Put(fmt.Sprintf("key%d", i), i)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		cache.Get(fmt.Sprintf("key%d", rand.Intn(b.N)))
	}
}

func BenchmarkCache_GetParallel(b *testing.B) {
	cache := NewCache(func(key string) interface{} {
		return key
	}, 1000)

	for i := 0; i < b.N; i++ {
		cache.Put(fmt.Sprintf("key%d", i), i)
	}

	b.ResetTimer()

	b.Run("Get", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			cache.Get(fmt.Sprintf("key%d", rand.Intn(b.N)))
		}
	})
}

func BenchmarkCache_Put(b *testing.B) {
	cache := NewCache(func(key string) interface{} {
		return key
	}, 100)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		cache.Put(fmt.Sprintf("key%d", i), i)
	}
}

func BenchmarkCache_Clear(b *testing.B) {
	cache := NewCache(func(key string) interface{} {
		return key
	}, 100)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		cache.Clear()
	}
}

func BenchmarkCache_ClearParallel(b *testing.B) {
	cache := NewCache(func(key string) interface{} {
		return key
	}, 100)

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			cache.Clear()
		}
	})
}

func BenchmarkCache_GetWithLoader(b *testing.B) {
	cache := NewCache(func(key string) interface{} {
		time.Sleep(time.Millisecond)
		return key
	}, 100)

	for i := 0; i < b.N; i++ {
		cache.Get("key")
	}
}

func BenchmarkCache_GetWithLoaderParallel(b *testing.B) {
	cache := NewCache(func(key string) interface{} {
		time.Sleep(time.Millisecond)
		return key
	}, 100)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			cache.Get("key")
		}
	})
}

func BenchmarkCache_GetWithListeners(b *testing.B) {
	cache := NewCache(func(key string) interface{} {
		return key
	}, 100)

	cache.Listeners = append(cache.Listeners, func(key string, value interface{}) {
		// do nothing
	})

	for i := 0; i < b.N; i++ {
		cache.Get("key")
	}
}

func BenchmarkCache_GetWithListenersParallel(b *testing.B) {
	cache := NewCache(func(key string) interface{} {
		return key
	}, 100)

	cache.Listeners = append(cache.Listeners, func(key string, value interface{}) {
		// do nothing
	})

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			cache.Get("key")
		}
	})
}
