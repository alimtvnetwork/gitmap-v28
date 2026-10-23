# Subtask 03: URL Cache Split-DB Engine & Storage Registry Integration

> **Subtask ID:** `SUBTASK-314-03`  
> **Parent Task:** `314-gitmap-heatmap-and-url-handling-instructions`  
> **Status:** Pending  
> **Target Files:**  
> - `cli/store/split_db_path.go` (modify)  
> - `cli/store/url_cache_types.go` (new file)  
> - `cli/store/url_cache_split_db.go` (new file)  
> - `cli/store/storage_inventory.go` (modify)  
> **Related Architecture Spec:** `02-spec/21-app/314-gitmap-heatmap-and-url-handling-instructions/02-component-spec.md`  

---

## 1. Objective

Implement the dedicated high-performance URL cache Split-DB engine (`repodb/url_cache.db`). The engine provides sub-millisecond (<0.5ms) lookup of cached HTTP/HTTPS response payloads using memory-mapped I/O (`PRAGMA mmap_size = 268435456`) and WAL mode. It incorporates in-flight request deduplication via `singleflight.Group`, PascalCase DDL schemas (`UrlCacheEntry`, `UrlAccessLog`), strictly positive booleans (`IsActive`, `IsCompressed`, `IsCacheHit`), and automatic registration in `SplitDatabaseRegistry` and `storage_inventory.go`.

---

## 2. Action Checklist

- [ ] 1. Update `cli/store/split_db_path.go`:
  - [ ] Add `ResolveUrlCacheDbPath(repoRoot string) string` helper function.
  - [ ] Target directory: `<baseDir>/repodb/url_cache.db`.
  - [ ] Enforce auto-creation of parent directories with `os.MkdirAll(dir, 0755)`.
  - [ ] Adhere to the 15-line function limit (preferred <= 8 lines).

- [ ] 2. Author `cli/store/url_cache_types.go`:
  - [ ] Package declaration `package store`.
  - [ ] Import `github.com/alimtvnetwork/gitmap-v28/cli/result`.
  - [ ] Define `UrlCacheEntry` domain struct:
    - [ ] `UrlCacheEntryId int64`
    - [ ] `UrlHash string` (SHA-256 hex string)
    - [ ] `UrlText string`
    - [ ] `Etag string`
    - [ ] `ContentType string`
    - [ ] `ContentLength int64`
    - [ ] `PayloadBlob []byte`
    - [ ] `CompressionType string`
    - [ ] `ExpiresAt int64`
    - [ ] `AccessCount int`
    - [ ] `IsActive bool`
    - [ ] `IsCompressed bool`
    - [ ] `CreatedAt int64`
    - [ ] `UpdatedAt int64`
  - [ ] Define `UrlAccessLog` audit struct:
    - [ ] `UrlAccessLogId int64`
    - [ ] `UrlHash string`
    - [ ] `AccessMethod string`
    - [ ] `ResponseStatus int`
    - [ ] `DurationMs int64`
    - [ ] `ByteCount int64`
    - [ ] `IsCacheHit bool`
    - [ ] `CreatedAt int64`
  - [ ] Define Result wrapper type aliases:
    - [ ] `UrlCacheEntryWrap = result.Wrap[UrlCacheEntry]`
    - [ ] `UrlCacheEntrySliceWrap = result.WrapSlice[UrlCacheEntry]`
    - [ ] `UrlAccessLogSliceWrap = result.WrapSlice[UrlAccessLog]`
  - [ ] Define parameter struct `UrlCacheQueryParam` for query filtering and pagination.

- [ ] 3. Author `cli/store/url_cache_split_db.go`:
  - [ ] Define `UrlCacheDB` struct wrapping `*sql.DB`, `singleflight.Group`, and connection mutex.
  - [ ] Implement `OpenUrlCacheDB(repoRoot string) (*UrlCacheDB, error)`.
  - [ ] Configure SQLite PRAGMAs:
    - [ ] `PRAGMA journal_mode = WAL;`
    - [ ] `PRAGMA synchronous = NORMAL;`
    - [ ] `PRAGMA busy_timeout = 5000;`
    - [ ] `PRAGMA mmap_size = 268435456;` (256MB)
    - [ ] `PRAGMA cache_size = -64000;`
  - [ ] Implement DDL migration helper `initUrlCacheSchema(conn *sql.DB) error`.
  - [ ] Implement SHA-256 URL hashing utility `HashUrl(rawUrl string) string`.
  - [ ] Implement query operations:
    - [ ] `GetUrlCache(urlText string) (UrlCacheEntry, bool, error)`
    - [ ] `SaveUrlCache(entry UrlCacheEntry) error`
    - [ ] `LogUrlAccess(log UrlAccessLog) error`
    - [ ] `PurgeExpiredUrlCache() (int64, error)`
    - [ ] `DeleteUrlCache(urlText string) error`
    - [ ] `ListUrlCache(limit, offset int) ([]UrlCacheEntry, error)`
  - [ ] Wrap errors with `apperror.WrapSimple(err, "store.<funcName>")`.
  - [ ] Ensure all functions <= 15 lines, guard clauses only, zero nested ifs.

