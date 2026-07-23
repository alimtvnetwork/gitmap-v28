# Refactor: release/compress.go

## Problem
`compress.go` is 204 lines with two responsibilities: zip-based compression (CompressAssets orchestration, zip creation, verbose logging, dry-run descriptions) and tar.gz-based compression (tar.gz creation, tar entry writing).

## Target Layout

### compress.go (~135 lines) — Orchestration & Zip
Stays:
- `CompressAssets()`
- `logCompressedArchive()`
- `compressSingle()`
- `isWindowsBinary()`
- `createZip()`
- `addFileToZip()`
- `DescribeCompression()`

### compresstar.go (~72 lines) — Tar.gz
Moves:
- `createTarGz()`
- `addFileToTar()`

Imports: `archive/tar`, `compress/gzip`, `fmt`, `io`, `os`, `path/filepath`

## Migration Rules
- No behaviour changes, no signature renames.
- Package remains `release`.
- Deduplicate imports per file.
- Blank line before every `return`.

## Acceptance Criteria
- Both files ≤ 200 lines.
- `go build ./...` succeeds.
- All existing tests pass unchanged.

---

## See Also

**Same package (`release/`) refactors:**

- [58-refactor-workflowfinalize.md](58-refactor-workflowfinalize.md) — pipeline, metadata, zip, GitHub upload
- [91-refactor-ziparchive.md](91-refactor-ziparchive.md) — zip I/O, dry-run, archive building
- [61-refactor-autocommit.md](61-refactor-autocommit.md) — auto-commit, git operations
- [63-refactor-workflowbranch.md](63-refactor-workflowbranch.md) — branch workflow, pending releases
- [64-refactor-workflow.md](64-refactor-workflow.md) — main workflow, validation
- [65-refactor-assets.md](65-refactor-assets.md) — cross-compilation, build helpers

**Related `cmd/` refactors:**
- [70-refactor-listreleases.md](70-refactor-listreleases.md) — release listing
- [71-refactor-listversions.md](71-refactor-listversions.md) — version listing
