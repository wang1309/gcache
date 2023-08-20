package gcache

import "time"

type CacheItem struct {
	value  interface{}
	expire time.Time
}
