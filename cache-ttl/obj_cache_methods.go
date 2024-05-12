package cachettl

import (
	"errors"
	"fmt"
	"time"
)

func (c *InMemoryCache) Set(dto *DTO) error {
	if dto == nil {
		return errors.New("passed DTO is nil")
	}

	c.mu.Lock()

	c.cache[dto.Key] = dto.Data

	c.mu.Unlock()

	return nil
}

// Get returns a DTO if it finds one for the passed key.
// defer not used due to the performance penalty
func (c *InMemoryCache) Get(key int64) (*DTO, error) {
	c.mu.Lock()

	serializedData, exists := c.cache[key]
	if !exists {
		c.mu.Unlock()

		return nil, fmt.Errorf("no cache entry found for key: `%d`", key)
	}

	c.mu.Unlock()

	return &DTO{
			Key:         key,
			DomainModel: c.domainModel,
			Data:        serializedData,
		},
		nil
}

func (c *InMemoryCache) Delete(key int64) error {
	c.mu.Lock()

	delete(c.cache, key)

	c.mu.Unlock()

	return nil
}

// TODO: move to time duration
func (c *InMemoryCache) isTimeExpired(epochNano int64, secondsTTL uint, now time.Time) bool {
	ttl := time.Duration(time.Second * time.Duration(secondsTTL))

	return now.UnixNano() >= epochNano+ttl.Nanoseconds()
}
