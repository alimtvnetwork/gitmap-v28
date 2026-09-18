# Subtask 148.5: Archive Types Extraction & Result Aliases

> **Parent Plan:** [148-types-go-extraction-and-generic-result-centralization.md](../../pending/148-types-go-extraction-and-generic-result-centralization.md)
> **Target Subsystem:** `cli/archive/`
> **Status:** COMPLETED

## Scope of Work

1. Create `cli/archive/types.go`:
   - Move `FormatType`, `Format`, format constants from `archive.go`.
   - Move `CompressionModeType`, `CompressionMode`, mode constants from `create.go`.
   - Move `CreateOptions`, `CreateResult`, `ArchiveWriteParams` from `create.go`.
   - Move `ExtractResult`, `CompactExtractParams`, `ArchiveExtractParams`, `CopyDirEntryParams` from `extract.go`.
   - Move `Entry`, `ListExtractParams` from `list.go`.
   - Move `SourceKindType`, `SourceKind`, source constants, `ResolvedSource`, `Aria2cDownloadParams` from `source.go`.
   - Declare canonical alias: `type BoolResult = result.Result[bool]`.
2. Refactor function signatures in `cli/archive/`:
   - `isEntryIncluded(...) BoolResult` in `create.go`.
   - `matchPattern(...) BoolResult` in `create.go`.
   - `isHTTPURL(...) BoolResult` in `source.go`.
3. Verify callers and tests compile cleanly with `go test ./cli/archive/...`.
