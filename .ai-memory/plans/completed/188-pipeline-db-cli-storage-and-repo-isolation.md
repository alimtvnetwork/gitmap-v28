# Plan 188: Pipeline Database CLI Co-Location, Directory Renaming & Repo-Slug Isolation

> **Status:** Completed
> **Category:** Architecture / Storage / Database
> **Completed At:** 2026-09-16

---

## 1. Problem Statement

In previous implementations, `ResolvePipelineDbPath` mistakenly routed pipeline SQLite databases into the repository-local `.gitmap/data/pipeline.db` path, causing:
1. Working-directory pollution (`.gitmap/data/pipeline.db`).
2. Loss of per-repository slug isolation across concurrent repositories.
3. Terminal summary displaying the relative path `.gitmap/data/pipeline.db` instead of the canonical CLI installation database directory.
4. Legacy directory naming inconsistency (`pipeline_db` vs `pipeline`).

The user explicitly mandated:
- Co-locate the pipeline database with the CLI binary installation database directory (`store.BinaryDataDir()`).
- Anchor databases inside a dedicated folder named `pipeline` (near `gitmap.db`).
- Isolate databases per repository based on the repo slug (e.g. `pipeline_<slug>.db` or `<slug>.db`).
- Display the canonical CLI database path in all pipeline summaries and error logs.
- Perform a release ceremony.

---

## 2. Root Cause Analysis (RCA)

1. **Working Directory Fallback**: `ResolvePipelineDbPath` previously called `resolveRepoScopedPath(root)`, defaulting to `<repoRoot>/.gitmap/data/pipeline.db`.
2. **Display Truncation**: `FormatRelativeDbPath` forced `.gitmap/data/pipeline.db` when an empty path was provided and looked for `.gitmap/` substring in paths.
3. **Directory Name Divergence**: The binary store used `pipeline_db/` while specification and user requirement called for `pipeline/`.

---

## 3. Implementation Details

### A. CLI Binary Data Directory Anchoring (`cli/pipelinedb/`)
- `PipelineDbDir()`: Points to `<BinaryDataDir>/pipeline/` (with fallback detection for `<binaryDir>/pipeline`).
- `ResolvePipelineDbPath(repoSlug string)`: Resolves to `<PipelineDbDir>/<slug>.db` or `<PipelineDbDir>/pipeline_<slug>.db`. Includes auto-migration from legacy `pipeline_db/` directory.
- `RepoScopedPipelineDbDir(_ string)`: Redirects to `PipelineDbDir()`.
- Decomposed schema execution into `cli/pipelinedb/pipeline_split_schema_exec.go` and connection management into `cli/pipelinedb/pipeline_split_conn.go`, keeping all files strictly $\le 100$ lines.

### B. Display Path Formatting (`cli/cmdpipeline/pipeline_sync_cache.go`)
- `FormatRelativeDbPath`: Returns repo-relative path ONLY when the file is genuinely located inside the repository root. When anchored to the CLI installation directory, returns the clean absolute path.
- `resolveDbStatTarget`: Defaults to `pipelinedb.ResolvePipelineDbPath(resolveCurrentRepoSlug())`.

### C. Registry and Inventory Synchronization (`cli/store/`)
- `syncPipelineDBs()`: Scans `<dataDir>/pipeline/` first, followed by legacy `pipeline_db/`.
- `storage_inventory.go`: Collects databases from both `pipeline` and `pipeline_db` directories.
- `cmddb_status.go`: Updates section header to `3. Split Pipeline Databases (pipeline/):`.

### D. Specifications & Helptext
- Updated `02-spec/21-app/10-pipeline-and-repo-split-db/01-pipeline-split-db-architecture.md` to reference `<BinaryDataDir>/pipeline/`.
- Updated `cli/helptext/agy-fix-pipeline.md`.

---

## 4. Verification

- Ran `go vet ./pipelinedb/... ./cmdpipeline/... ./store/... ./cmddb/...`: 100% clean (exit code 0).
- Generated test inventory: 3,536 tests indexed across 139 packages.
- Migrated 103MB pipeline database from `.gitmap/data/pipeline.db` to `C:\Users\Administrator\AppData\Local\gitmap-cli\data\pipeline\pipeline_alimtvnetwork-gitmap-v28.db`.
