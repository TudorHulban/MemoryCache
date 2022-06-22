package cache

import (
	"context"
	"sync"
)

// InMemoryCache contains cached domain model data.
type InMemoryCache struct {
	domainModel DomainModel
	ctx         context.Context
	cache       map[int64][]byte

	mu     sync.Mutex
	chStop chan struct{}

	secondsTTL             uint
	secondsBetweenCleanUps uint
}

// NewCache is constructor for in memory cache.
func NewCache(ctx context.Context, domain DomainModel, config ...CacheOption) *InMemoryCache {
	res := InMemoryCache{
		domainModel:            domain,
		ctx:                    ctx,
		secondsBetweenCleanUps: 3,
		cache:                  make(map[int64][]byte),
	}

	for _, option := range config {
		option(&res)
	}

	if res.secondsTTL > 0 {
		res.chStop = make(chan struct{})

		go res.Clean()
	}

	return &res
}

// Close releases resources held by the memory cache.
func (c *InMemoryCache) Close() {
	if c.secondsTTL == 0 {
		return
	}

	// release resources
	c.stop()
}
