# Subtask 05: Enhanced Error Diagnostics & Operational Stack Traces

## Scope
- Update `gitmap/cliexit/handle.go`:
  - Fix bug on line 73: change `if e.Cause != nil && e.Message != ""` to `if e.Cause != nil`.
  - Render origin caller (`e.Caller`), command context, exit code, log path, and stderr snippet.
  - Print diagnostic operational stack traces for fatal execution errors (`E9000:EXECUTION`) or when running with `--debug` / `GITMAP_DEBUG=1`.
- Update `gitmap/cmd/installtools.go`:
  - Refactor `handleInstallError` to wrap a rich `AppError` with exit code, command, and stderr rather than returning bare `apperror.NewSimple("fatal error", "E9000")`.

## Files Touched
- `gitmap/cliexit/handle.go`
- `gitmap/cmd/installtools.go`
