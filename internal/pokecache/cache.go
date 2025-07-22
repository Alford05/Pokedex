package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	cache    map[string]cacheEntry
	mux      sync.Mutex
	interval time.Duration
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) *Cache {
	new_cache := Cache{
		cache:    make(map[string]cacheEntry),
		mux:      sync.Mutex{},
		interval: interval,
	}
	go new_cache.reapLoop()
	return &new_cache
}

func (new_cache *Cache) Add(key string, value []byte) {
	new_cache.mux.Lock()
	defer new_cache.mux.Unlock()
	new_cache.cache[key] = cacheEntry{
		createdAt: time.Now(),
		val:       value,
	}
}

func (new_cache *Cache) Get(key string) ([]byte, bool) {
	new_cache.mux.Lock()
	defer new_cache.mux.Unlock()
	entry, found := new_cache.cache[key]
	if found {
		return entry.val, true
	} else {
		return nil, false
	}
}

func (new_cache *Cache) reapLoop() {
	ticker := time.NewTicker(new_cache.interval)
	for range ticker.C {
		new_cache.mux.Lock()
		for key, entry := range new_cache.cache {
			if time.Since(entry.createdAt) > new_cache.interval {
				delete(new_cache.cache, key)
			}
		}
		new_cache.mux.Unlock()
	}
}
