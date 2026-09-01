# 22 — Scan Release Import

## Purpose

Enhance the `scan` command to discover `.gitmap/release/*.json` metadata files
in the scan root and upsert them into the `Releases` database table.
This keeps the DB in sync with on-disk release history without requiring
users to run `gitmap release` for every past version.

## Trigger

The import runs automatically at the end of every `gitmap scan`,
after repos are upserted and before the output folder is opened.

## Behavior

1. Resolve the `.gitmap/release/` directory path relative to the scan root.
2. If the directory does not exist, skip silently (no error).
3. Read all files matching the glob `.gitmap/release/v*.json`.
   Skip `latest.json` and any non-`v`-prefixed files.
4. For each file, unmarshal into `release.ReleaseMeta`.
   Skip files that fail to parse (log a warning, continue).
5. Map each `ReleaseMeta` to a `model.ReleaseRecord`.
6. Call `db.UpsertRelease()` for each record.
7. Print a summary: `"Releases imported: %d from .gitmap/release/\n"`.
8. If zero files were found, print nothing (silent skip).

## Field Mapping

| ReleaseMeta field | ReleaseRecord field | Notes                        |
|-------------------|---------------------|------------------------------|
| Version           | Version             | Direct copy                  |
| Tag               | Tag                 | Direct copy                  |
| Branch            | Branch              | Direct copy                  |
| SourceBranch      | SourceBranch        | Direct copy                  |
| Commit            | CommitSha           | Renamed field                |
| Changelog         | Changelog           | Join notes with `\n` if list |
| Draft             | Draft               | Direct copy                  |
| PreRelease        | PreRelease          | Direct copy                  |
| IsLatest          | IsLatest            | Direct copy                  |
| CreatedAt         | CreatedAt           | Direct copy                  |

## Edge Cases

| Condition                        | Behavior                                    |
|----------------------------------|---------------------------------------------|
| `.gitmap/release/` missing              | Skip silently, no error                     |
| File fails to parse              | Log warning, skip file, continue            |
| Duplicate tag in DB              | Upsert overwrites (ON CONFLICT(Tag))        |
| `latest.json` in directory       | Ignored (not a release file)                |
| Pre-release files (e.g. rc)      | Imported normally, `IsLatest = false`        |
| Empty `.gitmap/release/` directory      | Skip silently, no output                    |

## Implementation Files

| File                            | Responsibility                               |
|---------------------------------|----------------------------------------------|
| `cmd/scan.go`                   | Call `importReleases()` after `upsertToDB`   |
| `cmd/scanimport.go` (new)       | `importReleases()` orchestration             |
| `release/metadata.go`           | `ReadReleaseMeta()` (new: read single file)  |
| `constants/constants_messages.go` | `MsgReleasesImported` format string        |

## New Functions

### `release/metadata.go`

```
ReadReleaseMeta(path string) (ReleaseMeta, error)
```

Reads and unmarshals a single `.gitmap/release/vX.Y.Z.json` file.

### `cmd/scanimport.go`

```
importReleases(scanDir, outputDir string)
```

Orchestrates: glob → parse → map → upsert → summary.

```
mapMetaToRecord(m ReleaseMeta) model.ReleaseRecord
```

Maps `ReleaseMeta` fields to `ReleaseRecord`.

## Integration Point

In `cmd/scan.go`, add the call after `upsertToDB`:

```go
upsertToDB(records, outputDir)
importReleases(absDir, outputDir)   // ← new
```

## Code Style

All functions ≤ 15 lines. Positive logic. Blank line before every return.
No magic strings. No switch statements. PascalCase for SQL column names.
