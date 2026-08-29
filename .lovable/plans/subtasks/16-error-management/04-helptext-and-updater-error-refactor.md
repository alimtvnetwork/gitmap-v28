# Subtask 04: Helptext & Gitmap-Updater Error Handling Refactoring

> **Parent Plan:** `16-error-management-audit.md`  
> **Files:** `gitmap/helptext/print.go`, `gitmap-updater/cmd/*.go`

## Objective

Refactor bare `os.Exit` calls in help text and the updater subsystem into structured error handling and centralized exit calls.

## Action Steps

1. [x] Refactor `gitmap/helptext/print.go` to use `cliexit.HandleError` or standard error return.
2. [x] Refactor `gitmap-updater/cmd/*.go` to construct structured error reports before exiting.

## Status: COMPLETED

All helptext, updater, and cmd error handling refactored cleanly to use structured AppError and cliexit dispatch.
