package gcache

import (
	"sync"
	"time"

	"golang.org/x/sync/singleflight"
)

// Todo 解决缓存雪崩问题 -- 异步加载缓存
// Todo 多分片 map 减少锁冲突
// 减少内存分配
// 减少GC

const DefaultMaxCheckKeyNumber = 20

var single singleflight.Group

func NewCache(loader func(key string) interface{}, maxItems int) *Cache {
	cache := &Cache{
		items:    make(map[string]CacheItem),
		loader:   loader,
		maxItems: maxItems,
		queue:    make([]string, 0, maxItems),
	}

	go cache.startCleanup()
	return cache
}

type Cache struct {
	mutex     sync.RWMutex
	Stats     CacheStats
	items     map[string]CacheItem
	loader    func(key string) interface{}
	Listeners []func(key string, value interface{})
	maxItems  int      // 缓存最大容量
	queue     []string // 缓存项队列，记录缓存项的访问顺序
}

func (c *Cache) Items() map[string]CacheItem {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.items
}

func (c *Cache) Get(key string) interface{} {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if item, ok := c.items[key]; ok {
		if item.expire.Before(time.Now()) {
			delete(c.items, key)
			c.removeFromQueue(key)
		} else {
			c.Stats.Hits++
			c.moveToFront(key)
			for _, listener := range c.Listeners {
				listener(key, item.value)
			}
			return item.value
		}
	}

	value, _, _ := single.Do(key, func() (interface{}, error) {
		return c.loader(key), nil
	})

	c.items[key] = CacheItem{value, time.Now().Add(time.Second * 10)}
	c.addToQueue(key)
	c.Stats.Misses++
	for _, listener := range c.Listeners {
		listener(key, value)
	}

	if len(c.items) > c.maxItems {
		c.evict()
	}

	return value
}

func (c *Cache) Put(key string, value interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if item, ok := c.items[key]; ok {
		item.value = value
		item.expire = time.Now().Add(time.Second * 10)
		c.items[key] = item
		c.moveToFront(key)
		for _, listener := range c.Listeners {
			listener(key, value)
		}
		return
	}

	c.items[key] = CacheItem{value, time.Now().Add(time.Second * 10)}
	c.addToQueue(key)
	c.Stats.Misses++
	for _, listener := range c.Listeners {
		listener(key, value)
	}
	if len(c.items) >= c.maxItems {
		c.evict()
	}
}

func (c *Cache) startCleanup() {
	ticker := time.NewTicker(time.Minute)
	for {
		select {
		case <-ticker.C:
			c.cleanup()
		}
	}
}

// cleanup clean expire key
func (c *Cache) cleanup() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

Loop:
	now := time.Now()
	cursor := 0
	evictNumber := 0
	for key, item := range c.items {
		if item.expire.Before(now) {
			delete(c.items, key)
			c.removeFromQueue(key)
			evictNumber++
		}
		cursor++

		duration := time.Now().Sub(now)
		if cursor > DefaultMaxCheckKeyNumber {
			if float32(evictNumber/DefaultMaxCheckKeyNumber) >= 0.2 && duration.Seconds() < 1 {
				goto Loop
			}

			break
		}
	}
}

func (c *Cache) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.items = make(map[string]CacheItem)
	c.queue = make([]string, 0, c.maxItems)
}

func (c *Cache) AddListener(fun func(key string, value interface{})) {
	c.Listeners = append(c.Listeners, fun)
}

func (c *Cache) addToQueue(key string) {
	c.queue = append(c.queue, key)
}

func (c *Cache) removeFirstFromQueue() string {
	if len(c.queue) == 0 {
		return ""
	}
	key := c.queue[0]
	c.queue = c.queue[1:]
	return key
}

func (c *Cache) removeLastFromQueue() string {
	if len(c.queue) == 0 {
		return ""
	}
	key := c.queue[len(c.queue)-1]
	c.queue = c.queue[:len(c.queue)-1]
	return key
}

func (c *Cache) removeItemFromQueue(key string) {
	for i, k := range c.queue {
		if k == key {
			c.queue = append(c.queue[:i], c.queue[i+1:]...)
			break
		}
	}
}

func (c *Cache) removeFirstIfFull() {
	if len(c.items) >= c.maxItems {
		key := c.removeFirstFromQueue()
		delete(c.items, key)
	}
}

func (c *Cache) evict() {
	key := c.removeLastFromQueue()
	delete(c.items, key)
}

func (c *Cache) moveToFront(key string) {
	c.removeItemFromQueue(key)
	c.addToQueue(key)
}

func (c *Cache) removeFromQueue(key string) {
	c.removeItemFromQueue(key)
	c.removeFirstIfFull()
}
