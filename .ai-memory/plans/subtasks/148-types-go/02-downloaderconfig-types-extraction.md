# Subtask 148.2: Downloaderconfig Types Extraction & Result Centralization

> **Parent Plan:** [148-types-go-extraction-and-generic-result-centralization.md](../../pending/148-types-go-extraction-and-generic-result-centralization.md)
> **Target Subsystem:** `cli/downloaderconfig/`
> **Status:** COMPLETED

## Scope of Work

1. Create `cli/downloaderconfig/types.go`:
   - Declare `ConfigKeyType string` and `type ConfigKey = ConfigKeyType`.
   - Declare key constants: `KeyDefaultSplitSize`, `KeyLargeFileSplitSize`, `KeyLargeFileThreshold`, `KeyTinyFileThreshold`, `KeyTinyFileSplitSize`.
   - Declare `Document` struct.
   - Declare `DownloaderConfig` struct with affirmative booleans (`AllowFallback`, `OverwriteUserConfig`).
   - Declare `DatabaseVersion` struct.
   - Declare canonical aliases:
     - `type DocumentResult = result.Result[Document]`
     - `type BytesResult = result.Result[[]byte]`
2. Remove extracted types and constants from `cli/downloaderconfig/downloaderconfig.go`.
3. Refactor function signatures in `cli/downloaderconfig/downloaderconfig.go`:
   - `LoadFile(path string) DocumentResult`
   - `Parse(raw []byte) DocumentResult`
   - `Marshal(doc Document) BytesResult`
4. Verify callers (`cli/cmd/downloaderconfig.go`, `cli/store/downloader_seed.go`, `cli/store/downloader_settings.go`) and tests compile cleanly with `go test ./cli/downloaderconfig/...`.
