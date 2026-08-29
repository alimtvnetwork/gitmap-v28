# Subtask 02: CLI Exit Central Handler Integration

master-plan: 15-centralized-error-handling-and-exit-architecture
subtask: 02-cliexit-central-handler
status: pending

## Goal
Enhance `gitmap/cliexit/cliexit.go` to provide centralized error handling, strategy selection (exit vs panic), pipe draining, and user-actionable reporting.

## Specifications
1. Implement `cliexit.HandleError(err error, defaultCode int)`:
   - Inspects if `err` is `*apperror.AppError`.
   - If not, wraps it into `apperror.WrapSimple(err, "unknown")`.
   - Extracts `Code`, `Op`, `Creator`, `Severity`, `Ctx`, and `Cause`.
   - Formats to `os.Stderr`:
     `gitmap [<Code>] <Op>: <Message> (creator: <Creator>, context: <Ctx>)`
   - Drains registered flushers (`runFlushers()`).
   - If environment variable `GITMAP_ERROR_PANIC=1` is set, invokes `panic(err)` for debug/assertion runs.
   - Otherwise, invokes `os.Exit(code)`.
2. Implement helper `cliexit.FailAppError(appErr *apperror.AppError, code int)` as a direct semantic dispatcher.
3. Keep all functions under 15 lines.

## Verification
- Unit tests in `gitmap/cliexit/cliexit_test.go` covering formatted reporting and panic mode.
