# Refactor: release/assets.go

## Problem
`assets.go` is 261 lines with two responsibilities: public cross-compilation orchestration (types, detection, compile loop, staging) and low-level build helpers (single-target compilation, output naming, environment setup, file existence checks).

## Target Layout

### assets.go (~147 lines) — Orchestration & Public API
Stays:
- `type BuildTarget`
- `type CrossCompileResult`
- `DetectGoProject()`
- `ReadModuleName()`
- `BinaryName()`
- `FindMainPackages()`
- `CrossCompile()`
- `resolveBinName()`
- `CollectSuccessfulBuilds()`
- `EnsureStagingDir()`
- `CleanupStagingDir()`

### assetsbuild.go (~90 lines) — Build Helpers
Moves:
- `buildSingleTarget()`
- `formatOutputName()`
- `buildEnv()`
- `setEnv()`
- `fileExists()`

Imports: `fmt`, `os`, `os/exec`, `path/filepath`, `strings`

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
- [78-refactor-compress.md](78-refactor-compress.md) — zip and tar.gz compression

**Related `cmd/` refactors:**
- [70-refactor-listreleases.md](70-refactor-listreleases.md) — release listing
- [71-refactor-listversions.md](71-refactor-listversions.md) — version listing
