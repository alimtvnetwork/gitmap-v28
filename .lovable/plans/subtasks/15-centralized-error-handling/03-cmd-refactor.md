# Subtask 03: Command Handlers Centralized Error Refactor

master-plan: 15-centralized-error-handling-and-exit-architecture
subtask: 03-cmd-refactor
status: pending

## Goal
Refactor all error exit points across `gitmap/cmd/*.go` to construct structured `apperror.AppError` instances and dispatch them via `cliexit.HandleError`.

## Scope
1. `gitmap/cmd/root.go`:
   - Zero arguments: returns cleanly (exits 0 with usage).
   - Unknown command: raises `apperror.NewWithDetails("cmd.dispatch", "E1001", "unknown command", ...)` and calls `cliexit.HandleError(err, 1)`.
2. `gitmap/cmd/releasepull.go`:
   - Not in repo: raises `apperror.NewWithDetails("release.pull", "E2001", ...)` and dispatches.
3. `gitmap/cmd/reinstall.go`:
   - Reinstall aborted: raises `apperror.NewWithDetails("cmd.reinstall", "E2002", ...)` and dispatches.
   - No repo linked: raises `apperror.NewWithDetails("cmd.reinstall", "E2003", ...)` and dispatches.
4. `gitmap/cmd/rootadd.go`:
   - Insufficient args: raises `apperror.NewWithDetails("cmd.add", "E2004", ...)` and dispatches.
   - Unknown subcommand: raises `apperror.NewWithDetails("cmd.add", "E2005", ...)` and dispatches.
5. All other command sites (`releaserebase.go`, `reset.go`, `revert.go`, `revertscript.go`, `reverttxn.go`, `reverttxn_lastn.go`, `scanresolve.go`, `selfuninstallhandoff.go`, `seowrite.go`, `seowritecsv.go`, `seowriteloop.go`, `seowritetemplate.go`, `sshcat.go`).

## Verification
- Code builds cleanly with `go vet ./...` and `go test -run=^$ ./...`.
