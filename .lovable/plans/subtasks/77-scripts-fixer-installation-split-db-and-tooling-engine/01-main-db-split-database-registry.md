# Subtask 01: Main DB Split Database Registry

**Parent Plan:** [77-scripts-fixer-installation-split-db-and-tooling-engine.md](../../pending/77-scripts-fixer-installation-split-db-and-tooling-engine.md)
**Status:** Completed
**Target Files:**
- `gitmap/constants/constants_split_db_sql.go`
- `gitmap/store/split_database_registry.go`
- `gitmap/store/split_database_registry_test.go`
- `gitmap/store/store.go`

## Objectives
1. Define DDL `SQLCreateSplitDatabaseRegistry` in `gitmap/constants/constants_split_db_sql.go`.
2. Implement `SplitDatabaseEntry` model and store methods in `gitmap/store/split_database_registry.go`:
   - `RegisterSplitDB(entry SplitDatabaseEntry) error`
   - `GetSplitDB(dbType, dbKey string) (*SplitDatabaseEntry, error)`
   - `ListSplitDBs(dbType string) ([]SplitDatabaseEntry, error)`
   - `SyncKnownSplitDatabases() error`
3. Integrate schema migration and automatic split DB synchronization in `gitmap/store/store.go:Migrate()`.
4. Add unit test `gitmap/store/split_database_registry_test.go` verifying registration, lookup, conflict updates, and sync.
