# Subtask 94.02: repodb & db ORM Connection, Entity Generation, and PascalCase Id Normalization

## Goal
Declare strongly-typed model structs for `gitmap/repodb` (`RepoFile`, `SearchCache`, `FileSequence`, `SequenceHistory`, `RepoScanLog`, `IndexedRepo`), normalize table primary keys to PascalCase `<Entity>Id` (`RepoFileId`, `SearchCacheId`, etc.), generate typed enums and repositories via the code generator, and connect them with `dbengine.DbWrapper`.

## Files Impacted
- `gitmap/repodb/repo_db.go`
- `gitmap/repodb/models.go` (NEW)
- `gitmap/repodb/enums/` (NEW)
- `gitmap/repodb/consts.go` (NEW)
- `gitmap/repodb/root.go`

## Acceptance Criteria
1. PascalCase primary keys: `RepoFileId`, `SearchCacheId`, `FileSequenceId`, `SequenceHistoryId`, `RepoScanLogId`, `IndexedRepoId` (zero bare `Id` or all-caps `ID`).
2. Generates type-safe column enums (`*FieldType`), singleton registries (`*DbRegistry`), null-safe row scanners (`Scan*`), and typed repositories (`*DbRepo`) using `30-db-struct-enum-generator.py`.
3. Integrates `dbengine.WrapDb(conn, dbengine.DbSQLite)` in `repodb.OpenRepoDB`.
4. All functions $\le 15$ lines, zero nested ifs, affirmative booleans only.
