# Refactor: cmd/seowriteloop.go

## Problem

`seowriteloop.go` is 340 lines — well over the 200-line limit. It contains
three distinct responsibilities mixed into one file:

1. **Commit loop orchestration** — `runCommitLoop`, `commitOne`,
   `runRotation`, `rotateLoop`, signal handling, timing
2. **Git operations** — staging, committing (with author override),
   pushing, file append/revert
3. **Output formatting** — header, commit line, rotation line, done
   summary, duration formatting

---

## Target Layout

Split into two files. Each stays under 200 lines.

### 1. `seowriteloop.go` — Loop orchestration + file resolution (≈170 lines)

Keeps the commit loop, rotation logic, file resolution, and timing.

**Retains:**
- `type commitMessage` (lines 19–22)
- `runCommitLoop()` (lines 25–45)
- `commitOne()` (lines 48–54)
- `runRotation()` (lines 57–71)
- `rotateLoop()` (lines 74–91)
- `resolvePendingFiles()` (lines 94–115)
- `pickFile()` (lines 118–124)
- `resolveRotateFile()` (lines 127–139)
- `autoDetectRotateFile()` (lines 142–151)
- `parseInterval()` (lines 154–169)
- `waitRandom()` (lines 172–182)
- `setupSignalHandler()` (lines 185–197)
- `shouldStop()` (lines 200–212)

**Estimated:** ~170 lines (including imports and spacing).

---

### 2. `seowritegit.go` — Git ops + output formatting (≈130 lines)

Git commands, file manipulation, and all print helpers.

**Moves:**
- `gitStage()` (lines 215–220)
- `gitCommit()` (lines 223–225)
- `gitCommitWithAuthor()` (lines 228–245)
- `resolveAuthorFlag()` (lines 248–260)
- `gitPush()` (lines 263–268)
- `appendToFile()` (lines 271–279)
- `revertFile()` (lines 282–290)
- `printHeader()` (lines 293–301)
- `printCommitLine()` (lines 304–312)
- `printRotationLine()` (lines 315–323)
- `printDone()` (lines 326–328)
- `formatDuration()` (lines 331–340)

**Imports needed:** `fmt`, `os`, `os/exec`, `strings`, `time`,
`constants`

**Estimated:** ~130 lines.

---

## Migration Rules

1. **No behavior changes.** Pure file-level extraction — no renames,
   no signature changes, no logic modifications.
2. **Package stays `cmd`.** All files remain in `gitmap-v27/cmd/`.
3. **Import deduplication.** Each new file declares only its own imports.
4. **Blank line before return** rule applies to all moved functions.

---

## Dependency Graph

```
seowriteloop.go (orchestrator)
  ├── calls commitOne()             → (self)
  ├── calls runRotation()           → (self)
  ├── calls gitStage()              → seowritegit.go
  ├── calls gitCommitWithAuthor()   → seowritegit.go
  ├── calls gitPush()               → seowritegit.go
  ├── calls appendToFile()          → seowritegit.go
  ├── calls revertFile()            → seowritegit.go
  ├── calls printHeader()           → seowritegit.go
  ├── calls printCommitLine()       → seowritegit.go
  ├── calls printRotationLine()     → seowritegit.go
  └── calls printDone()             → seowritegit.go

seowritegit.go (git + output)
  ├── calls resolveAuthorFlag()     → (self)
  └── calls formatDuration()        → (self)
```

No circular dependencies. All calls are within the same package.

---

## Verification

1. `go build ./...` compiles without errors.
2. `go vet ./cmd/` passes.
3. All existing tests pass: `go test ./...`
4. No file in `cmd/` exceeds 200 lines after the split.
5. `wc -l cmd/seowrite*.go` confirms each file is within budget.

---

## Acceptance Criteria

1. `seowriteloop.go` ≤ 200 lines.
2. One new file created (`seowritegit.go`), ≤ 200 lines.
3. Zero functional changes — identical binary output.
4. All imports minimal and correctly scoped per file.
5. `go test ./...` passes.

---

## See Also

**Same package (`cmd/`) refactors:**

- [90-refactor-root-dispatch.md](90-refactor-root-dispatch.md) — dispatch splitting
- [66-refactor-zipgroupops.md](66-refactor-zipgroupops.md) — zip group CRUD and display
- [68-refactor-aliasops.md](68-refactor-aliasops.md) — alias CRUD and suggest
- [69-refactor-tempreleaseops.md](69-refactor-tempreleaseops.md) — temp release branch ops
- [72-refactor-sshgen.md](72-refactor-sshgen.md) — SSH key generation
- [73-refactor-scanprojects.md](73-refactor-scanprojects.md) — project detection
- [74-refactor-amendexec.md](74-refactor-amendexec.md) — git amend operations
- [75-refactor-status.md](75-refactor-status.md) — status display
- [76-refactor-exec.md](76-refactor-exec.md) — batch execution

**Related `release/` refactors:**
- [58-refactor-workflowfinalize.md](58-refactor-workflowfinalize.md) — pipeline, metadata, zip, GitHub upload
- [64-refactor-workflow.md](64-refactor-workflow.md) — main workflow, validation
