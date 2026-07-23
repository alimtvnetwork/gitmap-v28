# Refactor: cmd/amendexec.go

## Problem
`amendexec.go` is 223 lines with two responsibilities: git operations for amending (listing commits, parsing output, filter-branch, checkout) and output/display helpers (env-filter construction, author string building, force push, progress/dry-run printing).

## Target Layout

### amendexec.go (~130 lines) — Git Operations
Stays:
- `listCommitsForAmend()`
- `parseCommitLines()`
- `detectPreviousAuthor()`
- `getCurrentBranch()`
- `switchBranch()`
- `runFilterBranch()`
- `runAmendHead()`

### amendexecprint.go (~100 lines) — Output & Display
Moves:
- `buildEnvFilter()`
- `buildAuthorString()`
- `runForcePush()`
- `printAmendHeader()`
- `printAmendProgress()`
- `printAmendDryRun()`

Imports: `fmt`, `os`, `os/exec`, `strings`, `constants`, `model`

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
- [73-refactor-scanprojects.md](73-refactor-scanprojects.md) — project detection
- [75-refactor-status.md](75-refactor-status.md) — status display
- [76-refactor-exec.md](76-refactor-exec.md) — batch execution

**Related `release/` refactors:**
- [58-refactor-workflowfinalize.md](58-refactor-workflowfinalize.md) — pipeline, metadata, zip, GitHub upload
- [64-refactor-workflow.md](64-refactor-workflow.md) — main workflow, validation
