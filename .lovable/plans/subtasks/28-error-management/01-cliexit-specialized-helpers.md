# Subtask 01 - Cliexit Specialized Helpers & ExitCodeType Enum

## Parent Specification
[28-error-management-audit.md](.lovable/plans/pending/28-error-management-audit.md)

## Acceptance Criteria & Requirements
- Define `ExitCodeType` enum in `gitmap/cliexit/handle.go`:
  - `ExitCodeSuccess ExitCodeType = 0`
  - `ExitCodeGeneralError ExitCodeType = 1`
  - `ExitCodeUsageError ExitCodeType = 2`
  - `ExitCodePartialFailure ExitCodeType = 3`
  - `ExitCodeNotFound ExitCodeType = 4`
  - `ExitCodeValidationError ExitCodeType = 5`
- Add specialized exit helpers:
  - `HandleValidationError(err error)`
  - `HandleUsageError(err error)`
  - `HandleGeneralError(err error)`
  - `HandleSuccess()`
- Maintain strict function length <= 15 lines, positive naming, and blank lines before returns.
- Verify compilation via `go build -C gitmap .`.
