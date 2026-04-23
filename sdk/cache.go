package sdk

import (
	"sync"
	"time"
)

type cacheEntry struct {
	enabled   bool
	expiresAt time.Time
}

type localCache struct {
	data    sync.Map
	ttl     time.Duration
	maxSize int
	size    int64
}

func newLocalCache(maxSize int, ttl time.Duration) *localCache {
	return &localCache{
		maxSize: maxSize,
		ttl:     ttl,
	}
}

func cacheKey(flagKey, projectID, environmentID string) string {
	return flagKey + ":" + projectID + ":" + environmentID
}

func (c *localCache) get(flagKey, projectID, environmentID string) (bool, bool) {
	key := cacheKey(flagKey, projectID, environmentID)
	val, ok := c.data.Load(key)
	if !ok {
		return false, false
	}

	entry := val.(cacheEntry)
	if time.Now().After(entry.expiresAt) {
		c.data.Delete(key)
		return false, false
	}
	return entry.enabled, true
}

func (c *localCache) set(flagKey, projectID, environmentID string, enabled bool) {
	key := cacheKey(flagKey, projectID, environmentID)
	c.data.Store(key, cacheEntry{
		enabled:   enabled,
		expiresAt: time.Now().Add(c.ttl),
	})
}

func (c *localCache) invalidate(flagKey, projectID, environmentID string) {
	key := cacheKey(flagKey, projectID, environmentID)
	c.data.Delete(key)
}

func (c *localCache) invalidateAll() {
	c.data.Range(func(key, _ any) bool {
		c.data.Delete(key)
		return true
	})
}
