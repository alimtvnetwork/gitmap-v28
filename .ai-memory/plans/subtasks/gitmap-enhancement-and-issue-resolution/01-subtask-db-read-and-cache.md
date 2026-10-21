# Subtask 01: Database Read Introspection & 5-10s In-Memory Table Cache

## Status: COMPLETED
- **Parent Plan**: `.ai-memory/plans/gitmap-enhancement-and-issue-resolution.md`
- **Specification Reference**: `02-spec/21-app/gitmap-enhancement-and-issue-resolution/01-architecture-spec.md`
- **Owner**: Worker 01 (Database Operations & Cache Engineer)
- **Target Subsystem**: `cli/cmddb`
- **Owned Files**:
  - `cli/cmddb/cmddb_table_cache.go` (new file)
  - `cli/cmddb/cmddb_table_cache_test.go` (new file)
  - `cli/cmddb/cmd_read.go`
  - `cli/cmddb/cmd_read_helpers.go`
  - `cli/cmddb/cmd_add.go`
  - `cli/cmddb/cmd_edit.go`
  - `cli/cmddb/cmd_record_del.go`
  - `cli/cmddb/cmd_rm.go`
  - `cli/cmddb/rest_handlers_crud.go`
  - `cli/cmddb/rest_handlers_ddl.go`

---

## 1. Context & Objectives

1. **Enhanced `gitmap db read <path>`**:
   - When a user or agent invokes `gitmap db read <path>` with only a database path or alias (no query or subcommand), inspect and display all table names, types (`table`/`view`), schema DDL, column counts, and row counts via `sqlite_master`.
   - Provide both formatted terminal table output and structured JSON envelope output (`--json`).
2. **In-Memory TTL Table Cache (`cmddb_table_cache.go`)**:
   - Implement an in-memory cache with a 5-10 second TTL (default 8 seconds).
   - Fast validation in <0.04ms using a dual check: file `mtime`/size via `os.Stat` and SQLite `PRAGMA data_version;`.
   - Thread-safe access via `sync.RWMutex`.
3. **Proactive Mutation Invalidation**:
   - Hook into CLI commands: `RunDbAdd`, `RunDbEdit`, `RunDbRecordDel`, `RunDbRm`.
   - Hook into REST mutation handlers: record insertions, updates, deletions, and schema DDL mutations.
   - Immediate cache eviction for the modified database path to guarantee read-your-own-writes consistency.
4. **Code Quality & FastGate Discipline**:
   - All Go functions must be <= 15 lines (<= 8 lines preferred).
   - Zero nested `if` statements (guard clauses and early returns required).
   - Positive boolean naming only (`is*` and `has*`).
   - AppError wrapping on all error return paths.

---

## 2. Step-by-Step Implementation Plan

### Step 2.1: Implement `cli/cmddb/cmddb_table_cache.go`
1. Define cache data models:
   ```go
   type TableMetadata struct {
       Name        string `json:"name"`
       Type        string `json:"type"`
       ColumnCount int    `json:"columnCount"`
       RowCount    int64  `json:"rowCount"`
       Schema      string `json:"schema"`
   }

   type TableCacheEntry struct {
       DbPath      string
       DataVersion int64
       ModTime     time.Time
       FileSize    int64
       CachedAt    time.Time
       ExpiresAt   time.Time
       Tables      []TableMetadata
   }
   ```
2. Implement thread-safe global cache storage with `sync.RWMutex`.
3. Implement `ValidateCacheFreshness(entry *TableCacheEntry, conn *sql.DB, dbPath string) bool`:
   - Guard 1: Verify `time.Now().Before(entry.ExpiresAt)`. Return false if expired.
   - Guard 2: Call `os.Stat(dbPath)`. Return false if `statErr != nil`, or `info.ModTime() != entry.ModTime`, or `info.Size() != entry.FileSize`.
   - Guard 3: Query `PRAGMA data_version;`. Return false if query fails or `currentDataVersion != entry.DataVersion`.
   - Return true if all three checks pass (<0.04ms execution).
4. Implement `GetCachedTableMetadata(conn *sql.DB, dbPath string) ([]TableMetadata, bool)`:
   - Acquire read lock `RLock()`.
   - Check if entry exists and is fresh. Return cached slice if valid.
