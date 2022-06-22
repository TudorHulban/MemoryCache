package cache

import (
	"time"
)

// Stop stops the cleaning of the cache and releases the resources in use.
func (c *InMemoryCache) stop() {
	c.chStop <- struct{}{}

	close(c.chStop)
}

// Clean invokes the cache deletion periodically.
func (c *InMemoryCache) Clean() {
	ticker := time.NewTicker(time.Duration(c.secondsBetweenCleanUps * uint(time.Second)))

	for {
		select {
		case <-c.ctx.Done():
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
func (c *InMemoryCache) DeleteExpired() int {
	now := time.Now()
	var howManyDeleted int

	for keyDTO := range c.cache {
		c.mu.Lock()

		if c.isTimeExpired(keyDTO, c.secondsTTL, now) {
			delete(c.cache, keyDTO)

			howManyDeleted++
		}

		c.mu.Unlock()
	}

	return howManyDeleted
}
