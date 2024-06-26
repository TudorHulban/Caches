package cachelruany

// CfgCache consolidates the cache capacity and key based on which it can be retrieved.
type CfgCache struct {
	Name     string
	Capacity uint16
}

// NewCachesForMethods is a constructor for the map holding the caches.
func NewCaches(cacheConfigurations ...CfgCache) map[string]*CacheLRU {
	res := make(map[string]*CacheLRU, len(cacheConfigurations))

	for _, cacheConfiguration := range cacheConfigurations {
		res[cacheConfiguration.Name] = NewCacheLRU(cacheConfiguration.Capacity)
	}

	return res
}
