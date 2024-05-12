package cachevanilla

import (
	"fmt"
	"sync"
)

type MemoryCache struct {
	cache map[string][]byte

	mu sync.RWMutex
}

func NewCache() *MemoryCache {
	return &MemoryCache{
		cache: make(map[string][]byte),
	}
}

func (c *MemoryCache) Set(key string, value []byte) error {
	c.mu.Lock()

	c.cache[key] = value

	c.mu.Unlock()

	return nil
}

func (c *MemoryCache) Get(key string) ([]byte, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	serializedValue, exists := c.cache[key]
	if !exists {
		return nil,
			fmt.Errorf(
				"no cache entry found for key: `%s`",
				key,
			)
	}

	return serializedValue, nil
}

func (c *MemoryCache) Delete(key string) error {
	c.mu.Lock()

	delete(c.cache, key)

	c.mu.Unlock()

	return nil
}
