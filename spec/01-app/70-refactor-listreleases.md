# Refactor: cmd/listreleases.go

## Problem
`listreleases.go` is 229 lines with two responsibilities: command orchestration with flag parsing/filtering and data loading with format conversion (repo meta files, DB fallback, record sorting).

## Target Layout

### listreleases.go (~140 lines) — Command & Display
Stays:
- `runListReleases()`
- `parseListReleasesSource()`
- `filterBySource()`
- `hasListReleasesJSONFlag()`
- `parseListReleasesLimit()`
- `applyReleaseLimit()`
- `loadReleases()`
- `printReleasesTerminal()`
- `printReleaseRow()`
- `printReleasesJSON()`

### listreleasesload.go (~100 lines) — Data Loading
Moves:
- `loadReleasesFromRepo()`
- `convertMetasToRecords()`
- `metaToRecord()`
- `sortRecordsByDate()`
- `markLatestRecord()`
- `loadReleasesFromDB()`

Imports: `fmt`, `os`, `sort`, `strings`, `constants`, `model`, `release`

## Migration Rules
- No behaviour changes, no signature renames.
- Package remains `cmd`.
- Deduplicate imports per file.
- Blank line before every `return`.

## Acceptance Criteria
- Both files ≤ 200 lines.
- `go build ./...` succeeds.
- All existing tests pass unchanged.

---

## See Also

**Same package (`release/cmd/`) refactors:**

- [71-refactor-listversions.md](71-refactor-listversions.md) — version listing

**Related `release/` refactors:**
- [58-refactor-workflowfinalize.md](58-refactor-workflowfinalize.md) — pipeline, metadata, zip, GitHub upload
- [63-refactor-workflowbranch.md](63-refactor-workflowbranch.md) — branch workflow, pending releases
