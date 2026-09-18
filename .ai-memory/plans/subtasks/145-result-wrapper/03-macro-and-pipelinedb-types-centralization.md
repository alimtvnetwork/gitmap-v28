# Subtask 03: Macro & PipelineDB Types Centralization

Parent Plan: [145-result-wrapper-and-types-go-centralization-audit.md](../../completed/145-result-wrapper-and-types-go-centralization-audit.md)

## Goals
1. Update `cli/macro/types.go`:
   - Define single reusable Result aliases:
     `type MacroSliceResult = result.ResultSlice[Macro]`
     `type MacroResult = result.Result[Macro]`
     `type MacroStepsMapResult = result.ResultMap[string, []MacroStep]`
2. Update `cli/macro/import.go`:
   - `ParseImportJSON`, `ParseImportYAML`, `ParseImportZIP`, `extractMacrosFromZip`, `ParseImportFile`, `filterAndValidateImportTargets` to return `MacroSliceResult`.
3. Update `cli/macro/import_sqlite.go`:
   - `ParseImportSQLite`, `queryAllMacrosWithSteps`, `scanMacroRows` to return `MacroSliceResult`.
   - `queryAllMacroSteps`, `scanMacroStepsMap` to return `MacroStepsMapResult`.
4. Update `cli/macro/storage.go`:
   - `ListMacros` to return `MacroSliceResult`.
5. Create `cli/pipelinedb/types.go`:
   - Define single reusable Result aliases:
     `type PipelineRunSliceResult = result.ResultSlice[PipelineRunRecord]`
     `type PipelineErrorSliceResult = result.ResultSlice[PipelineErrorRecord]`
     `type PipelineCompactErrorSliceResult = result.ResultSlice[PipelineCompactErrorRecord]`
     `type PipelineRunIdSliceResult = result.ResultSlice[uint64]`
     `type PipelineRunIdMapResult = result.ResultMap[uint64, bool]`
     `type PipelineDbStatsResult = result.Result[PipelineDbStats]`
6. Update `cli/pipelinedb/pipeline_split_ops.go`:
   - Update 14 functions returning raw generic ResultSlice/ResultMap to use canonical `Pipeline*Result` types from `types.go`.
7. Run `go test -v ./macro/... ./pipelinedb/...` in `cli/` to verify 100% pass.

## Acceptance Criteria
- [x] `cli/macro/types.go` and `cli/pipelinedb/types.go` define all single reusable Result aliases.
- [x] 24 functions refactored to return canonical single type aliases from `types.go`.
- [x] Unit tests pass with 0 errors.
