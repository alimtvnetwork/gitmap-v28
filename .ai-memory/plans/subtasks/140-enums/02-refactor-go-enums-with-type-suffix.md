# Subtask 02: Refactor Go Enums to Enforce *Type Suffix

Parent Plan: [140-constants-and-enums-architecture.md](../../pending/140-constants-and-enums-architecture.md)

## Goals
1. Refactor `cli/cmd/commitin/enums.go`:
   - `ConflictMode` -> `ConflictModeType` + alias
   - `InputKind` -> `InputKindType` + alias
   - `RunStatus` -> `RunStatusType` + alias
   - `CommitOutcome` -> `CommitOutcomeType` + alias
   - `SkipReason` -> `SkipReasonType` + alias
   - `ExclusionKind` -> `ExclusionKindType` + alias
   - `MessageRuleKind` -> `MessageRuleKindType` + alias
   - `FunctionIntelLanguage` -> `FunctionIntelLanguageType` + alias
2. Refactor `cli/cmd/commitin/workspace/source.go`:
   - `SourceKind` -> `SourceKindType` + alias
3. Refactor `cli/dbengine_old/operators.go`:
   - `SqlOperator` -> `SqlOperatorType` + alias
4. Refactor `scripts/changelog/internal/runner/args.go`:
   - `Mode` -> `ModeType` + alias

## Acceptance Criteria
- [x] All 11 enums end with `Type` suffix.
- [x] Clean compilation with zero build breaks (`go vet ./...` exits 0).
- [x] Subtask completed and logged.
