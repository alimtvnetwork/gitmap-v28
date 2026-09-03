# Unified Database Status and Optimization

## 1. Overview

The `gitmap db` command hierarchy provides an operational control plane over all SQLite databases in the GitMap ecosystem:

1. **Master SQLite Database:** `gitmap.db` (Primary metadata, repositories, schedules, version manifests).
2. **Split Repository Databases:** `repo_search/<slug>-<id>.db` (Isolated file sequence indexes and full-text search caches).
3. **Split Pipeline Databases:** `pipeline_db/pipeline_<slug>.db` (Isolated CI/CD run baselines, segment logs, error diagnostics).

## 2. Command Specifications

### `gitmap db status`

Prints a consolidated health overview across all three tiers:
- Primary Master DB location, size, and registered repository count.
- Split Repo DB count, total disk space consumed, and directory path.
- Split Pipeline DB count, total disk space consumed, and directory path.
- Supports `--json` for machine-readable status pipelines.

### `gitmap db optimize`

Iteratively executes vacuum, WAL truncation, and table analysis across:
- Primary database (`gitmap.db`)
- Every discovered database in `repo_search/`
- Every discovered database in `pipeline_db/`

Reports before/after disk footprints and total bytes reclaimed.
