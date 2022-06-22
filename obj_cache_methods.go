package cache

import (
	"fmt"
	"time"
)

// Set stores encoded DTO items in memory.
func (c *InMemoryCache) Set(dto *DTO) error {
	c.mu.Lock()

	c.cache[dto.Key] = dto.Data

	c.mu.Unlock()
	return nil
}

// Get returns a DTO if it finds one for the passed key.
func (c *InMemoryCache) Get(key int64) (*DTO, error) {
	c.mu.Lock()

	serializedData, exists := c.cache[key]
	if !exists {
		return nil, fmt.Errorf("no cache entry found for key: `%d`", key)
	}

	c.mu.Unlock()

	return &DTO{
		Key:         key,
		DomainModel: c.domainModel,
		Data:        serializedData,
	}, nil
}

// Delete deletes a key value by passed key.
func (c *InMemoryCache) Delete(key int64) error {
	c.mu.Lock()

	delete(c.cache, key)

	c.mu.Unlock()
	return nil
}

// isKeyExpired method returns nil if cache key expired.
func (c *InMemoryCache) isTimeExpired(epochNano int64, secondsTTL uint, now time.Time) bool {
	ttl := time.Duration(time.Second * time.Duration(secondsTTL))

	return epochNano+ttl.Nanoseconds() <= now.UnixNano()
}
