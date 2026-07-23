---
name: scan-folder-and-version-probe
description: ScanFolder/VersionProbe schema (v3.7.0+). One row per scan-root absolute path; nullable Repo.ScanFolderId FK; gitmap sf add/list/rm subcommands. Hybrid HEAD-then-clone probe wired in Phase 2.3.
type: feature
---
# ScanFolder + VersionProbe (Phase 2.1 + 2.2)

## Schema (v3.7.0)

### `ScanFolder`
| Column | Type | Notes |
|---|---|---|
| ScanFolderId | INTEGER PK AUTOINCREMENT | |
| AbsolutePath | TEXT NOT NULL | UNIQUE via `IdxScanFolder_AbsolutePath` |
| Label | TEXT '' | Optional human label, settable via `--label` |
| Notes | TEXT '' | Free-form notes, settable via `--notes` |
| LastScannedAt | TEXT CURRENT_TIMESTAMP | Bumped on every `EnsureScanFolder` upsert |
| CreatedAt | TEXT CURRENT_TIMESTAMP | |

### `Repo` ALTER (idempotent)
- `ScanFolderId INTEGER DEFAULT NULL` — points at the *most recent* scan that discovered this repo. SQLite cannot add a `REFERENCES` clause via `ALTER`; the column stores the FK value without a declared FOREIGN KEY constraint, and store code enforces validity. **No backfill** — old repos stay NULL until the next `gitmap scan` re-discovers them.

### `VersionProbe` (table created in Phase 2.1, populated in Phase 2.3)
| Column | Type | Notes |
|---|---|---|
| VersionProbeId | INTEGER PK | |
| RepoId | INTEGER NOT NULL REFERENCES Repo(RepoId) ON DELETE CASCADE | |
| ProbedAt | TEXT CURRENT_TIMESTAMP | |
| NextVersionTag | TEXT '' | e.g. `v2`, `v3-alpha` |
| NextVersionNum | INTEGER 0 | |
| Method | TEXT '' | `head`, `clone`, `none` |
| IsAvailable | INTEGER 0 | Boolean: probe found a higher version |
| Error | TEXT '' | Probe failure detail |
- Index `IdxVersionProbe_RepoId` on `(RepoId, ProbedAt DESC)` for fast latest-probe lookups.

## CLI: `gitmap sf <add|list|rm>`

| Subcommand | Behaviour |
|---|---|
| `sf add <path> [--label X] [--notes Y]` | Resolves to absolute path, upserts ScanFolder row, prints `(id=N)` outcome. Re-running with the same path bumps `LastScannedAt` and only overwrites Label/Notes when the new values are non-empty. |
| `sf list` (alias `ls`) | Prints `[id] absolute-path` + `label / repos / last-scanned` per row. Newest-scanned first. Repo count is `SELECT COUNT(*) FROM Repo WHERE ScanFolderId = ?`. |
| `sf rm <path-or-id>` (alias `remove`) | Detaches every linked repo (`Repo.ScanFolderId = NULL`) then deletes the row. Reports `(id=N, R repos detached)`. |

`sf` is dispatched from `dispatchUtility`. Help line `HelpSf` appears in the **Navigation** group of `printUsage` (alongside group/multi-group/alias/diff-profiles), since scan folders are a navigation/scope concept rather than a data import.

## Files added/changed

| File | Change |
|---|---|
| `gitmap/constants/constants_scan_folder.go` (new) | Tables, indexes, CRUD SQL, error/message strings, CLI tokens |
| `gitmap/model/scan_folder.go` (new) | `ScanFolder` + `VersionProbe` types |
| `gitmap/store/scan_folder.go` (new) | `EnsureScanFolder`, `ListScanFolders`, `CountReposInScanFolder`, `RemoveScanFolderByPath/ByID` |
| `gitmap/store/store.go` | Wired `SQLCreateScanFolder/...PathIndex/VersionProbe/...RepoIndex` into `Migrate()` statement list; added `migrateRepoScanFolderID` ALTER step; added drops to `Reset()` |
| `gitmap/cmd/sf.go` (new) | `runSf` dispatcher, `runSfAdd/List/Remove`, helper `extractSfFlags`, `openSfDB` |
| `gitmap/cmd/rootutility.go` | Routes `CmdSf` to `runSf` |
| `gitmap/cmd/rootusage.go` | `printGroupNavigation` includes `HelpSf` |
| `gitmap/constants/constants_cli.go` | `CmdSf = "sf"` (top-level marker), `HelpSf` line |

## What's NOT in this phase
- No HEAD probe, no parallel `git clone --depth 1 --no-checkout` fallback (Phase 2.3).
- `ScanFolderId` is **not yet populated** by `gitmap scan` itself. Phase 2.3 will add `EnsureScanFolder(absDir, "", "")` to `executeScan` and pass the resulting id through `UpsertRepos`.
- `VersionProbe` table exists but has zero readers/writers yet.

## Reset ordering
`Reset()` drops `VersionProbe` and `ScanFolder` AFTER `RepoVersionHistory` and BEFORE `Repo`/`Repos` so the FK cascade order stays clean.
