# Refactor: cmd/status.go

## Problem
`status.go` is 219 lines with two responsibilities: command orchestration with scope resolution (flag parsing, alias/group/DB/JSON loading) and terminal display (banner, table headers, summary formatting with colored output).

## Target Layout

### status.go (~133 lines) — Command & Data Loading
Stays:
- `runStatus()`
- `parseStatusFlags()`
- `loadStatusByScope()`
- `loadRecordsByGroup()`
- `loadAllRecordsDB()`
- `loadRecordsJSONFallback()`
- `loadStatusRecords()`
- `type statusSummary`

### statusprint.go (~100 lines) — Display & Formatting
Moves:
- `printStatusBanner()`
- `printStatusTable()`
- `printStatusTableTracked()`
- `printStatusHeader()`
- `printStatusSummary()`
- `buildSummaryParts()`
- `appendSummaryPart()`

Imports: `fmt`, `strings`, `cloner`, `constants`, `model`

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
- [74-refactor-amendexec.md](74-refactor-amendexec.md) — git amend operations
- [76-refactor-exec.md](76-refactor-exec.md) — batch execution

**Related `release/` refactors:**
- [58-refactor-workflowfinalize.md](58-refactor-workflowfinalize.md) — pipeline, metadata, zip, GitHub upload
- [64-refactor-workflow.md](64-refactor-workflow.md) — main workflow, validation
