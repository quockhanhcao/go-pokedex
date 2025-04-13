package pokecache

import (
	"sync"
	"time"
)

type CacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	cacheEntries map[string]CacheEntry
	mutex        sync.Mutex
}

func (c *Cache) Add(key string, val []byte) {
    c.mutex.Lock()
    c.cacheEntries[key] = CacheEntry{
        createdAt: time.Now(),
        val:       val,
    }
    c.mutex.Unlock()
}

func (c *Cache) Get(key string) ([]byte, bool) {
    c.mutex.Lock()
    defer c.mutex.Unlock()
	data, ok := c.cacheEntries[key]
	if ok {
		return data.val, true
	}
	return nil, false
}

func (c *Cache) reapLoop(interval time.Duration) {
    if interval == 0 {
        interval = 5 * time.Second
    }
    ticker := time.NewTicker(interval)
    done := make(chan bool)
    go func() {
        for {
            select {
            case <-done:
                return
            case <-ticker.C:
                c.mutex.Lock()
                for key, entry := range c.cacheEntries {
                    if time.Since(entry.createdAt) > interval {
                        delete(c.cacheEntries, key)
                    }
                }
                c.mutex.Unlock()
            }
        }
    }()
}

func NewCache(interval time.Duration) *Cache {
	cache := &Cache{
		cacheEntries: make(map[string]CacheEntry),
		mutex:        sync.Mutex{},
	}
	cache.reapLoop(interval)
	return cache
}
