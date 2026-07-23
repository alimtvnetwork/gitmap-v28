# Refactor: cmd/listversions.go

## Problem
`listversions.go` is 226 lines with two responsibilities: command orchestration with flag parsing/filtering and data collection with output formatting (git tag parsing, changelog loading, terminal/JSON rendering).

## Target Layout

### listversions.go (~130 lines) — Command & Filtering
Stays:
- `type versionEntry`
- `runListVersions()`
- `parseListVersionsSource()`
- `filterVersionsBySource()`
- `hasListVersionsJSONFlag()`
- `parseListVersionsLimit()`
- `applyVersionLimit()`
- `collectVersionEntries()`
- `loadVersionSourceMap()`

### listversionsutil.go (~105 lines) — Data Collection & Output
Moves:
- `collectVersionTags()`
- `parseVersionTags()`
- `loadChangelogMap()`
- `printVersionEntriesTerminal()`
- `type lvJSONEntry`
- `printVersionEntriesJSON()`

Imports: `encoding/json`, `fmt`, `os`, `os/exec`, `sort`, `strings`, `constants`, `release`

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

- [70-refactor-listreleases.md](70-refactor-listreleases.md) — release listing

**Related `release/` refactors:**
- [58-refactor-workflowfinalize.md](58-refactor-workflowfinalize.md) — pipeline, metadata, zip, GitHub upload
- [63-refactor-workflowbranch.md](63-refactor-workflowbranch.md) — branch workflow, pending releases
