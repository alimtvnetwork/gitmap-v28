# Refactor: cmd/aliasops.go

## Problem
`aliasops.go` is 235 lines with two responsibilities: core CRUD operations (set, remove, list, show) and the suggest workflow (auto-propose aliases for unaliased repos with interactive prompts).

## Target Layout

### aliasops.go (~155 lines) — CRUD Operations
Stays:
- `runAliasSet()`
- `executeAliasSet()`
- `runAliasRemove()`
- `runAliasList()`
- `printAliasList()`
- `runAliasShow()`
- `isLegacyDataError()`

### aliassuggest.go (~100 lines) — Suggest Workflow
Moves:
- `runAliasSuggest()`
- `parseAliasSuggestFlags()`
- `suggestAliases()`
- `promptAliasSuggestion()`
- `createSuggestedAlias()`

Imports: `bufio`, `flag`, `fmt`, `os`, `strings`, `constants`, `store`

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
- [69-refactor-tempreleaseops.md](69-refactor-tempreleaseops.md) — temp release branch ops
- [72-refactor-sshgen.md](72-refactor-sshgen.md) — SSH key generation
- [73-refactor-scanprojects.md](73-refactor-scanprojects.md) — project detection
- [74-refactor-amendexec.md](74-refactor-amendexec.md) — git amend operations
- [75-refactor-status.md](75-refactor-status.md) — status display
- [76-refactor-exec.md](76-refactor-exec.md) — batch execution

**Related `release/` refactors:**
- [58-refactor-workflowfinalize.md](58-refactor-workflowfinalize.md) — pipeline, metadata, zip, GitHub upload
- [64-refactor-workflow.md](64-refactor-workflow.md) — main workflow, validation
