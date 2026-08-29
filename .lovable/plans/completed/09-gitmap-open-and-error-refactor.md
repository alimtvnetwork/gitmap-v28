# Gitmap Open Command & Global Error Management Refactor

## Parent Task Goal

Implement `gitmap open <file|folder|url|mail>` for cross-platform execution and refactor the entire `cmd` package to return typed `*apperror.AppError` values from command handlers (`runXxx` functions), preventing abrupt `os.Exit` calls and allowing centralized error logging.

## Core Requirements

1. **Error Management**: 
   - Refactor `dispatchEntry` in `rootdispatch.go` to use `handler func() *apperror.AppError` (or `error`).
   - Refactor all `runXxx` functions across `gitmap/cmd/*.go` to return `error` or `*apperror.AppError` instead of `void`.
   - Update `runDispatchTable` to capture errors, wrap them with `apperror.Wrap`, log them using `cliexit`, and exit gracefully.
   - Implement the `ErrorLog` database table if requested, or integrate with existing logs.
2. **Gitmap Open**:
   - Implement `gitmap open <file|folder|url|mail>` handling OS-specific openers (`start`, `open`, `xdg-open`).
   - Add help text and CLI bindings.

## Subtasks Execution Strategy

1. **Subtask 1: Central Dispatch Refactor**
   - Update `rootdispatch.go` and `root.go` to support error-returning dispatchers.
   - Modify the `runDispatch` pipeline to handle top-level errors and log them.

2. **Subtask 2: AST Refactor Script for Handlers**
   - Create a Go AST/Regex script in `.lovable/temp-scripts/refactor_cmds.go`.
   - Execute the script to automatically rewrite all `runXxx` function signatures and `dispatchEntry` registrations.
   - Fix compilation errors.

3. **Subtask 3: Implement `gitmap open`**
   - Write `gitmap/cmd/open.go` containing `runOpen(args []string) *apperror.AppError`.
   - Add cross-platform OS openers.
   - Create helptext and CLI mappings.

4. **Subtask 4: Validation & Release**
   - Ensure all `go test ./...` pass.
   - Verify `gitmap open` works correctly.
   - Update README and What-to-Read.
   - Group commits and push.
