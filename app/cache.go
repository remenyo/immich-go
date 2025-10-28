package app

import "sync"

// UploadCache is a thread-safe cache that stores the full names of uploaded assets.
type UploadCache struct {
	mu    sync.RWMutex
	items map[string]struct{}
}

// NewUploadCache creates a new UploadCache.
func NewUploadCache() *UploadCache {
	return &UploadCache{
		items: make(map[string]struct{}),
	}
}

// Add adds an item to the cache.
func (c *UploadCache) Add(item string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[item] = struct{}{}
}

// Exists checks if an item exists in the cache.
func (c *UploadCache) Exists(item string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, found := c.items[item]
	return found
}
