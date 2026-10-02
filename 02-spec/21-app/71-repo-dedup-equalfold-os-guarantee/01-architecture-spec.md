# Architecture Specification: Redundant Repos, EqualFold, OS Path Sensitivity & SQLite Uniqueness

- **Task Slug:** 71-repo-dedup-equalfold-os-guarantee
- **Version:** 6.462.0
- **Scope:** Cross-platform GitMap CLI, SQLite split-db engine, and CLI command options.

## 1. Architectural Principles
1. **Zero-Allocation String Equality:** Prefer `strings.EqualFold(a, b)` over `strings.ToLower(a) == strings.ToLower(b)` or `strings.ToLower(a) == "literal"` to eliminate memory allocations in hot loops and flag evaluation.
2. **OS-Aware Path Sensitivity:**
   - Windows filesystems are case-insensitive. Canonical path keys fold to lowercase, and SQLite indexes use `COLLATE NOCASE`.
   - Unix filesystems are case-sensitive. Paths differing only by case are legally distinct. Byte-level case preservation must be strictly enforced.
3. **Database-Driven Redundancy Pruning:**
   - Redundancy checks must execute directly against SQLite (`gitmap.db` Repo table and `repo_search` databases) rather than traversing physical filesystems.
4. **Guaranteed Ingestion Uniqueness:**
   - Tables tracking repositories and files must enforce unique indexes and conflict-resolution upserts (`ON CONFLICT(...) DO UPDATE SET`).
