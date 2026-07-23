# Issue 09 — list-releases reads from DB instead of current repo

## Problem

Running `gitmap-v28 lr` inside a git repo with `.gitmap/release/v*.json` files would
ignore those files and query only the gitmap-v28 SQLite database. This meant
releases created in the current repo were invisible unless a `gitmap-v28 scan`
had been run to import them.

## Root Cause

`loadReleases()` in `cmd/listreleases.go` unconditionally opened the
database via `openDB()` and called `db.ListReleases()`. It never checked
for `.gitmap/release/` files in the working directory.

## Fix

Changed `loadReleases()` to use a dual-source resolution:

1. **Repo-first**: call `release.ListReleaseMetaFiles()` to read all
   `.gitmap/release/v*.json` files. Convert each `ReleaseMeta` to a
   `model.ReleaseRecord` with `Source = "repo"`, sort by `CreatedAt DESC`,
   and mark `IsLatest` from `latest.json`.
2. **DB fallback**: only if no `.gitmap/release/` files are found, fall back to
   `db.ListReleases()`.

Added `model.SourceRepo = "repo"` constant alongside existing
`SourceRelease` and `SourceImport`.

## Files Changed

- `cmd/listreleases.go` — new `loadReleasesFromRepo()`, `metaToRecord()`,
  `sortRecordsByDate()`, `markLatestRecord()` helpers.
- `model/release.go` — added `SourceRepo` constant.
- `spec/01-app/21-list-releases.md` — updated data source section.
- `helptext/list-releases.md` — updated usage and examples.

## Status

Fixed in v2.33.0.
