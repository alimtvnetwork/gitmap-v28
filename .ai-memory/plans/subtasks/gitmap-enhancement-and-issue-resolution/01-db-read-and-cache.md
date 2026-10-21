# Subtask 01: Database Read Introspection & In-Memory TTL Cache

## Status: READY
- **Subtask Slug**: `01-db-read-and-cache`
- **Parent Plan**: `.ai-memory/plans/gitmap-enhancement-and-issue-resolution.md`
- **Architecture Spec**: `02-spec/21-app/gitmap-enhancement-and-issue-resolution/01-architecture-spec.md`
- **Owner**: Worker 01 (Database Operations & Cache Engineer)
- **Target Subsystem**: `cli/cmddb`
- **Exclusive Owned Files**:
  - `cli/cmddb/cmd_read.go`
  - `cli/cmddb/cmd_read_helpers.go`
  - `cli/cmddb/cmddb_cache.go` (new file)
  - `cli/cmddb/cmddb_cache_test.go` (new file)

---

## 1. Context & Objectives

Currently, invoking `gitmap db read <path>` with only a database file path or alias defaults to `handleReadLs`, which only retrieves table and view names from SQLite. Developers and automated agents require comprehensive table introspection:
1. **Bare Path / Alias Introspection**:
   - Inspect all table names, types (`table` or `view`), column counts, row counts (`SELECT COUNT(*)`), and schema DDL (`sqlite_master`).
   - Render clean tabular output in the terminal and a structured envelope when `--json` is supplied.
2. **In-Memory 5–10s Metadata Cache (`cmddb_cache.go`)**:
   - Default TTL of 8 seconds (range 5–10s).
   - Sub-millisecond cache validation (<0.04ms execution target) using a dual-check strategy: `os.Stat` file modification time/size check and SQLite `PRAGMA data_version;` verification.
   - Thread-safe storage with `sync.RWMutex`.
   - Public eviction API (`InvalidateDbTableCache`, `InvalidateAllDbTableCaches`).
3. **Coding Discipline & FastGate Standards**:
   - Every function <= 15 lines (<= 8 lines preferred).
   - Zero nested `if` statements (guard clauses only).
   - Strict positive boolean naming (`is*`, `has*`).
   - All errors wrapped via `appfault` / `apperror`.
   - 100% relative paths only.

---

## 2. Step-by-Step Implementation Work Breakdown

### Step 1: Architect In-Memory Metadata Cache in `cli/cmddb/cmddb_cache.go`

1. **Define Core Data Models**:
   ```go
   package cmddb

   import (
       "database/sql"
       "os"
       "sync"
       "time"
   )

   // TableMetadata encapsulates table schema, column counts, and row counts.
   type TableMetadata struct {
       Name        string `json:"name"`
       Type        string `json:"type"`
       ColumnCount int    `json:"columnCount"`
       RowCount    int64  `json:"rowCount"`
       Schema      string `json:"schema"`
   }

   // TableCacheEntry holds cached table metadata alongside cache validation tokens.
   type TableCacheEntry struct {
       DbPath      string
       DataVersion int64
       ModTime     time.Time
       FileSize    int64
       CachedAt    time.Time
       ExpiresAt   time.Time
       Tables      []TableMetadata
   }

   // TableMetadataCache provides thread-safe in-memory caching for SQLite table metadata.
   type TableMetadataCache struct {
       mu      sync.RWMutex
       entries map[string]*TableCacheEntry
       ttl     time.Duration
   }
   ```

2. **Initialize Global Cache Singleton**:
   ```go
   var (
       globalTableCache     *TableMetadataCache
       onceGlobalTableCache sync.Once
   )

   // GetGlobalTableCache returns the singleton cache instance with 8-second default TTL.
   func GetGlobalTableCache() *TableMetadataCache {
       onceGlobalTableCache.Do(func() {
           globalTableCache = &TableMetadataCache{
               entries: make(map[string]*TableCacheEntry),
               ttl:     8 * time.Second,
           }
       })
       return globalTableCache
   }
   ```

3. **Implement Freshness Validation (<0.04ms)**:
   - Function: `ValidateCacheFreshness(entry *TableCacheEntry, conn *sql.DB, dbPath string) bool`
   - Guard 1: Verify `time.Now().Before(entry.ExpiresAt)`. Return false if expired.
   - Guard 2: Call `os.Stat(dbPath)`. Return false if stat fails, `info.ModTime() != entry.ModTime`, or `info.Size() != entry.FileSize`.
   - Guard 3: Query `PRAGMA data_version;`. Return false if query fails or `currentDataVersion != entry.DataVersion`.
   - Return true if all checks pass.

4. **Implement Cache Lookup & Store**:
   - `GetCachedTableMetadata(conn *sql.DB, dbPath string) ([]TableMetadata, bool)`:
     - Read-lock `RLock()`.
     - Check map entry. If present and `ValidateCacheFreshness` returns true, return cached copy and `true`.
     - Otherwise return `nil, false`.
   - `SetCachedTableMetadata(conn *sql.DB, dbPath string, tables []TableMetadata)`:
     - Stat file and query `PRAGMA data_version;`.
     - Write-lock `Lock()`.
     - Store `TableCacheEntry` with `ExpiresAt = time.Now().Add(c.ttl)`.

