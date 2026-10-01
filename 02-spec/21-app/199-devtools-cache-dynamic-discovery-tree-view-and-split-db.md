# Spec 199: Devtools Cache Dynamic Discovery, Tree View Rendering & Split-DB Persistence

## 1. Overview & Problem Statement

Users running `gitmap clear devtools` (aliases: `gitmap clean devtools`, `gitmap devtool clear`) identified critical operational gaps:
1. **0.01 MB Metric Anomaly:** Wiping commands (`go clean`, `npm cache clean`) execute prior to directory measurement, zeroing out files before size accounting occurs.
2. **Missing Custom Cache Paths:** Standard clean sweeps only inspect static default paths (e.g. `go/pkg/mod`) and fail to locate real developer caches (e.g. `C:\dev-tool\go\cache`, `C:\dev-tool\go\pkg\mod`, custom `GOCACHE`, `GOMODCACHE`, `D:\dev-tool\pnpm\store\v10`).
3. **No Dynamic CLI or Environment Querying:** Tools such as `go env`, `npm config get cache`, `pnpm store path`, `pip cache dir`, and environment variables like `GOCACHE`, `GOPATH`, `CARGO_HOME`, `PIP_CACHE_DIR` are bypassed.
4. **Uncolored & Opaque Terminal Output:** Output renders a plain, uncolored flat table with no ANSI styling, does not display which specific paths are being deleted, and lacks an expandable directory tree view.
5. **Absence of Cache Persistence with `--force` Reset:** Path discovery overhead repeats on every invocation unless cached in a dedicated SQLite Split-DB with a `--force` (`-f`) flag for user-directed cache re-scanning.

---

## 2. Technical Specifications

### 2.1 Dynamic Discovery Engine (3-Tier Resolution)
- **Tier 1 — Dynamic CLI Query:**
  - Go: `go env GOCACHE`, `go env GOMODCACHE`, `go env GOPATH`
  - Node: `pnpm store path`, `npm config get cache`, `yarn cache dir`, `bun pm cache`
  - Python: `pip cache dir`
  - Rust: `cargo cache --dir`
  - .NET: `dotnet nuget locals all -l`
- **Tier 2 — Environment Variables:**
  - `GOCACHE`, `GOMODCACHE`, `GOPATH`
  - `PNPM_HOME`, `npm_config_cache`, `YARN_CACHE_FOLDER`, `BUN_INSTALL`
  - `PIP_CACHE_DIR`, `UV_CACHE_DIR`, `POETRY_CACHE_DIR`
  - `CARGO_HOME`, `RUSTUP_HOME`
  - `NUGET_PACKAGES`, `GRADLE_USER_HOME`, `M2_HOME`
- **Tier 3 — Drive & Filesystem Heuristics:**
  - Probing drives `C:\`, `D:\`, `E:\` for `<drive>:\dev-tool\*` (e.g. `C:\dev-tool\go\cache`, `C:\dev-tool\go\pkg\mod`, `D:\dev-tool\pnpm\store`)
  - Standard OS cache directories: `%LOCALAPPDATA%\go-build`, `%LOCALAPPDATA%\npm-cache`, `~/.cargo/registry/cache`, `~/.gradle/caches`, `~/.m2/repository`.

### 2.2 Split-DB Cache Persistence (`cli/store/devtools_cache_split_db.go`)
- Database path: `.gitmap/data/devtools/cache/sql.db` (resolved via `ResolveSplitDbPath("devtools", "cache", "")`).
- SQLite Table:
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
- Methods: `OpenDevtoolsCacheSplitDB()`, `EnsureDevtoolsCacheTable()`, `SaveDiscoveredPaths()`, `GetCachedPaths()`, `InvalidateCache()`.
- When `--force` (`-f`) is passed, invalidates cached entries and re-executes deep discovery.

### 2.3 ANSI Colored Output & Tree View (`--tree`, `-t`)
- Rich colored category cards with ANSI status badges: `[✔ Cleaned]`, `[⚡ Reclaimable]`, `[💾 Cached]`, `[🔄 Fresh Scan]`.
- Detailed path listings showing exact directory paths targeted with byte sizes.
- Hierarchical tree view (`--tree`, `-t`):
  - Renders connector branches (`├──`, `└──`, `│  `).
  - Itemizes parent category nodes, root paths, sub-caches (e.g. `cache`, `testcache`, `pkg/mod`), and individual file counts.

---

## 3. CLI Flags & Contract

| Flag | Shorthand | Description |
|---|---|---|
| `--force` | `-f` | Invalidate Split-DB cache and force fresh deep filesystem discovery |
| `--tree` | `-t` | Render hierarchical directory tree view with file counts and byte sizes |
| `--dry-run` | `-n`, `-d` | Preview files and space to be freed without deleting anything |
| `--yes` | `-y`, `/y` | Bypass confirmation prompt |
| `--only` | — | Filter target ecosystems (e.g. `go,npm,pnpm,cargo`) |
| `--json` | — | Output structured JSON for automation and AI agents |

---

## 4. Verification Gates
1. Rule R1: Zero `go build` or `go test` calls.
2. Rule R11: Strict relative Git paths and lowercase conventions.
3. Coding Guidelines: All functions <= 15 lines, files < 100 lines, affirmative booleans only.
