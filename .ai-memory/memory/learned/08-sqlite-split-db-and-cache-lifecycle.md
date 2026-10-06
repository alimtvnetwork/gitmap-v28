# 08 — SQLite Split-DB Architecture and Cache Lifecycle

- **Subsystem:** Storage Engine & Split Database Design
- **Status:** Authoritative Reference

## 1. Three-Tier Split-DB Architecture
- `gitmap.db`: Core repository catalog, machine settings, and persistent aliases.
- `installation.db`: Package versions, pinned releases, and installer metadata.
- `pipeline.db` / `repodb/*.db`: Transient CI/CD telemetry, build stage timings, and repository-scoped git caches.

## 2. Concurrency & Performance
- Enables Write-Ahead Logging (WAL) and memory-mapped I/O (`PRAGMA mmap_size`).
- Reentrant mutex locks protect connection pools against database locking errors.
- Prepared SQL statements eliminate query parsing overhead.
