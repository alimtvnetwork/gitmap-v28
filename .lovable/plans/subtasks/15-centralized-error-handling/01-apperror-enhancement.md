# Subtask 01: AppError Enhancement

master-plan: 15-centralized-error-handling-and-exit-architecture
subtask: 01-apperror-enhancement
status: pending

## Goal

Enhance `gitmap/apperror/apperror.go` to provide comprehensive error taxonomy, creator attribution, and metadata context.

## Specifications

1. Define `ErrorType` enum:
   - `ErrorTypeValidation`
   - `ErrorTypePrecondition`
   - `ErrorTypeNotFound`
   - `ErrorTypeExecution`
   - `ErrorTypeAbort`
   - `ErrorTypeInternal`
2. Define `SeverityType` enum:
   - `SeverityInfo`
   - `SeverityWarn`
   - `SeverityError`
   - `SeverityFatal`
3. Expand `AppError` struct:
   - `Op string` (Operation label)
   - `Code string` (Error code, e.g. E9001)
   - `Type ErrorType`
   - `Severity SeverityType`
   - `Creator string` (Component/package that created the error)
   - `Message string` (Human-readable explanation)
   - `Ctx map[string]any` (Diagnostic parameters)
   - `Cause error` (Underlying wrapped error)
4. Add constructors adhering to the 15-line limit per function.
5. Provide backward-compatible wrappers (`New`, `NewSimple`, `Wrap`, `WrapSimple`).

## Verification

- Unit tests in `gitmap/apperror/apperror_test.go` verifying formatting, unwrapping, and field access.
