package cmddb

import (
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const defaultCacheTTL = 8 * time.Second

type cachedMetadataEntry struct {
	tables    []TableInfo
	columns   map[string][]ColumnInfo
	expiresAt time.Time
}

// TableMetadataCache provides thread-safe in-memory caching for database table metadata.
type TableMetadataCache struct {
	mu      sync.RWMutex
	ttl     time.Duration
	entries map[string]cachedMetadataEntry
}

// NewTableMetadataCache creates a new TableMetadataCache with the specified TTL.
func NewTableMetadataCache(ttl time.Duration) *TableMetadataCache {
	effTTL := ttl
	if effTTL <= 0 {
		effTTL = defaultCacheTTL
	}

	return &TableMetadataCache{
		ttl:     effTTL,
		entries: make(map[string]cachedMetadataEntry),
	}
}

var globalTableMetadataCache = NewTableMetadataCache(defaultCacheTTL)

// GetGlobalTableMetadataCache returns the singleton TableMetadataCache.
func GetGlobalTableMetadataCache() *TableMetadataCache {
	return globalTableMetadataCache
}

// InvalidateTableMetadata invalidates cached table metadata for the specified database path.
func InvalidateTableMetadata(dbPath string) {
	if dbPath == "" {
		return
	}

	globalTableMetadataCache.Invalidate(dbPath)
	resolved, err := ResolveDbPath(dbPath)
	if err == nil && resolved != "" {
		globalTableMetadataCache.Invalidate(resolved)
	}
}

// InvalidateAllTableMetadata purges all entries in the metadata cache.
func InvalidateAllTableMetadata() {
	globalTableMetadataCache.InvalidateAll()
}

// Get retrieves cached tables and columns for a database path.
func (c *TableMetadataCache) Get(dbPath string) ([]TableInfo, map[string][]ColumnInfo, bool) {
	key := cleanCacheKey(dbPath)
	if key == "" {
		return nil, nil, false
	}

	entry, hasEntry := c.lookupEntry(key)
	if !hasEntry {
		return nil, nil, false
	}

	return c.validateEntry(dbPath, entry)
}

func (c *TableMetadataCache) lookupEntry(key string) (cachedMetadataEntry, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, hasEntry := c.entries[key]

	return entry, hasEntry
}

func (c *TableMetadataCache) validateEntry(dbPath string, entry cachedMetadataEntry) ([]TableInfo, map[string][]ColumnInfo, bool) {
	isExpired := time.Now().After(entry.expiresAt)
	if isExpired {
		c.Invalidate(dbPath)

		return nil, nil, false
	}

	return entry.tables, entry.columns, true
}

// Set stores tables and columns metadata in the cache for the given database path.
func (c *TableMetadataCache) Set(dbPath string, tables []TableInfo, columns map[string][]ColumnInfo) {
	key := cleanCacheKey(dbPath)
	if key == "" {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries[key] = cachedMetadataEntry{
		tables:    tables,
		columns:   columns,
		expiresAt: time.Now().Add(c.ttl),
	}
}

// Invalidate purges cached metadata for the given database path.
func (c *TableMetadataCache) Invalidate(dbPath string) {
	key := cleanCacheKey(dbPath)
	if key == "" {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.entries, key)
}

// InvalidateAll clears all cached entries across all databases.
func (c *TableMetadataCache) InvalidateAll() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries = make(map[string]cachedMetadataEntry)
}

func cleanCacheKey(dbPath string) string {
	trimmed := strings.TrimSpace(dbPath)
	if trimmed == "" {
		return ""
	}

	return strings.ToLower(filepath.Clean(trimmed))
}
