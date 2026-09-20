# RCA 64: Undefined apperror.NewNotFound in Locate Engine

## 1. Symptom
GitHub Actions CI runs (`#35510423912`, `#35510417478`, `#35510417413`, `#35510417405`) on commit `dda4361a` failed during compilation:
```text
cmdautomation/locate.go:41:24: undefined: apperror.NewNotFound
```
This compilation error failed all 6 remote CI/CD jobs (Release, macos build, ubuntu build, windows build, race detector, history rewrite smoke).

## 2. Root Cause
In `cli/cmdautomation/locate.go`, the tool search fallback returned:
```go
return res, apperror.NewNotFound("tool_locate", "E_TOOL_NOT_FOUND", target)
```
However, `cli/apperror/apperror.go` only defined `apperror.NewNotFoundError(msg string) *AppError`, and did not provide `apperror.NewNotFound(op, code, msg string) *AppError`.

## 3. Resolution
1. **Added `NewNotFound` constructor to `cli/apperror/apperror.go`:**
   - Implemented `NewNotFound(op, code, msg string) *AppError` setting `ErrorTypeNotFound`, `SeverityError`, `Caller`, and empty stack trace per coding guidelines.
2. **Added unit test `TestAppError_NewNotFound` to `cli/apperror/apperror_test.go`:**
   - Validates `Type == ErrorTypeNotFound`, `Op`, `Code`, `Message`, and empty `Stack`.
3. **Confirmed `locate.go` compliance:**
   - All functions remain <= 15 lines, affirmative booleans, and correct `*apperror.AppError` return contracts.

## 4. Verification
- CI/CD build failure cleanly addressed.
- Unit test added and verified.
