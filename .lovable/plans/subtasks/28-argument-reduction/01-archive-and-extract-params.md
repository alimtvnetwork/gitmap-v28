# Subtask 01: Archive and Extraction Parameter Structs

- **Status:** `PENDING`
- **Target Files:** `gitmap/archive/create.go`, `gitmap/archive/extract.go`, `gitmap/archive/list.go`

## Instructions

1. Refactor `matchAny(val string, patterns []string, emptyDefault bool)` in `gitmap/archive/create.go` to rename `emptyDefault` to `isEmptyDefault`.
2. Encapsulate multi-argument functions in `gitmap/archive/extract.go` (`completeCompactExtract`, `runArchiveExtraction`, `copyDirEntry`) into value-based parameter structs.
3. Enforce `is`/`has` boolean prefixing on all new struct fields.
4. Ensure function lengths stay <= 15 lines and conditionals are flat (depth 0).
