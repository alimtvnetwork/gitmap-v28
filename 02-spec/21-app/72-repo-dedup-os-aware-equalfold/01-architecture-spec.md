# Architecture Specification: Redundant Repos, EqualFold, OS Path Sensitivity & SQLite Uniqueness

- **Task Slug:** 72-repo-dedup-os-aware-equalfold
- **Version:** 6.464.0
- **Scope:** Cross-platform GitMap CLI, SQLite split-db engine, path sensitivity, and string comparison efficiency.

## 1. Architectural Principles

1. **Zero-Allocation String Equality:**
   - Prefer `strings.EqualFold(a, b)` over `strings.ToLower(a) == strings.ToLower(b)` or `strings.ToLower(a) == "literal"` to eliminate memory allocations in hot loops and flag evaluation.
   - For multi-candidate checks, use centralized `strutil.EqualFoldAny` and `strutil.EqualFoldAnyTrim`.

2. **OS-Aware Path Sensitivity:**
   - Windows filesystems are case-insensitive. Canonical path keys fold to lowercase, and SQLite indexes use `COLLATE NOCASE`.
   - Unix filesystems are case-sensitive. Paths differing only by case are legally distinct. Byte-level case preservation must be strictly enforced. SQLite indexes on Unix must use standard binary collation without `COLLATE NOCASE`.

3. **Database-Driven Redundancy Pruning:**
   - Redundancy checks must execute directly against SQLite (`gitmap.db` Repo table and `repo_search` databases) rather than traversing physical filesystems.
   - Redundant repositories are cleaned before pulling worker channels via `OptimizeRedundantRepos` in `cli/cmdpull/pull_dedup.go`.

4. **Guaranteed Ingestion Uniqueness:**
   - Tables tracking repositories and files must enforce unique indexes and conflict-resolution upserts (`ON CONFLICT(...) DO UPDATE SET`).
   - `RepoFile` in Split DBs uses `ON CONFLICT(RelativePath) DO UPDATE SET` via `RepoFileDbRepo.Upsert()`.
