package cmddb

import (
	"database/sql"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// DefaultTableCacheTTL is the expiration TTL duration for in-memory table metadata.
const DefaultTableCacheTTL = 10 * time.Second

// TableMetadata is an alias for TableInfo for specification compatibility.
type TableMetadata = TableInfo

// CachedTableMetadata stores cached table introspection data with validation metadata.
type CachedTableMetadata struct {
	Tables      []TableInfo
	LoadedAt    time.Time
	DataVersion int64
	FileModTime time.Time
	FileSize    int64
}

var (
	tableCacheMu sync.RWMutex
	tableCache   = make(map[string]*CachedTableMetadata)
)

func normalizeDbCacheKey(dbPath string) string {
	cleaned := filepath.Clean(dbPath)

	return filepath.ToSlash(cleaned)
}

// InvalidateDbTableCache removes the cached table metadata for the given database path.
func InvalidateDbTableCache(dbPath string) {
	if dbPath == "" {
		return
	}

	key := normalizeDbCacheKey(dbPath)
	tableCacheMu.Lock()
	defer tableCacheMu.Unlock()

	delete(tableCache, key)
}

// ClearAllTableCaches clears all entries from the in-memory table cache.
func ClearAllTableCaches() {
	tableCacheMu.Lock()
	defer tableCacheMu.Unlock()

	tableCache = make(map[string]*CachedTableMetadata)
}

// InvalidateAllDbTableCaches clears all entries from the in-memory table cache.
func InvalidateAllDbTableCaches() {
	ClearAllTableCaches()
}

func lookupTableCache(key string) (*CachedTableMetadata, bool) {
	tableCacheMu.RLock()
	defer tableCacheMu.RUnlock()

	entry, hasEntry := tableCache[key]

	return entry, hasEntry
}

func isCacheTtlValid(loadedAt time.Time, ttl time.Duration) bool {
	return time.Since(loadedAt) <= ttl
}

func hasMatchingFileStat(dbPath string, modTime time.Time, size int64) bool {
	info, err := os.Stat(dbPath)
	if err != nil {
		return false
	}

	return info.ModTime().Equal(modTime) && info.Size() == size
}

func fetchPragmaDataVersion(conn *sql.DB) (int64, error) {
	var dataVersion int64
	err := conn.QueryRow("PRAGMA data_version;").Scan(&dataVersion)
	if err != nil {
		return 0, err
	}

	return dataVersion, nil
}

func hasMatchingDataVersion(conn *sql.DB, expected int64) bool {
	currentVersion, err := fetchPragmaDataVersion(conn)
	if err != nil {
		return false
	}

	return currentVersion == expected
}

// ValidateCacheFreshness checks TTL, file stat, and sqlite data_version.
func ValidateCacheFreshness(entry *CachedTableMetadata, conn *sql.DB, dbPath string) bool {
	if entry == nil {
		return false
	}
	if !isCacheTtlValid(entry.LoadedAt, DefaultTableCacheTTL) {
		return false
	}
	if !hasMatchingFileStat(dbPath, entry.FileModTime, entry.FileSize) {
		return false
	}

	return hasMatchingDataVersion(conn, entry.DataVersion)
}

func copyTableInfos(tables []TableInfo) []TableInfo {
	copied := make([]TableInfo, len(tables))
	copy(copied, tables)

	return copied
}

func buildCacheEntry(conn *sql.DB, dbPath string, tables []TableInfo) (*CachedTableMetadata, bool) {
	info, err := os.Stat(dbPath)
	if err != nil {
		return nil, false
	}

	version, err := fetchPragmaDataVersion(conn)
	if err != nil {
		return nil, false
	}

	return &CachedTableMetadata{
		Tables:      copyTableInfos(tables),
		LoadedAt:    time.Now(),
		DataVersion: version,
		FileModTime: info.ModTime(),
		FileSize:    info.Size(),
	}, true
}

func storeTableCacheEntry(conn *sql.DB, dbPath string, tables []TableInfo) {
	entry, hasEntry := buildCacheEntry(conn, dbPath, tables)
	if !hasEntry {
		return
	}

	key := normalizeDbCacheKey(dbPath)
	tableCacheMu.Lock()
	defer tableCacheMu.Unlock()

	tableCache[key] = entry
}

func refreshAndCacheTables(conn *sql.DB, dbPath string) ([]TableInfo, error) {
	tables, err := listTablesAndViews(conn)
	if err != nil {
		return nil, err
	}

	storeTableCacheEntry(conn, dbPath, tables)

	return tables, nil
}

func resolveCachedOrFreshTables(conn *sql.DB, dbPath string) ([]TableInfo, error) {
	key := normalizeDbCacheKey(dbPath)
	entry, hasEntry := lookupTableCache(key)
	if hasEntry && ValidateCacheFreshness(entry, conn, dbPath) {
		return copyTableInfos(entry.Tables), nil
	}

	return refreshAndCacheTables(conn, dbPath)
}

// GetCachedTablesAndViews retrieves table metadata from cache or reloads if stale.
func GetCachedTablesAndViews(conn *sql.DB, dbPath string) ([]TableInfo, error) {
	if conn == nil {
		return nil, apperror.NewExecutionError("nil database connection")
	}
	if dbPath == "" {
		return listTablesAndViews(conn)
	}

	return resolveCachedOrFreshTables(conn, dbPath)
}