5. **Implement Eviction Functions**:
   - `InvalidateDbTableCache(dbPath string)`:
     - Normalize path via `filepath.Clean`.
     - Acquire write lock and delete key from map.
   - `InvalidateAllDbTableCaches()`:
     - Acquire write lock and reset entries map.

---

### Step 2: Implement Introspection Helpers in `cli/cmddb/cmd_read_helpers.go`

1. **Implement `introspectDatabaseTables(conn *sql.DB) ([]TableMetadata, error)`**:
   - Query `sqlite_master`:
     `SELECT name, type, sql FROM sqlite_master WHERE type IN ('table', 'view') AND name NOT LIKE 'sqlite_%' ORDER BY name ASC;`
   - For each table/view:
     - Query column count via `PRAGMA table_info("<name>");`.
     - For tables (not views), query row count via `SELECT COUNT(*) FROM "<name>";`.
     - For views, set row count to `-1` or query count safely.
     - Build `TableMetadata` entry.

2. **Implement Output Renderers**:
   - `renderTableIntrospectionTerminal(w io.Writer, dbPath string, tables []TableMetadata, elapsed time.Duration, isFromCache bool)`:
     - Render database path and table summary header.
     - Output ASCII/Unicode table containing: `Table Name`, `Type`, `Columns`, `Rows`, `DDL Preview`.
     - Append latency footer: `Introspection completed in X.XXms (cached: <isFromCache>)`.
   - `renderTableIntrospectionJSON(w io.Writer, dbPath, alias string, tables []TableMetadata, elapsed time.Duration, isFromCache bool)`:
     - Format JSON envelope matching spec:
       ```json
       {
         "operation": "db-read:inspect",
         "status": "success",
         "database": dbPath,
         "alias": alias,
         "data": {
           "totalTables": len(tables),
           "isFromCache": isFromCache,
           "tables": tables
         },
         "durationMs": elapsed.Seconds() * 1000
       }
       ```

---

### Step 3: Wire Bare Path Dispatch in `cli/cmddb/cmd_read.go`

1. **Update `dispatchReadOperation(flags readFlags, posArgs []string) *appfault.AppError`**:
   - When `len(posArgs) == 1` and `posArgs[0]` does not match SQL tokens (`SELECT`, `INSERT`, etc.) and is not a subcommand (`ls`, `table`):
     - Replace call to `handleReadLs(flags, posArgs[0])` with `handleReadInspect(flags, posArgs[0])`.
   - When `len(posArgs) == 0` and `flags.aliasName != ""`:
     - Route to `handleReadInspect(flags, flags.aliasName)`.
   - When `len(posArgs) == 0` and default alias is active:
     - Route to `handleReadInspect(flags, "master")`.

2. **Implement `handleReadInspect(flags readFlags, dbArg string) *appfault.AppError`**:
   - Open read-only connection via `openReadOnlyResolvedDb(dbArg)`.
   - Check cache: `tables, hasCache := GetGlobalTableCache().GetCachedTableMetadata(conn, dbPath)`.
   - If not cached:
     - Fetch `tables, err = introspectDatabaseTables(conn)`.
     - Store `GetGlobalTableCache().SetCachedTableMetadata(conn, dbPath, tables)`.
   - Render output (terminal or JSON) based on `flags.isJson`.
   - Close connection cleanly and return `nil`.

---

### Step 4: Unit Testing & Performance Verification in `cli/cmddb/cmddb_cache_test.go`

1. **`TestTableCacheHitAndMiss`**:
   - Create temporary SQLite database with 3 tables.
   - Run first introspection -> verify cache miss (`hasCache == false`).
   - Run second introspection -> verify cache hit (`hasCache == true`).
2. **`TestTableCacheTtlExpiration`**:
   - Set cache TTL to 50ms.
   - Populate cache, sleep 60ms, query -> verify entry expired.
3. **`TestTableCacheExternalModification`**:
   - Populate cache.
   - Modify SQLite database using external write (`INSERT INTO ...`).
   - Query cache -> verify `PRAGMA data_version` detects modification and invalidates entry.
4. **`TestTableCacheExplicitInvalidation`**:
   - Populate cache.
   - Call `InvalidateDbTableCache(dbPath)`.
   - Query cache -> verify entry evicted.
5. **`BenchmarkTableCacheFreshnessValidation`**:
   - Benchmark `ValidateCacheFreshness`.
   - Assert mean execution time < 0.04ms (40 microseconds).

---

## 3. Verification & Quality Gates

- **Function Length Limit**: All functions <= 15 lines. Helper functions extracted for table scanning, row counting, and pragma checks.
- **Positive Booleans**: Use `isFromCache`, `hasCache`, `isValid`, `isTable`, `isJson`. No negative flags.
- **Guard Clauses**: Early returns on invalid handles, empty results, and failed pragma queries.
- **AppError Wrapping**: Return `*appfault.AppError` with clear error codes (`ERR_DB_NOT_FOUND`, `ERR_INTROSPECTION_FAILED`).
- **Relative Paths Only**: Strictly relative Git paths in comments and code.
