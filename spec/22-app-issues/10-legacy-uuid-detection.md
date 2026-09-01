# Issue 10 — Legacy UUID Data Detection

## Problem

After migrating all primary keys from `TEXT` (UUID) to `INTEGER PRIMARY KEY AUTOINCREMENT` (`int64`), existing databases still contain legacy UUID strings. Queries scanning these rows fail with:

```
sql: Scan error on column index 0, name "Id": converting driver.Value type string ("...") to a int64: invalid syntax
```

## Root Cause

SQLite does not enforce column types strictly. Old rows written with UUID strings remain valid in the database even after schema migration, but the Go `database/sql` driver fails when scanning a string into an `int64` field.

## Solution

Added `isLegacyDataError(err error)` in `cmd/projectrepos.go` that checks for the `"converting driver.Value type string"` substring in the error message. All DB query paths that could encounter legacy rows now intercept this error and print a recovery prompt instead of a raw SQL error.

### Guarded Paths

| File                    | Function              | Query                          |
|-------------------------|-----------------------|--------------------------------|
| `cmd/projectrepos.go`   | `printProjectCount`   | `CountProjectsByTypeKey`       |
| `cmd/projectrepos.go`   | `printProjectList`    | `SelectProjectsByTypeKey`      |
| `cmd/list.go`           | `runList`             | `ListRepos` / `ShowGroup`      |
| `cmd/listreleases.go`   | `loadReleasesFromDB`  | `ListReleases`                 |
| `cmd/groupshow.go`      | `executeGroupShow`    | `ShowGroup`                    |
| `cmd/grouplist.go`      | `runGroupList`        | `ListGroups`                   |
| `cmd/stats.go`          | `loadStats`           | `QueryOverallStats`            |
| `cmd/history.go`        | `loadHistory`         | `ListHistory` / `ListHistoryByCommand` |
| `cmd/status.go`         | `loadRecordsByGroup`  | `ShowGroup`                    |
| `cmd/status.go`         | `loadAllRecordsDB`    | `ListRepos`                    |
| `cmd/export.go`         | `loadExportData`      | `ExportAll`                    |

### Recovery Message

Defined in `constants/constants_project.go` as `MsgLegacyProjectData`:

```
Database contains legacy project data from a previous version.
To fix, run one of:

  gitmap-v28 rescan          Re-scan repos and rebuild project data
  gitmap-v28 db-reset --confirm   Reset the entire database
```

## Status

Resolved. All critical query paths are guarded.
