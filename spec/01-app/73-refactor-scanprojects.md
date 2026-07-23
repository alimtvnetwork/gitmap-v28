# Refactor: cmd/scanprojects.go

## Problem
`scanprojects.go` is 224 lines with two responsibilities: top-level project detection orchestration (detect, upsert, resolve IDs) and language-specific metadata persistence (Go runnables, C# project files, C# key files, stale cleanup).

## Target Layout

### scanprojects.go (~98 lines) — Orchestration
Stays:
- `detectAllProjects()`
- `upsertProjectsToDB()`
- `upsertProjectRecords()`
- `resolveDetectedProjectID()`
- `upsertProjectMetadata()`

### scanprojectsmeta.go (~130 lines) — Metadata Persistence
Moves:
- `upsertGoProjectMeta()`
- `upsertGoRunnables()`
- `upsertCSharpProjectMeta()`
- `upsertCSharpFiles()`
- `upsertCSharpKeyFiles()`
- `collectRepoIDs()`
- `cleanStaleProjects()`
- `collectKeepIDs()`

Imports: `fmt`, `os`, `constants`, `detector`, `model`, `store`

## Migration Rules
- No behaviour changes, no signature renames.
- Package remains `cmd`.
- Deduplicate imports per file.
- Blank line before every `return`.

## Acceptance Criteria
- Both files ≤ 200 lines.
- `go build ./...` succeeds.
- All existing tests pass unchanged.

---

## See Also

**Same package (`cmd/`) refactors:**

- [90-refactor-root-dispatch.md](90-refactor-root-dispatch.md) — dispatch splitting
- [62-refactor-seowriteloop.md](62-refactor-seowriteloop.md) — SEO write loop, git ops
- [66-refactor-zipgroupops.md](66-refactor-zipgroupops.md) — zip group CRUD and display
- [68-refactor-aliasops.md](68-refactor-aliasops.md) — alias CRUD and suggest
- [69-refactor-tempreleaseops.md](69-refactor-tempreleaseops.md) — temp release branch ops
- [72-refactor-sshgen.md](72-refactor-sshgen.md) — SSH key generation
- [74-refactor-amendexec.md](74-refactor-amendexec.md) — git amend operations
- [75-refactor-status.md](75-refactor-status.md) — status display
- [76-refactor-exec.md](76-refactor-exec.md) — batch execution

**Related `release/` refactors:**
- [58-refactor-workflowfinalize.md](58-refactor-workflowfinalize.md) — pipeline, metadata, zip, GitHub upload
- [64-refactor-workflow.md](64-refactor-workflow.md) — main workflow, validation
