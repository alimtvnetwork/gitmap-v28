// Package fsutil — discovery_cache.go provides in-memory discovery caching.
package fsutil

import "sync"

type DiscoveryCache struct {
	mu    sync.RWMutex
	cache map[string][]string
}

func NewDiscoveryCache() *DiscoveryCache {
	return &DiscoveryCache{
		cache: make(map[string][]string),
	}
}

func (dc *DiscoveryCache) Get(dir string) ([]string, bool) {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	repos, ok := dc.cache[dir]
	return repos, ok
}

func (dc *DiscoveryCache) Set(dir string, repos []string) {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	dc.cache[dir] = repos
}
