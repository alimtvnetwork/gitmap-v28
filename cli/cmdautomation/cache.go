package cmdautomation

import (
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
)

// CacheStats captures cache capacity and lookup metrics.
type CacheStats struct {
	TotalFiles int   `json:"totalFiles"`
	TotalBytes int64 `json:"totalBytes"`
	Hits       int64 `json:"hits"`
	Misses     int64 `json:"misses"`
}

// MemoryCache provides sub-millisecond in-memory caching with zero disk bloat.
type MemoryCache struct {
	mu     sync.RWMutex
	files  map[string][]byte
	hits   int64
	misses int64
}

var globalCache = &MemoryCache{
	files: make(map[string][]byte),
}

// GlobalCache returns the shared process-level memory cache.
func GlobalCache() *MemoryCache {
	return globalCache
}

func (c *MemoryCache) GetFile(path string) ([]byte, bool) {
	norm := filepath.Clean(path)
	c.mu.RLock()
	data, hasFile := c.files[norm]
	c.mu.RUnlock()

	if hasFile {
		atomic.AddInt64(&c.hits, 1)
		return data, true
	}

	atomic.AddInt64(&c.misses, 1)
	return nil, false
}

func (c *MemoryCache) SetFile(path string, data []byte) {
	norm := filepath.Clean(path)
	c.mu.Lock()
	defer c.mu.Unlock()
	c.files[norm] = data
}

func (c *MemoryCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.files = make(map[string][]byte)
}

func (c *MemoryCache) Stats() CacheStats {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var totalBytes int64
	for _, b := range c.files {
		totalBytes += int64(len(b))
	}

	return CacheStats{
		TotalFiles: len(c.files),
		TotalBytes: totalBytes,
		Hits:       atomic.LoadInt64(&c.hits),
		Misses:     atomic.LoadInt64(&c.misses),
	}
}

func (c *MemoryCache) Warm(dir string) int {
	files := collectSearchFiles(dir, nil)
	count := 0
	for _, f := range files {
		data, err := os.ReadFile(f)
		if err == nil && !HasBinaryContent(data) {
			c.SetFile(f, data)
			count++
		}
	}
	return count
}
