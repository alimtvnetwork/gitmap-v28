# Issue 12 — Legacy UUID to Integer ID Migration

## Problem

After migrating the codebase from UUID (`TEXT`) to `INTEGER PRIMARY KEY AUTOINCREMENT` (`int64`), existing databases retained the old schema because `CREATE TABLE IF NOT EXISTS` never alters existing tables. This caused:

1. `sql: Scan error` — scanning TEXT UUID into `int64` fields
2. `FOREIGN KEY constraint failed (787)` — `alignRecordsWithDB` silently failed (returning ID=0), so `DetectedProjects.RepoId` referenced a non-existent row

## Root Cause

SQLite's `CREATE TABLE IF NOT EXISTS` is a no-op when the table already exists. Old `Repos.Id` remained `TEXT` even though the new schema specifies `INTEGER PRIMARY KEY AUTOINCREMENT`. Upserts worked (they don't touch `Id`), but all ID-dependent reads and FK references broke.

## Solution

Added `store/migrateids.go` with `migrateLegacyIDs()`, called at the start of `Migrate()`:

1. **Detection**: `PRAGMA table_info(Repos)` checks if `Id` column type is `TEXT`
2. **Drop dependents**: Drops all FK-dependent tables (project detection, GroupRepos)
3. **Rebuild**: Renames `Repos` → `Repos_legacy`, creates new `Repos` with integer IDs, copies all data (new auto-increment IDs assigned), drops legacy table
4. **FK safety**: Temporarily disables `PRAGMA foreign_keys` during rebuild

### Data preserved

- All repo records (Slug, RepoName, paths, URLs, etc.)
- Groups (no FK to Repos)
- Releases (no FK to Repos)

### Data reset (repopulated by next scan)

- GroupRepos associations (FK to Repos.Id)
- DetectedProjects and all metadata tables (FK chain to Repos.Id)
- ProjectTypes (re-seeded by Migrate)

## Related

- Issue 10: Legacy UUID data detection (graceful error messages)
- v2.36.1: Fix release

## Status

Resolved.
