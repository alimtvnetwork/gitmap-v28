package dbengine

import (
	"sync"
)

// CompiledQueryCache provides a thread-safe in-memory cache for compiled SQL queries.
type CompiledQueryCache struct {
	mu    sync.RWMutex
	cache map[string]string
}

// NewCompiledQueryCache creates an initialized CompiledQueryCache.
func NewCompiledQueryCache() *CompiledQueryCache {
	return &CompiledQueryCache{
		cache: make(map[string]string),
	}
}

// Get retrieves a compiled SQL query by cache key.
func (c *CompiledQueryCache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	sqlStr, found := c.cache[key]
	return sqlStr, found
}

// Put stores a compiled SQL query in the cache.
func (c *CompiledQueryCache) Put(key string, sqlStr string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[key] = sqlStr
}

// Clear flushes all cached queries.
func (c *CompiledQueryCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache = make(map[string]string)
}

// Size returns the count of cached SQL queries.
func (c *CompiledQueryCache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.cache)
}

// GlobalQueryCache is the process-level default query compilation cache.
var GlobalQueryCache = NewCompiledQueryCache()