5. Implement `SetCachedTableMetadata(conn *sql.DB, dbPath string, tables []TableMetadata)`:
   - Read current `PRAGMA data_version;` and `os.Stat(dbPath)`.
   - Store entry with `ExpiresAt = time.Now().Add(8 * time.Second)`.
6. Implement eviction API:
   - `InvalidateDbTableCache(dbPath string)`: Evicts specific normalized database path.
   - `InvalidateAllDbTableCaches()`: Clears entire cache map.

### Step 2.2: Implement Mutation Invalidation Hooks
1. In `cli/cmddb/cmd_add.go`:
   - In `RunDbAdd`: after successful `executeInsert`, call `InvalidateDbTableCache(dbPath)`.
2. In `cli/cmddb/cmd_edit.go`:
   - In `RunDbEdit`: after successful `executeUpdate`, call `InvalidateDbTableCache(dbPath)`.
3. In `cli/cmddb/cmd_record_del.go`:
   - In `executeRecordDeletion`: after successful commit, call `InvalidateDbTableCache(dbPath)`.
4. In `cli/cmddb/cmd_rm.go`:
   - In `executeAliasRemoval`: after successful alias removal, call `InvalidateDbTableCache(resolvedPath)`.
5. In `cli/cmddb/rest_handlers_crud.go` and `rest_handlers_ddl.go`:
   - In `handleInsertRecord`, `handleUpdateRecord`, `handleDeleteRecord`: extract database path from `mgr` and call `InvalidateDbTableCache(dbPath)`.
   - In `handleCreateTable`, `handleAlter`, `handleMigrate`, `handleClone`: call `InvalidateDbTableCache(dbPath)`.

### Step 2.3: Enhance `cli/cmddb/cmd_read.go` & `cmd_read_helpers.go`
1. Update `dispatchReadOperation`:
   - When `len(posArgs) == 1` and `posArgs[0]` is not a SQL statement, route to `handleReadInspect(flags, posArgs[0])` instead of basic `handleReadLs`.
2. Implement `handleReadInspect(flags readFlags, dbArg string) *appfault.AppError`:
   - Open read-only database connection via `openReadOnlyResolvedDb(dbArg)`.
   - Check `GetCachedTableMetadata(conn, dbPath)`.
   - If cache miss:
     - Query `sqlite_master` for all tables and views (`type IN ('table', 'view') AND name NOT LIKE 'sqlite_%'`).
     - Query schema DDL (`sql`).
     - Query column count via `PRAGMA table_info` and row count via `SELECT COUNT(*)`.
     - Save to cache via `SetCachedTableMetadata(conn, dbPath, tables)`.
   - Render output:
     - If `flags.isJson`: output structured JSON envelope with `operation: "db-read:inspect"`, tables array, column counts, row counts, DDL schema, and `isFromCache` boolean flag.
     - If terminal: render formatted table with columns `Table Name`, `Type`, `Columns`, `Rows`, `DDL Preview`.

### Step 2.4: Implement Unit Tests in `cli/cmddb/cmddb_table_cache_test.go`
1. `TestTableCacheHitAndMiss`: Verify initial query populates cache and second query returns cached data.
2. `TestTableCacheTtlExpiration`: Verify entries expire after configured TTL window.
3. `TestTableCacheExternalFileMutation`: Modify file `mtime`/size and verify cache detects stale entry.
4. `TestTableCachePragmaDataVersion`: Execute external write and verify `PRAGMA data_version` invalidates cache.
5. `TestTableCacheProactiveInvalidation`: Verify calling `InvalidateDbTableCache` immediately evicts the entry.
6. `BenchmarkTableCacheFreshnessValidation`: Benchmark `ValidateCacheFreshness` to verify execution time remains <0.04ms.

---

## 3. Verification & Quality Gates

- **Function Length**: <= 15 lines per function (<= 8 lines preferred).
- **Branching**: Zero nested `if` statements. Guard clauses only.
- **Naming**: Strict positive booleans (`is*` and `has*`).
- **Error Propagation**: Return direct `error` or `*appfault.AppError`. No calls to `os.Exit`.
- **Relative Path Hygiene**: Strictly relative paths in code, comments, and tests.
