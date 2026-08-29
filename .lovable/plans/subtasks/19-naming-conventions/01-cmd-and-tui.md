# Subtask 19.01: Refactor Bare `ok` in `cmd/` and `tui/`

## Target
- Files: `gitmap/cmd/*.go`, `gitmap/tui/*.go`

## Violations to Fix
- [ ] Replace `keyMsg, ok := msg.(tea.KeyMsg)` with `keyMsg, isKeyMsg := msg.(tea.KeyMsg)`.
- [ ] Replace `appErr, ok := err.(*apperror.AppError)` with `appErr, isAppErr := err.(*apperror.AppError)`.
- [ ] Replace map lookups `val, ok := map[key]` with `val, isFound := map[key]`.

## Acceptance Criteria
- [ ] Zero bare `ok` identifiers in `gitmap/cmd/` and `gitmap/tui/`.
- [ ] `go test ./gitmap/cmd/... ./gitmap/tui/...` passes cleanly.
