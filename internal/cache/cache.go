package cache

import (
	"sync"
	"time"
)

type entry struct {
	value     any
	expiresAt time.Time
}

type Cache struct {
	mu   sync.RWMutex
	data map[string]entry
	ttl  time.Duration
	stop chan struct{}
}

func New(ttl time.Duration) *Cache {
	c := &Cache{
		data: make(map[string]entry),
		ttl:  ttl,
		stop: make(chan struct{}),
	}
	go c.cleanup(ttl * 3)
	return c
}

func (c *Cache) Set(key string, value any) {
	c.mu.Lock()
	c.data[key] = entry{value: value, expiresAt: time.Now().Add(c.ttl)}
	c.mu.Unlock()
}

func (c *Cache) Get(key string) (any, bool) {
	c.mu.RLock()
	e, ok := c.data[key]
	c.mu.RUnlock()

	if !ok || time.Now().After(e.expiresAt) {
		return nil, false
	}
	return e.value, true
}

// GetStale returns cached data even if expired (but not yet cleaned up).
// Useful as last-resort fallback when all upstream APIs fail.
func (c *Cache) GetStale(key string) (any, bool) {
	c.mu.RLock()
	e, ok := c.data[key]
	c.mu.RUnlock()

	if !ok {
		return nil, false
	}
	return e.value, true
}

func (c *Cache) Close() {
	close(c.stop)
}

func (c *Cache) cleanup(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			now := time.Now()
			c.mu.Lock()
			for k, e := range c.data {
				if now.After(e.expiresAt) {
					delete(c.data, k)
				}
			}
			c.mu.Unlock()
		case <-c.stop:
			return
		}
	}
}
