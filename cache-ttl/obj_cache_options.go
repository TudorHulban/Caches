package cachettl

// CacheOption provides configuration options for the in memory cache.
type CacheOption func(c *InMemoryCache)

// WithTTL configuration option for time to live key values.
func WithTTL(seconds uint) CacheOption {
	return func(c *InMemoryCache) {
		c.secondsTTL = seconds
	}
}

// WithSecondsBetweenCleanUps configuration option for seconds between clean ups.
func WithSecondsBetweenCleanUps(seconds uint) CacheOption {
	return func(c *InMemoryCache) {
		c.secondsBetweenCleanUps = seconds
	}
}
