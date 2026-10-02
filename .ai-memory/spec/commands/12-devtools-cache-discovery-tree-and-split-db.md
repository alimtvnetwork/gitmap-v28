# Command Specification: Devtools Cache Dynamic Discovery, Tree View Rendering & Split-DB Persistence

Spec Reference: [02-spec/21-app/199-devtools-cache-dynamic-discovery-tree-view-and-split-db.md](../../../02-spec/21-app/199-devtools-cache-dynamic-discovery-tree-view-and-split-db.md)
Plan Reference: [.ai-memory/plans/pending/62-devtools-cache-discovery-tree-and-split-db.md](../../plans/pending/62-devtools-cache-discovery-tree-and-split-db.md)

## 1. Overview

This specification establishes the commands, flags, and operational contracts for:
1. Dynamic 3-tier developer cache discovery across CLI tools, environment variables, and filesystem heuristics (including `C:\dev-tool\*`, `D:\dev-tool\*`, `GOCACHE`, `GOMODCACHE`, `pnpm store path`).
2. High-speed Split-DB caching of discovered cache directories in `.gitmap/data/devtools/cache/sql.db` with `--force` cache invalidation.
3. ANSI colored category cards with status badges (`[✔ Cleaned]`, `[⚡ Reclaimable]`, `[💾 Cached]`, `[🔄 Fresh Scan]`).
4. Hierarchical directory tree visualization via `--tree` (`-t`) rendering parent nodes, sub-caches, file counts, and byte sizes.
5. Accurate pre-sweep disk size calculation eliminating the 0.01 MB metric anomaly.

---

## 2. Command Index & Usage

### 2.1 Devtools Cache Cleaning & Discovery
- **Commands:** `gitmap clear devtools`, `gitmap clean devtools`, `gitmap devtool clear`, `gitmap dev clean`, `gitmap clean-dev`
- **Description:** Sweeps, sizes, and cleans developer tool caches (Go, npm, pnpm, yarn, bun, pip, cargo, dotnet, gradle, maven).
- **Flags:**
  - `--force` (`-f`): Invalidates the Split-DB cache and executes a fresh deep multi-tier filesystem probe.
  - `--tree` (`-t`): Renders an expandable hierarchical directory tree showing directories, file counts, and byte sizes.
  - `--dry-run` (`-n`, `-d`): Previews reclaimable space and files without deleting anything.
  - `--yes` (`-y`, `/y`): Bypasses interactive confirmation for automated workflows and AI subagents.
  - `--only <ecosystems>`: Filters targets to specified comma-delimited ecosystems (e.g. `--only go,pnpm`).
  - `--json`: Outputs structured JSON payload with per-ecosystem metrics and path lists.
  - `-v`, `--verbose`: Prints individual path diagnostics and skipped entries.
  - `-h`, `--help`: Displays full command documentation and flag usage.

---

## 3. Split-DB Storage Schema

Stored at `.gitmap/data/devtools/cache/sql.db`:
- **Table:** `devtools_cache_paths`
- **Fields:** `id`, `path`, `ecosystem`, `size_bytes`, `files_count`, `dirs_count`, `last_verified_at`, `is_custom`, `is_active`, `created_at`, `updated_at`.
- **Query speed:** Sub-millisecond retrieval on warm cache, bypassing repeated external CLI queries.
