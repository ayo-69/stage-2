package utils

import (
	"sync"
	"time"
)

type CacheItem struct {
	Data      interface{}
	ExpiresAt time.Time
}

var cache = sync.Map{}

func Set(key string, value interface{}, ttl time.Duration) {
	cache.Store(key, CacheItem{
		Data:      value,
		ExpiresAt: time.Now().Add(ttl),
	})
}

func Get(key string) (interface{}, bool) {
	item, ok := cache.Load(key)
	if !ok {
		return nil, false
	}
	c := item.(CacheItem)
	if time.Now().After(c.ExpiresAt) {
		cache.Delete(key)
		return nil, false
	}
	return nil, false
}
