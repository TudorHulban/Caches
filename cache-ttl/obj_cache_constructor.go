package cachettl

import (
	"context"
	"sync"
)

// InMemoryCache contains cached domain model data.
type InMemoryCache struct {
	domainModel DomainModel
	cache       map[int64][]byte

	mu     sync.Mutex
	chStop chan struct{}

	secondsTTL             uint // TODO: move to time duration
	secondsBetweenCleanUps uint // TODO: move to time duration
}

// NewCache is constructor for in memory cache.
func NewCache(ctx context.Context, domain DomainModel, config ...CacheOption) *InMemoryCache {
	result := InMemoryCache{
		domainModel: domain,
		cache:       make(map[int64][]byte, 100),
	}

	for _, option := range config {
		option(&result)
	}

	if result.secondsTTL > 0 {
		result.chStop = make(chan struct{})

		go result.clean(ctx)
	}

	return &result
}

func (c *InMemoryCache) Close() {
	if c.secondsTTL == 0 {
		return
	}

	c.stop()
}
