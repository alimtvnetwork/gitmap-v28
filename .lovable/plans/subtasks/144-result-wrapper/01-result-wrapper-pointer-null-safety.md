# Subtask 01: Result Wrapper Pointer Null Safety & Core Predicates

Parent Plan: [144-result-wrapper-null-safety-and-single-return-audit.md](../../completed/144-result-wrapper-null-safety-and-single-return-audit.md)

## Goals
1. Refactor `cli/result/result.go`:
   - Attach all methods to pointer receiver `(r *Result[T])`.
   - Add explicit `if r == nil` guards returning canonical defaults on line 1 of every method.
   - Implement missing 4 core predicates: `Count() int`, `IsCountOtherThan(number int) bool`, `HasRecord() bool`, `HasRecords() bool`, `IsDefined() bool`.
2. Refactor `cli/result/result_slice.go`:
   - Attach all methods to pointer receiver `(r *ResultSlice[T])`.
   - Add explicit `if r == nil` guards returning canonical defaults on line 1 of every method.
   - Implement missing 4 core predicates: `IsCountOtherThan(number int) bool`, `HasRecord() bool`, `HasRecords() bool`, `IsDefined() bool`, `Items() []T`.
3. Refactor `cli/result/result_map.go`:
   - Attach all methods to pointer receiver `(r *ResultMap[K, V])`.
   - Add explicit `if r == nil` guards returning canonical defaults on line 1 of every method.
   - Implement missing 4 core predicates: `IsCountOtherThan(number int) bool`, `HasRecord() bool`, `HasRecords() bool`, `IsDefined() bool`.
4. Author comprehensive test suite in `cli/result/null_safety_test.go`:
   - Verify every method called on nil pointers `(*Result[T])(nil)`, `(*ResultSlice[T])(nil)`, and `(*ResultMap[K, V])(nil)` executes without panicking and returns safe defaults.
   - Verify non-nil values behave identically to existing contract.

## Acceptance Criteria
- [x] All methods on `Result[T]`, `ResultSlice[T]`, and `ResultMap[K, V]` have pointer receivers with `if r == nil` guards.
- [x] `IsCountOtherThan`, `IsEmpty`, `HasRecord`, `IsDefined` implemented on all three types.
- [x] `go test -v ./cli/result/...` passes with 100% green and zero panics on nil receivers.
- [x] `go vet ./...` in `cli/` passes with 0 errors.
