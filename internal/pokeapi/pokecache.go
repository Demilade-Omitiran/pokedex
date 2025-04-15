package pokeapi

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	mu       sync.Mutex
	Entry    map[string]cacheEntry
	interval time.Duration
}

func (c *Cache) Add(key string, val []byte) {
	entryVal := cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	c.Entry[key] = entryVal
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if entry, ok := c.Entry[key]; ok {
		return entry.val, true
	}

	return nil, false
}

func (c *Cache) reapLoop() {
	go func() {
		for {
			time.Sleep(c.interval)

			c.mu.Lock()
			for key, entry := range c.Entry {
				if time.Since(entry.createdAt) > c.interval {
					delete(c.Entry, key)
				}
			}
			c.mu.Unlock()
		}
	}()
}

func NewCache(interval time.Duration) *Cache {
	newCache := &Cache{
		interval: interval,
		Entry:    make(map[string]cacheEntry),
	}

	newCache.reapLoop()

	return newCache
}
