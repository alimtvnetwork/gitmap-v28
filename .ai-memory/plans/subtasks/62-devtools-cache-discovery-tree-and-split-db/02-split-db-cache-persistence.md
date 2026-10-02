# Subtask 02: Split-DB Cache Persistence & Force Invalidation

> **Parent Plan:** [62-devtools-cache-discovery-tree-and-split-db](../../completed/62-devtools-cache-discovery-tree-and-split-db.md)
> **Tracking Spec:** [199-devtools-cache-dynamic-discovery-tree-view-and-split-db.md](../../../../02-spec/21-app/199-devtools-cache-dynamic-discovery-tree-view-and-split-db.md)
> **Primary File Targets:** `cli/store/devtools_cache_split_db.go`, `cli/store/devtools_cache_types.go`
> **Status:** `COMPLETED`
> **Owner:** Worker 02

---

## 1. Objective & Context

Dynamic developer cache discovery involves querying external runtime CLI binaries (`go env`, `npm config get cache`, `pnpm store path`, `pip cache dir`, `cargo cache --dir`, etc.) and probing multi-drive heuristics (e.g. `C:\dev-tool\*`, `D:\dev-tool\*`). Executing full discovery on every invocation introduces avoidable subprocess and filesystem traversal latency.

To optimize performance while preserving accuracy:
1. **Split-DB Persistence:** Discovered paths, ecosystem tags, sizes, and file metrics must be persisted into a dedicated SQLite Split-DB table (`devtools_cache_paths`) located at `.gitmap/data/devtools/cache/sql.db`.
2. **Sub-second Re-clean / Dry-run:** On subsequent runs, cached paths are retrieved immediately from Split-DB without repeating expensive external subprocess calls.
3. **`--force` (`-f`) Invalidation:** When the user explicitly requests cache invalidation via `--force`, the database records are invalidated or purged, triggering fresh multi-tier discovery and updating Split-DB with new probe results.

---

## 2. SQLite Schema Specification

The database file is resolved using GitMap's canonical split-database resolution helper:
`ResolveSplitDbPath("devtools", "cache", "")` -> `.gitmap/data/devtools/cache/sql.db`.

### 2.1 Table Definition (`devtools_cache_paths`)

```sql
CREATE TABLE IF NOT EXISTS devtools_cache_paths (
    id               INTEGER PRIMARY KEY AUTOINCREMENT,
    path             TEXT NOT NULL UNIQUE,
    ecosystem        TEXT NOT NULL,
    size_bytes       INTEGER NOT NULL DEFAULT 0,
    files_count      INTEGER NOT NULL DEFAULT 0,
    dirs_count       INTEGER NOT NULL DEFAULT 0,
    last_verified_at INTEGER NOT NULL DEFAULT 0,
    is_custom        INTEGER NOT NULL DEFAULT 0,
    is_active        INTEGER NOT NULL DEFAULT 1,
    created_at       INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
    updated_at       INTEGER NOT NULL DEFAULT (strftime('%s', 'now'))
);

CREATE INDEX IF NOT EXISTS idx_devtools_cache_ecosystem ON devtools_cache_paths(ecosystem);
CREATE INDEX IF NOT EXISTS idx_devtools_cache_active ON devtools_cache_paths(is_active);
```

### 2.2 Column Semantics

| Column | Type | Constraints | Description |
| :--- | :--- | :--- | :--- |
| `id` | INTEGER | PRIMARY KEY AUTOINCREMENT | Unique row identifier |
| `path` | TEXT | NOT NULL UNIQUE | Canonical slash-separated absolute directory path (e.g. `C:/dev-tool/go/cache`) |
| `ecosystem` | TEXT | NOT NULL | Target ecosystem identifier (`go`, `npm`, `pnpm`, `pip`, `cargo`, `dotnet`, `gradle`, etc.) |
| `size_bytes` | INTEGER | NOT NULL DEFAULT 0 | Last measured disk usage in bytes prior to purge |
| `files_count` | INTEGER | NOT NULL DEFAULT 0 | Count of files residing in the target directory |
| `dirs_count` | INTEGER | NOT NULL DEFAULT 0 | Count of subdirectories residing in the target directory |
| `last_verified_at`| INTEGER | NOT NULL DEFAULT 0 | Epoch timestamp (seconds) when path existence was last confirmed |
| `is_custom` | INTEGER | NOT NULL DEFAULT 0 | `1` for non-default paths (e.g. `C:\dev-tool\...`), `0` for OS defaults |
| `is_active` | INTEGER | NOT NULL DEFAULT 1 | `1` for valid active cache targets; `0` when disabled or soft-deleted |
| `created_at` | INTEGER | NOT NULL DEFAULT epoch | Creation timestamp in Unix epoch seconds |
| `updated_at` | INTEGER | NOT NULL DEFAULT epoch | Last modification timestamp in Unix epoch seconds |

---

## 3. Go Models & Split-DB Method Signatures

To comply with the coding guideline constraint of `<= 100 lines per file` and `<= 15 lines per function`, the implementation is cleanly modularized into:
1. `cli/store/devtools_cache_types.go` (Domain model & DTOs)
2. `cli/store/devtools_cache_split_db.go` (Database connection, schema init, CRUD operations)

