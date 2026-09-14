---
name: apperror-stacktrace-skip-and-filter
description: Autonomously implement and verify AppError stack trace skip configuration, default skip constants, dynamic skip increments, and internal frame filtering to eliminate internal error creation frames from stack traces across GitMap.
---

# AppError Stack Trace Skip, Default Configuration and Defensive Filtering

## Overview
Provides guidelines and invariants for capturing stack traces and caller locations in `cli/apperror/`.
Ensures that internal `apperror` constructor frames (e.g. `NewSimple`, `New`, `Wrap`, etc.) are completely omitted from stack traces, provides a configurable default skip count, supports runtime/per-call skip elevation to pierce wrapper layers, and enforces defensive frame filtering.

## Core Rules & Invariants
1. **Default Skip Constants**:
   - `DefaultStackTraceSkip = 3`: By default skips `Callers`, `captureStackTrace`, and constructor frame to immediately record the caller's call site.
   - `DefaultCallerSkip = 2`: Skips `captureCaller` and constructor frame to isolate caller.
2. **Dynamic Skip Elevation**:
   - Support `WithSkip(additional int)` and `WithAdditionalSkip(additional int)` on `*AppError` to recalculate frames when intermediate wrappers are introduced.
   - Support `SetDefaultStackTraceSkip(skip int)` and `SetDefaultCallerSkip(skip int)` for global reconfiguration.
   - Provide `CaptureStackTrace(skip int) string` and `CaptureCaller(skip int) string`.
   - Provide `NewWithSkip(skip int, op string, code string) *AppError` and `WrapWithSkip(skip int, err error, op string) *AppError`.
3. **Defensive Frame Filtering**:
   - In `appendStackFrame`, verify `isAppErrorInternal(frame)`.
   - Any stack frame matching internal `apperror.go` or `cli/apperror.` constructors must be suppressed from stack traces.
4. **Coding Guidelines Invariants**:
   - All functions <= 15 lines (target <= 8 lines).
   - Affirmative booleans only (`is*`, `has*`).
   - Zero nested ifs (guard clauses and early returns).
   - Strict Unix LF line endings.