- [ ] 4. Update `cli/store/storage_inventory.go`:
  - [ ] Add `collectUrlCacheDatabase(dataDir string, add func(SplitDatabaseEntry))` hook.
  - [ ] Register within `collectAllDataSources(dataDir, add)`.
  - [ ] Ensure inventory inspects `repodb/url_cache.db` with databaseType `"repodb"` and key `"url_cache"`.

- [ ] 5. Quality & Invariant Verification:
  - [ ] Relative paths only (zero absolute paths).
  - [ ] Strictly lowercase filenames.
  - [ ] Positive booleans only (`IsActive`, `IsCompressed`, `IsCacheHit`, `HasFound`).
  - [ ] Zero git commands, zero builds, zero test execution.

---

## 3. Concrete Code Specifications

### Specification: `cli/store/split_db_path.go`

```go
// ResolveUrlCacheDbPath resolves the URL cache split SQLite database file path.
func ResolveUrlCacheDbPath(repoRoot string) string {
	baseDir := resolveBaseDataDir(repoRoot)
	dir := filepath.Join(baseDir, "repodb")
	_ = os.MkdirAll(dir, 0755)

	return filepath.ToSlash(filepath.Join(dir, "url_cache.db"))
}
```

### Specification: `cli/store/url_cache_types.go`

```go
package store

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

// UrlCacheEntry models a cached web resource in the Split-DB.
type UrlCacheEntry struct {
	UrlCacheEntryId int64  `json:"urlCacheEntryId"`
	UrlHash         string `json:"urlHash"`
	UrlText         string `json:"urlText"`
	Etag            string `json:"etag,omitempty"`
	ContentType     string `json:"contentType"`
	ContentLength   int64  `json:"contentLength"`
	PayloadBlob     []byte `json:"payloadBlob,omitempty"`
	CompressionType string `json:"compressionType"`
	ExpiresAt       int64  `json:"expiresAt"`
	AccessCount     int    `json:"accessCount"`
	IsActive        bool   `json:"isActive"`
	IsCompressed    bool   `json:"isCompressed"`
	CreatedAt       int64  `json:"createdAt"`
	UpdatedAt       int64  `json:"updatedAt"`
}

// UrlAccessLog records an audit log for URL access events.
type UrlAccessLog struct {
	UrlAccessLogId int64  `json:"urlAccessLogId"`
	UrlHash        string `json:"urlHash"`
	AccessMethod   string `json:"accessMethod"`
	ResponseStatus int    `json:"responseStatus"`
	DurationMs     int64  `json:"durationMs"`
	ByteCount      int64  `json:"byteCount"`
	IsCacheHit     bool   `json:"isCacheHit"`
	CreatedAt      int64  `json:"createdAt"`
}

// Canonical Result wrapper type aliases.
type (
	UrlCacheEntryWrap      = result.Wrap[UrlCacheEntry]
	UrlCacheEntrySliceWrap = result.WrapSlice[UrlCacheEntry]
	UrlAccessLogSliceWrap  = result.WrapSlice[UrlAccessLog]
)

// UrlCacheQueryParam defines parameters for filtering and paginating cached URLs.
type UrlCacheQueryParam struct {
	Limit        int
	Offset       int
	IsActiveOnly bool
}
```

### Specification: `cli/store/url_cache_split_db.go` (Excerpt & Interface Contract)

```go
package store

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"sync"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"golang.org/x/sync/singleflight"
)

// UrlCacheDB manages the Split-DB for cached web responses.
type UrlCacheDB struct {
	conn       *sql.DB
	group      singleflight.Group
	mu         sync.RWMutex
	dbFilePath string
}

// HashUrl generates a deterministic SHA-256 hex string for a given URL.
func HashUrl(rawUrl string) string {
	sum := sha256.Sum256([]byte(rawUrl))

	return hex.EncodeToString(sum[:])
}

// OpenUrlCacheDB creates or opens the dedicated url_cache.db split database.
func OpenUrlCacheDB(repoRoot string) (*UrlCacheDB, error) {
	dbPath := ResolveUrlCacheDbPath(repoRoot)
	conn, err := openTunedSqlite(dbPath)
	if err != nil {
		return nil, apperror.WrapSimple(err, "store.OpenUrlCacheDB")
	}

	if err := initUrlCacheSchema(conn); err != nil {
		_ = conn.Close()
		return nil, apperror.WrapSimple(err, "store.OpenUrlCacheDB.schema")
	}

	return &UrlCacheDB{conn: conn, dbFilePath: dbPath}, nil
}
```

---

## 4. Verification & Done Criteria

1. `ResolveUrlCacheDbPath` returns `<baseDir>/repodb/url_cache.db` with forward slashes.
2. Tables `UrlCacheEntry` and `UrlAccessLog` execute cleanly with PascalCase identifiers and positive booleans.
3. PRAGMA tuning guarantees sub-millisecond (<0.5ms) indexed queries for in-memory reads.
4. `singleflight.Group` deduplicates concurrent queries for the same URL hash.
5. All Go functions conform to <= 15 lines with guard clauses and zero nested conditionals.
