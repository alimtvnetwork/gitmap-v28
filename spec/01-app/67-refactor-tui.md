# Refactor: tui/tui.go

## Problem
`tui.go` is 245 lines with two responsibilities: model definition, initialization, and update logic; and view rendering (tabs, content dispatch, status bar).

## Target Layout

### tui.go (~157 lines) — Model & Update
Stays:
- View index constants
- `type rootModel`
- `Run()`
- `newRootModel()`
- `Init()`
- `Update()`
- `updateActiveView()`

### tuiview.go (~97 lines) — View Rendering
Moves:
- `View()`
- `renderTabs()`
- `renderContent()`
- `renderStatusBar()`

Imports: `strings`, `constants`

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

- [77-refactor-logs.md](77-refactor-logs.md) — log model, search, view

**Related `cmd/` refactors:**
- [75-refactor-status.md](75-refactor-status.md) — status display
- [76-refactor-exec.md](76-refactor-exec.md) — batch execution
