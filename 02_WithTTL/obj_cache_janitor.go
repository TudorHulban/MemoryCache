package cache

import (
	"context"
	"time"
)

// stop stops the cleaning of the cache and releases the resources in use.
func (c *InMemoryCache) stop() {
	c.chStop <- struct{}{}

	close(c.chStop)
}

// Clean invokes the cache deletion periodically.
func (c *InMemoryCache) Clean(ctx context.Context) {
	ticker := time.NewTicker(time.Duration(c.secondsBetweenCleanUps * uint(time.Second)))

	for {
		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
			c.DeleteExpired()

		case <-c.chStop:
			ticker.Stop()
			return
		}
	}
}

// DeleteExpired deletes expired DTOs.
func (c *InMemoryCache) DeleteExpired() {
	now := time.Now()

	c.mu.Lock()

	for keyDTO := range c.cache {
		if c.isTimeExpired(keyDTO, c.secondsTTL, now) {
			delete(c.cache, keyDTO)
		}
	}

	c.mu.Unlock()
}