### 3.1 Domain Structs (`devtools_cache_types.go`)

```go
package store

// DevtoolsCacheRecord represents a persisted cache directory entry.
type DevtoolsCacheRecord struct {
	ID             int64  `json:"id"`
	Path           string `json:"path"`
	Ecosystem      string `json:"ecosystem"`
	SizeBytes      int64  `json:"sizeBytes"`
	FilesCount     int    `json:"filesCount"`
	DirsCount      int    `json:"dirsCount"`
	LastVerifiedAt int64  `json:"lastVerifiedAt"`
	IsCustom       bool   `json:"isCustom"`
	IsActive       bool   `json:"isActive"`
	CreatedAt      int64  `json:"createdAt"`
	UpdatedAt      int64  `json:"updatedAt"`
}

// DevtoolsCacheFilter defines filtering options for querying cached paths.
type DevtoolsCacheFilter struct {
	Ecosystems []string
	IsActiveOnly bool
}
```

### 3.2 Method Signatures (`devtools_cache_split_db.go`)

All functions return structured `*appfault.AppError` for predictable error propagation:

```go
// OpenDevtoolsCacheSplitDB opens the split SQLite database at .gitmap/data/devtools/cache/sql.db.
func OpenDevtoolsCacheSplitDB() (*sql.DB, *appfault.AppError)

// EnsureDevtoolsCacheSchema creates the devtools_cache_paths table and indexes if missing.
func EnsureDevtoolsCacheSchema(db *sql.DB) *appfault.AppError

// SaveDiscoveredDevtoolsPaths saves or updates discovered paths using an atomic transaction.
func SaveDiscoveredDevtoolsPaths(db *sql.DB, records []DevtoolsCacheRecord) *appfault.AppError

// GetActiveDevtoolsPaths returns all active cache paths, optionally filtered by ecosystems.
func GetActiveDevtoolsPaths(db *sql.DB, ecosystems []string) ([]DevtoolsCacheRecord, *appfault.AppError)

// InvalidateDevtoolsCache marks existing records inactive or wipes them for a fresh scan.
func InvalidateDevtoolsCache(db *sql.DB) *appfault.AppError

// UpdateDevtoolsPathMetrics updates the size, file counts, and verification time for a path.
func UpdateDevtoolsPathMetrics(db *sql.DB, path string, sizeBytes int64, files, dirs int) *appfault.AppError
```

---

## 4. Invalidation & Refresh Lifecycle (`--force`)

```
      +-----------------------------+
      |  User runs:                 |
      |  gitmap clear devtools [-f] |
      +--------------+--------------+
                     |
            Is --force present?
            /               \
         YES                 NO
         /                     \
+--------------------+   +--------------------------------+
| Invalidate Cache:  |   | Query Split-DB:                |
| InvalidateDevtools |   | GetActiveDevtoolsPaths()       |
| Cache()            |   +---------------+----------------+
+---------+----------+                   |
          |                      Are records found?
          |                      /                \
          |                   YES                  NO
          |                   /                      \
          v                  v                        v
+-------------------------------+             +-------------------------------+
| Run Multi-Tier Deep Discovery |             | Use Cached Paths Directly:    |
| (Worker 01 Engine)            |             | Skip external CLI probing     |
+---------------+---------------+             +-------------------------------+
                |
+---------------+---------------+
| Persist to Split-DB:          |
| SaveDiscoveredDevtoolsPaths() |
+-------------------------------+
```

---

## 5. Architectural & Coding Guideline Guardrails

1. **File Size Limit (< 100 Lines):** Keep `devtools_cache_split_db.go` and `devtools_cache_types.go` strictly under 100 lines. Extract query helpers if line count approaches threshold.
2. **Function Line Limit (<= 15 Lines):** Every function must not exceed 15 lines of code. Split complex SQL loops or scanning logic into small private helper functions.
3. **Affirmative Boolean Naming:** Use `IsCustom`, `IsActive`, `IsActiveOnly` (no negative booleans like `DisableCache` or `SkipCustom`).
4. **Normalized Paths:** Store paths with forward slashes (`filepath.ToSlash`) to avoid Windows backslash escaping inconsistencies in SQLite queries.
5. **WAL & Pragmas:** Inherit canonical pragmas (`_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)&_pragma=synchronous(NORMAL)`) via `openSqliteConn`.

---

## 6. Verification & Acceptance Criteria

- **AC-1:** Invoking `OpenDevtoolsCacheSplitDB()` ensures directory `.gitmap/data/devtools/cache/` is created and initializes `devtools_cache_paths` with valid indexes.
- **AC-2:** `SaveDiscoveredDevtoolsPaths()` successfully inserts discovered paths (e.g. `C:/dev-tool/go/cache`, `C:/dev-tool/go/pkg/mod`) with `is_custom = 1`.
- **AC-3:** `GetActiveDevtoolsPaths()` filters by ecosystem correctly when `--only` is specified (e.g. `go,npm`).
- **AC-4:** When `--force` is passed, `InvalidateDevtoolsCache()` resets cache entries, forcing deep discovery to re-populate fresh data.
- **AC-5:** No `go build` or `go test` invocations per Rule R1; passes `python 03-ai-scripts/05-guideline-autofixer.py` with zero guideline violations.
