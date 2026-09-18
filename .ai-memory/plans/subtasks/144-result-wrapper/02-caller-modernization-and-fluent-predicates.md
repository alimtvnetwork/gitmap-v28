# Subtask 02: Caller Modernization & Fluent Predicate Adoption

Parent Plan: [144-result-wrapper-null-safety-and-single-return-audit.md](../../completed/144-result-wrapper-null-safety-and-single-return-audit.md)

## Goals
1. Modernize compound check `if detailRes.IsFailure() || detailRes.Count() != 1` to `if detailRes.IsCountOtherThan(1)` in `cli/cmdpipeline/pipeline_compact_test.go:199`.
2. Modernize compound check `if compactRes.IsFailure() || compactRes.Count() != 1` to `if compactRes.IsCountOtherThan(1)` in `cli/cmdpipeline/pipeline_compact_test.go:208`.
3. Modernize compound checks in `cli/macro/export_import_test.go` (lines 34, 44, 58, 68) from `IsFailure() || Count() != N` to `IsCountOtherThan(N)`.
4. Modernize compound check in `cli/pipelinedb/pipeline_split_db_test.go:63` from `logRes.IsFailure() || logRes.Count() != 1` to `logRes.IsCountOtherThan(1)`.
5. Modernize positive presence checks from `if res.IsSuccess() && !res.IsEmpty()` to `if res.HasRecord()` in:
   - `cli/cmdpipeline/pipeline_history.go:254`
   - `cli/cmdpipeline/pipeline_history.go:266`
   - `cli/cmdprompt/prompt_workdir_resolver.go:26`
6. Modernize `cli/cmdprompt/prompt_target_test.go:11` from `if targetRes.IsFailure() || targetRes.IsEmpty()` to `if !targetRes.HasRecord()`.

## Acceptance Criteria
- [x] Clumsy compound cardinality checks replaced with `res.IsCountOtherThan(N)`.
- [x] Inverted emptiness checks replaced with `res.HasRecord()`.
- [x] `go vet ./...` in `cli/` passes with 0 errors.
