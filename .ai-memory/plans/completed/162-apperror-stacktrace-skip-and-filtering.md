# Plan 162: AppError Stack Trace Skip, Default Configuration & Defensive Filtering Suite

## Overview
Autonomously implement stack trace caller skip configuration, default skip constants, dynamic skip increments, and internal frame filtering to eliminate internal error creation frames (e.g. `at .../cli/apperror.NewSimple (apperror/apperror.go:159)`) from stack traces everywhere across GitMap.

## Key Goals
1. **Configurable Default Skip Constants**:
   - `DefaultStackTraceSkip = 3`: By default skips `runtime.Callers`, `captureStackTrace`, and constructor frame (e.g. `NewSimple`, `New`, `Wrap`, etc.) to immediately record the caller's call site (e.g. `cmd.registerGHDesktop`).
   - `DefaultCallerSkip = 2`: Skips `captureCaller` and constructor frame to isolate caller.
   - Global setters: `SetDefaultStackTraceSkip(skip int)` and `SetDefaultCallerSkip(skip int)`.
2. **Defensive Frame Filtering**:
   - In `appendStackFrame`, verify `isAppErrorInternal(frame)`.
   - Filter out any stack frames matching `apperror.go` or `cli/apperror.` internal constructor frames across all stack trace rendering.
   - Ensures internal helper lines never appear in user-facing stack traces even if skip numbers are altered or functions inlined.
3. **Dynamic Skip Elevation**:
   - Support `err.WithSkip(additionalSkip int) *AppError` and `err.WithAdditionalSkip(additionalSkip int) *AppError` to pierce wrapper layers when intermediate helper functions wrap errors.
   - Export `CaptureStackTrace(skip int) string` and `CaptureCaller(skip int) string` for custom callers.
   - Provide `NewWithSkip(skip int, op, code string) *AppError` and `WrapWithSkip(skip int, err error, op, code string) *AppError`.
4. **Constructors Modernization**:
   - Update all constructors in `cli/apperror/apperror.go`: `New`, `NewSimple`, `NewWithDetails`, `NewValidationError`, `NewExecutionError`, `Wrap`, `WrapSimple`, `WrapWithDetails` to use `captureStackTrace(DefaultStackTraceSkip)` and `captureCaller(DefaultCallerSkip)`.
5. **Coding Guidelines Invariants**:
   - Functions strictly <= 15 lines (target <= 8 lines).
   - Affirmative booleans only (`is*`, `has*`). Zero nested ifs.
   - Strict Unix LF line endings.

## Custom Rules
1. Functions strictly <= 15 lines (target <= 8 lines).
2. Affirmative booleans only (isValid, isInternal). No negatives.
3. Universal AppError wrapping on all errors.
4. Absolute ban on `go test` and `go build` during routine turns.

> **Task Origin**: User requested skipping `apperror.NewSimple (apperror/apperror.go:159)` stack trace line everywhere with configurable default skip number that can be increased to fix issues.
> **Total Loops**: 1 continuous loop.

## Consolidated Subtasks Detail

# Subtask 162-01: AppError Stack Trace Skip, Default Values, Dynamic Elevation & Defensive Filtering

## Target Files
- cli/apperror/apperror.go
- cli/apperror/apperror_test.go

## Implemented Features
1. Defined configurable default variables:
   ```go
   var (
       DefaultStackTraceSkip = 3
       DefaultCallerSkip     = 2
   )
   ```
2. Added global setters `SetDefaultStackTraceSkip(skip int)` and `SetDefaultCallerSkip(skip int)`.
3. Exported `CaptureStackTrace(skip int) string` and `CaptureCaller(skip int) string`.
4. Added `isAppErrorInternal(frame runtime.Frame) bool` and `isAppErrorFile(file string) bool`.
5. In `appendStackFrame`: unconditionally drops frames where `isAppErrorInternal(frame)` is true.
6. In `captureCaller`: checks callers up to 5 levels to safely skip any intermediate `apperror.go` frames.
7. Added `(e *AppError) WithSkip(additional int) *AppError` and `(e *AppError) WithAdditionalSkip(additional int) *AppError`.
8. Added `NewWithSkip(skip int, op string, code string) *AppError` and `WrapWithSkip(skip int, err error, op string, code string) *AppError`.
9. Updated all standard constructors (`New`, `NewSimple`, `NewWithDetails`, `NewValidationError`, `NewExecutionError`, `Wrap`, `WrapSimple`, `WrapWithDetails`) to use `captureCaller(DefaultCallerSkip)` and `captureStackTrace(DefaultStackTraceSkip)`.
10. Added unit tests in `cli/apperror/apperror_test.go` asserting that `apperror.NewSimple` and `apperror.go:` lines are absent from stack traces, callers point to the test call site, `WithSkip` elevates skip frames, and global setters update defaults.
11. Passed all linters: `check-nested-ifs.py`, `check-enum-and-boolean.py`, `check-error-management.py`.
