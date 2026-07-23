# Refactor: tui/logs.go

## Problem
`logs.go` is 215 lines with two responsibilities: model state management (struct, initialization, update, search, filtering, key handling) and view rendering (list view, detail view, row formatting).

## Target Layout

### logs.go (~137 lines) — Model & Update
Stays:
- `type logsModel`
- `newLogsModel()`
- `Update()`
- `updateSearch()`
- `applyFilter()`
- `matchesLogQuery()`
- `handleKey()`

### logsview.go (~88 lines) — View Rendering
Moves:
- `View()`
- `viewList()`
- `viewDetail()`

Imports: `fmt`, `strings`, `constants`

## Migration Rules
- No behaviour changes, no signature renames.
- Package remains `tui`.
- Deduplicate imports per file.
- Blank line before every `return`.

## Acceptance Criteria
- Both files ≤ 200 lines.
- `go build ./...` succeeds.
- All existing tests pass unchanged.

---

## See Also

**Same package (`tui/`) refactors:**

- [67-refactor-tui.md](67-refactor-tui.md) — model, update, view rendering

**Related `cmd/` refactors:**
- [75-refactor-status.md](75-refactor-status.md) — status display
- [76-refactor-exec.md](76-refactor-exec.md) — batch execution
