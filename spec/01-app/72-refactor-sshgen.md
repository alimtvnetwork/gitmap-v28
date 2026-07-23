# Refactor: cmd/sshgen.go

## Problem
`sshgen.go` is 224 lines with two responsibilities: command orchestration (flag parsing, key generation, DB storage, user prompts) and utility functions (keygen validation, email resolution, fingerprint reading, path helpers).

## Target Layout

### sshgen.go (~160 lines) — Command Orchestration
Stays:
- `runSSHGenerate()`
- `parseSSHGenFlags()`
- `handleExistingKey()`
- `generateAndStore()`

### sshgenutil.go (~78 lines) — Utilities
Moves:
- `validateSSHKeygen()`
- `resolveGitEmail()`
- `readFingerprint()`
- `removeKeyFiles()`
- `defaultSSHKeyPath()`
- `expandHome()`
- `ensureSSHDir()`

Imports: `os`, `os/exec`, `path/filepath`, `strings`, `constants`

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
- [73-refactor-scanprojects.md](73-refactor-scanprojects.md) — project detection
- [74-refactor-amendexec.md](74-refactor-amendexec.md) — git amend operations
- [75-refactor-status.md](75-refactor-status.md) — status display
- [76-refactor-exec.md](76-refactor-exec.md) — batch execution

**Related `release/` refactors:**
- [58-refactor-workflowfinalize.md](58-refactor-workflowfinalize.md) — pipeline, metadata, zip, GitHub upload
- [64-refactor-workflow.md](64-refactor-workflow.md) — main workflow, validation
