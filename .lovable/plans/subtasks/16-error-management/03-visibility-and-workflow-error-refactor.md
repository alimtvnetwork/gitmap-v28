# Subtask 03: Visibility & Workflow Commands Error Handling Refactoring

> **Parent Plan:** `16-error-management-audit.md`  
> **Files:** `gitmap/cmd/visibilitymakelast.go`, `visibilityredo.go`, `visibilityresolve.go`, `visibilityundo.go`, `visibilityundoflags.go`, `workflow_open_pr.go`, `workflow_recent_todo.go`

## Objective

Refactor bare `os.Exit` calls to use `apperror.NewWithDetails` + `cliexit.HandleError` while preserving explicit visibility exit codes (`constants.ExitVisBadFlag`, `constants.ExitVisAuthFailed`, `constants.ExitVisNotARepo`, etc.).

## Action Steps

1. Inspect all visibility command handlers and workflow command handlers.
2. Route exits through `cliexit.HandleError(appErr, exitCode)`.
