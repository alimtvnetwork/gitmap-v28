# Subtask 148.3: Movemerge Types Centralization & Result Aliases

> **Parent Plan:** [148-types-go-extraction-and-generic-result-centralization.md](../../pending/148-types-go-extraction-and-generic-result-centralization.md)
> **Target Subsystem:** `cli/movemerge/`
> **Status:** COMPLETED

## Scope of Work

1. Expand `cli/movemerge/types.go`:
   - Move `DiffKindType` and enums (`DiffMissingLeft`, etc.) from `diff.go`.
   - Move `DiffEntry` struct from `diff.go`.
   - Move `FileMeta` struct from `walk.go`.
   - Move `ChoiceType` and enums from `conflict.go`.
   - Move `Resolver` struct from `conflict.go`.
   - Declare canonical aliases:
     - `type DiffEntryResult = result.Result[DiffEntry]`
     - `type FileMetaMapResult = result.ResultMap[string, FileMeta]`
2. Refactor function signatures:
   - `classifyOne(...) DiffEntryResult` in `diff.go`.
   - `classifyBoth(...) DiffEntryResult` in `diff.go`.
   - `IndexTree(root string, opts Options) FileMetaMapResult` in `walk.go`.
3. Verify tests compile cleanly with `go test ./cli/movemerge/...`.
