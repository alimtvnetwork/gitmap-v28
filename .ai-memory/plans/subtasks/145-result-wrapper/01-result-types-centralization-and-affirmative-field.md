# Subtask 01: Result Package Types Centralization & Affirmative Field Naming

Parent Plan: [145-result-wrapper-and-types-go-centralization-audit.md](../../completed/145-result-wrapper-and-types-go-centralization-audit.md)

## Goals
1. Create `cli/result/types.go`:
   - Declare `Result[T any] struct { ... }`
   - Declare `ResultSlice[T any] struct { ... }`
   - Declare `ResultMap[K comparable, V any] struct { ... }`
   - Declare single reusable aliases:
     `Wrap[T any] = Result[T]`
     `Slice[T any] = ResultSlice[T]`
     `Map[K comparable, V any] = ResultMap[K, V]`
2. Remove struct definitions from `result.go`, `result_slice.go`, and `result_map.go` to prevent duplicate declarations.
3. Fix affirmative boolean field naming in `Result[T]`:
   - Rename `defined bool` to `isDefined bool`.
   - Update `cli/result/result.go:115` (`r.defined` -> `r.isDefined`).
   - Update `cli/result/result.go:202` (`defined: true` -> `isDefined: true`).
4. Modernize clumsy checks in `cli/result/null_safety_test.go`:
   - Line 292: `sliceOk.IsFailure() || sliceOk.Count() != 2` -> `sliceOk.IsCountOtherThan(2)`
   - Line 304: `sliceEmpty.IsFailure() || sliceEmpty.Count() != 0` -> `sliceEmpty.IsCountOtherThan(0)`
   - Line 313: `mapOk.IsFailure() || mapOk.Count() != 1` -> `mapOk.IsCountOtherThan(1)`
5. Run `go test -v ./result/...` in `cli/` to verify 100% pass.

## Acceptance Criteria
- [x] `cli/result/types.go` created with core container types and aliases.
- [x] Affirmative `isDefined bool` enforced.
- [x] Clumsy caller checks replaced with `IsCountOtherThan(N)`.
- [x] Unit tests pass with 0 errors.
