package services

import (
	"fmt"
	"sync"
	"time"
)

type Cache interface {
	Set(key string, value any) error
	Get(key string) (any, bool)
	Del(key string) bool
	StopCleanup()

	cleanup()
}

type item struct {
	value any
	ttl   int64
}

type InMemoryCache struct {
	cacheMap   map[string]item
	mu         sync.Mutex
	quit       chan struct{}
	defaultTTL time.Duration
}

func NewInMemoryCache(
	defaultTTL, cleanupInterval time.Duration,
) *InMemoryCache {
	c := &InMemoryCache{
		cacheMap:   make(map[string]item),
		quit:       make(chan struct{}),
		defaultTTL: defaultTTL,
	}

	go c.janitor(cleanupInterval)

	return c
}

func (c *InMemoryCache) Set(key string, value any) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if key == "" || value == nil {
		return fmt.Errorf("cache key/value is invalid")
	}

	var expiry int64
	if c.defaultTTL > 0 {
		expiry = time.Now().Add(c.defaultTTL).UnixNano()
	}

	c.cacheMap[key] = item{value: value, ttl: expiry}

	return nil
}

func (c *InMemoryCache) Get(key string) (any, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	v, ok := c.cacheMap[key]
	if !ok {
		return nil, false
	}

	if v.ttl > 0 && time.Now().UnixNano() > v.ttl {
		delete(c.cacheMap, key)
		return nil, false
	}

	return v.value, true
}

func (c *InMemoryCache) Del(key string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, ok := c.cacheMap[key]
	if !ok {
		return false
	}

	delete(c.cacheMap, key)

	return true
}

func (c *InMemoryCache) StopCleanup() {
	close(c.quit)
}

func (c *InMemoryCache) janitor(cleanupInterval time.Duration) {
	ticker := time.NewTicker(cleanupInterval)
	for {
		select {
		case <-ticker.C:
			c.cleanup()
		case <-c.quit:
			ticker.Stop()
			return
		}
	}
}

func (c *InMemoryCache) cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now().UnixNano()
	for k, v := range c.cacheMap {
		if v.ttl > 0 && now > v.ttl {
			delete(c.cacheMap, k)
		}
	}
}
