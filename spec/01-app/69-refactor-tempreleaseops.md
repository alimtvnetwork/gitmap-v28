# Refactor: cmd/tempreleaseops.go

## Problem
`tempreleaseops.go` is 234 lines with two responsibilities: branch creation workflow (create, pattern parsing, sequence management, dry-run) and listing/display (list, print, flag detection).

## Target Layout

### tempreleaseops.go (~170 lines) — Create Workflow
Stays:
- `runTempReleaseCreate()`
- `executeTRCreate()`
- `parseVersionPattern()`
- `resolveAutoStart()`
- `validateSequenceRange()`
- `formatSeq()`
- `createTempBranches()`
- `pushTempBranches()`
- `printTRDryRun()`

### tempreleaselist.go (~80 lines) — List & Display
Moves:
- `runTempReleaseList()`
- `printTRList()`
- `hasTRListFlag()`

Imports: `encoding/json`, `fmt`, `os`, `constants`, `model`

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
- [72-refactor-sshgen.md](72-refactor-sshgen.md) — SSH key generation
- [73-refactor-scanprojects.md](73-refactor-scanprojects.md) — project detection
- [74-refactor-amendexec.md](74-refactor-amendexec.md) — git amend operations
- [75-refactor-status.md](75-refactor-status.md) — status display
- [76-refactor-exec.md](76-refactor-exec.md) — batch execution

**Related `release/` refactors:**
- [58-refactor-workflowfinalize.md](58-refactor-workflowfinalize.md) — pipeline, metadata, zip, GitHub upload
- [64-refactor-workflow.md](64-refactor-workflow.md) — main workflow, validation
