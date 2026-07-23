# 90 — ScanFolder & VersionProbe schema (Phase 2.1 + 2.2)

## Goal

Track every absolute root that `gitmap-v27 scan` was invoked against, give each repo a pointer to the most recent root that discovered it, and add a `VersionProbe` table for the upcoming Phase 2.3 hybrid HEAD-then-clone version probe. Expose registration/inspection through `gitmap-v27 sf add | list | rm`.

Shipped in v3.7.0. No probe logic yet — that lands in Phase 2.3.

## Schema

### Table: `ScanFolder`

```sql
CREATE TABLE IF NOT EXISTS ScanFolder (
    ScanFolderId  INTEGER PRIMARY KEY AUTOINCREMENT,
    AbsolutePath  TEXT NOT NULL,
    Label         TEXT DEFAULT '',
    Notes         TEXT DEFAULT '',
    LastScannedAt TEXT DEFAULT CURRENT_TIMESTAMP,
    CreatedAt     TEXT DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX IF NOT EXISTS IdxScanFolder_AbsolutePath
  ON ScanFolder(AbsolutePath);
```

The unique index on `AbsolutePath` is what makes `EnsureScanFolder` an idempotent upsert. Re-running with the same path bumps `LastScannedAt` and only overwrites `Label`/`Notes` when the new value is non-empty (`CASE WHEN excluded.X = '' THEN ScanFolder.X ELSE excluded.X END`), so a manually set label survives subsequent automatic scans.

### Repo column addition

```sql
ALTER TABLE Repo ADD COLUMN ScanFolderId INTEGER DEFAULT NULL;
```

Added via `addColumnIfNotExists` so it's safe on every migration run. SQLite cannot attach a `REFERENCES` clause through `ALTER`, so this column stores the FK value without a declared `FOREIGN KEY` constraint. Application code (`removeScanFolderRow → SQLDetachReposFromScanFolder` then `SQLDeleteScanFolderByID`) enforces referential cleanup.

**No backfill.** Pre-v3.7.0 repos stay `ScanFolderId = NULL` until the next `gitmap-v27 scan` re-discovers them. This is intentional — back-filling would require guessing which historical scan-root each repo came from.

### Table: `VersionProbe`

```sql
CREATE TABLE IF NOT EXISTS VersionProbe (
    VersionProbeId  INTEGER PRIMARY KEY AUTOINCREMENT,
    RepoId          INTEGER NOT NULL REFERENCES Repo(RepoId) ON DELETE CASCADE,
    ProbedAt        TEXT DEFAULT CURRENT_TIMESTAMP,
    NextVersionTag  TEXT DEFAULT '',
    NextVersionNum  INTEGER DEFAULT 0,
    Method          TEXT DEFAULT '',     -- head | clone | none
    IsAvailable     INTEGER DEFAULT 0,
    Error           TEXT DEFAULT ''
);
CREATE INDEX IF NOT EXISTS IdxVersionProbe_RepoId
  ON VersionProbe(RepoId, ProbedAt DESC);
```

Empty in Phase 2.1 — created so Phase 2.3 doesn't need a separate migration. `(RepoId, ProbedAt DESC)` index keeps "latest probe per repo" lookups O(log n).

### Reset ordering

`Reset()` drops `VersionProbe` then `ScanFolder` between `RepoVersionHistory` and `Repo`/legacy `Repos`. FK cascade order is preserved.

## CLI surface

### `gitmap-v27 sf add <absolute-path> [--label X] [--notes Y]`
- Resolves the path with `filepath.Abs` (relative paths supported, but stored as absolute).
- Upserts via `EnsureScanFolder`.
- Prints `✓ Registered scan folder: <path> (id=N)` for first-time registration, or `✓ Scan folder already registered: <path> (id=N, last scanned T)` when the row already existed.

### `gitmap-v27 sf list` (alias `ls`)
- Prints `Scan folders (N):` header followed by one block per row:
  ```
  [3] /Users/foo/code
      label: work | repos: 12 | last scanned: 2026-04-19 14:32:01
  ```
- Empty state: `No scan folders registered. Run gitmap-v27 scan <dir> or gitmap-v27 sf add <dir>.`
- Sort: `LastScannedAt DESC, AbsolutePath ASC`.

### `gitmap-v27 sf rm <absolute-path|id>` (alias `remove`)
- Accepts either a path or a numeric id (`strconv.ParseInt`).
- Counts repos via `SELECT COUNT(*) FROM Repo WHERE ScanFolderId = ?` BEFORE detach.
- Detaches: `UPDATE Repo SET ScanFolderId = NULL WHERE ScanFolderId = ?`.
- Deletes: `DELETE FROM ScanFolder WHERE ScanFolderId = ?`.
- Prints `✓ Removed scan folder: <path> (id=N, R repos detached)`.

`sf` is routed from `dispatchUtility`; the navigation help group lists it next to `group` / `multi-group` / `alias` since it's a scope/navigation concept.

## Implementation map

| File | Role |
|---|---|
| `gitmap-v27/constants/constants_scan_folder.go` | All SQL, error formats, message formats, CLI tokens (`SFSubAdd`, `SFFlagLabel`, …) |
| `gitmap-v27/model/scan_folder.go` | `ScanFolder` + `VersionProbe` types |
| `gitmap-v27/store/scan_folder.go` | `EnsureScanFolder`, `ListScanFolders`, `CountReposInScanFolder`, `RemoveScanFolderByPath/ByID`, `removeScanFolderRow`, `findScanFolderByPath/ByID`, `scanOneScanFolder`, `scanScanFolderRows` |
| `gitmap-v27/store/store.go::Migrate` | Adds 4 CREATE statements + `migrateRepoScanFolderID` ALTER |
| `gitmap-v27/store/store.go::Reset` | Adds `SQLDropVersionProbe` and `SQLDropScanFolder` in FK order |
| `gitmap-v27/cmd/sf.go` | CLI dispatch + handlers (`runSfAdd/List/Remove`, `extractSfFlags`, `openSfDB`) |
| `gitmap-v27/cmd/rootutility.go` | `case CmdSf: runSf(os.Args[2:])` |
| `gitmap-v27/cmd/rootusage.go::printGroupNavigation` | Lists `HelpSf` |
| `gitmap-v27/constants/constants_cli.go` | `CmdSf = "sf"` + `HelpSf` |

## Out of scope (handled in later phases)

| Phase | Work |
|---|---|
| 2.3 | `executeScan` calls `EnsureScanFolder(absDir, "", "")`, threads the resulting id into `UpsertRepos`. Hybrid HEAD-then-clone probe runs after scan, populates `VersionProbe`. Scan blocks until all probes finish. Exit code reflects probe success. |
| 2.4 | `gitmap-v27 find-next` reads `VersionProbe` to print which repos have a higher version available. |
| 2.5 | `gitmap-v27 pull` parallelises pull + probe per scan folder. |
| 2.6 | `gitmap-v27 cn next all` bulk-clones every probed-available repo. |
| 2.7 | Final spec + ERD updates. |
