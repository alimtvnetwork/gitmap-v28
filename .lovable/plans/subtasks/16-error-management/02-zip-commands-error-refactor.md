# Subtask 02: Zip & ZipGroup Commands Error Handling Refactoring

> **Parent Plan:** `16-error-management-audit.md`
> **Files:** `gitmap/cmd/zipgroup.go`, `gitmap/cmd/zipgroupcreate.go`, `gitmap/cmd/zipgroupops.go`, `gitmap/cmd/zipgroupshow.go`, `gitmap/cmd/zip.go`

## Objective

Replace all `panic("error")` and bare `os.Exit` calls in zip-related commands with `apperror` and `cliexit.HandleError`.

## Action Steps

1. Inspect `zipgroup.go`, `zipgroupcreate.go`, `zipgroupops.go`, `zipgroupshow.go`, and `zip.go`.
2. Construct structured `apperror.NewWithDetails` / `apperror.WrapWithDetails` with full operation names (`cmd.zip.group`, `cmd.zip.create`, `cmd.zip.ops`, etc.).
3. Pass to `cliexit.HandleError(appErr, 1)`.
