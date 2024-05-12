package cachettl

import (
	"context"
	"time"
)

// stop stops the cleaning of the cache and releases the resources in use.
func (c *InMemoryCache) stop() {
	c.chStop <- struct{}{}

	close(c.chStop)
}

func (c *InMemoryCache) clean(ctx context.Context) {
	ticker := time.NewTicker(
		time.Duration(c.secondsBetweenCleanUps * uint(time.Second)),
	)

	for {
		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
			c.deleteExpired()

		case <-c.chStop:
			ticker.Stop()
			return
		}
	}
}

func (c *InMemoryCache) deleteExpired() {
	c.mu.Lock()

	for keyDTO := range c.cache {
		if c.isTimeExpired(
			keyDTO,
			c.secondsTTL,
			time.Now(),
		) {
			delete(c.cache, keyDTO)
		}
	}

	c.mu.Unlock()
}
