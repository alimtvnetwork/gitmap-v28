# Refactor: release/autocommit.go

## Problem

`autocommit.go` is 352 lines — well over the 200-line limit. It contains
three distinct responsibilities mixed into one file:

1. **Auto-commit orchestration** — entry point, file classification,
   user prompting, commit/push coordination
2. **Git staging/commit primitives** — staging files, committing, push
   with rebase-retry recovery
3. **Git output helpers** — porcelain parsing, error formatting,
   non-fast-forward detection, output trimming

---

## Target Layout

Split into two files. Each stays under 200 lines.

### 1. `autocommit.go` — Orchestration + classification (≈180 lines)

Keeps the public API, file classification, and commit workflows.

**Retains:**
- `type AutoCommitResult` (lines 14–19)
- `AutoCommit()` (lines 24–61)
- `listChangedFiles()` (lines 64–72)
- `parsePorcelainOutput()` (lines 75–91)
- `classifyFiles()` (lines 94–105)
- `commitReleaseOnly()` (lines 108–149)
- `promptAndCommit()` (lines 152–179)
- `commitAll()` (lines 182–223)

**Estimated:** ~180 lines (including imports and spacing).

---

### 2. `autocommitgit.go` — Git primitives + push recovery (≈120 lines)

Low-level git operations, push-with-rebase logic, and output helpers.

**Moves:**
- `stageFiles()` (lines 226–230)
- `stageAll()` (lines 233–235)
- `commitStaged()` (lines 238–240)
- `pushCurrentBranch()` (lines 243–259)
- `syncBranchAndRetryPush()` (lines 261–304)
- `runGitCmdCombined()` (lines 306–311)
- `isNonFastForwardPushError()` (lines 313–319)
- `formatGitCommandError()` (lines 321–328)
- `trimGitOutput()` (lines 330–337)
- `singleLineGitOutput()` (lines 339–341)
- `abortRebaseAfterFailure()` (lines 343–352)

**Imports needed:** `fmt`, `os/exec`, `strings`, `constants`, `verbose`

**Estimated:** ~120 lines.

---

## Migration Rules

1. **No behavior changes.** Pure file-level extraction — no renames,
   no signature changes, no logic modifications.
2. **Package stays `release`.** All files remain in `gitmap-v28/release/`.
3. **Import deduplication.** Each new file declares only its own imports.
4. **Blank line before return** rule applies to all moved functions.

---

## Dependency Graph

```
autocommit.go (orchestrator)
  ├── calls listChangedFiles()       → (self)
  ├── calls classifyFiles()          → (self)
  ├── calls stageFiles()             → autocommitgit.go
  ├── calls stageAll()               → autocommitgit.go
  ├── calls commitStaged()           → autocommitgit.go
  └── calls pushCurrentBranch()      → autocommitgit.go

autocommitgit.go (git primitives)
  ├── calls runGitCmdCombined()      → (self)
  ├── calls isNonFastForwardPushError() → (self)
  ├── calls syncBranchAndRetryPush() → (self)
  ├── calls formatGitCommandError()  → (self)
  ├── calls trimGitOutput()          → (self)
  ├── calls singleLineGitOutput()    → (self)
  └── calls abortRebaseAfterFailure() → (self)

workflowfinalize.go
  └── (no direct calls — AutoCommit called from workflow.go)

workflow.go
  └── calls AutoCommit()             → autocommit.go
```

No circular dependencies. All calls are within the same package.

---

## Verification

1. `go build ./...` compiles without errors.
2. `go vet ./release/` passes.
3. All existing tests pass: `go test ./...`
4. No file in `release/` exceeds 200 lines after the split.
5. `wc -l release/autocommit*.go` confirms each file is within budget.

---

## Acceptance Criteria

1. `autocommit.go` ≤ 200 lines.
2. One new file created (`autocommitgit.go`), ≤ 200 lines.
3. Zero functional changes — identical binary output.
4. All imports minimal and correctly scoped per file.
5. `go test ./...` passes.

---

## See Also

**Same package (`release/`) refactors:**

- [58-refactor-workflowfinalize.md](58-refactor-workflowfinalize.md) — pipeline, metadata, zip, GitHub upload
- [91-refactor-ziparchive.md](91-refactor-ziparchive.md) — zip I/O, dry-run, archive building
- [63-refactor-workflowbranch.md](63-refactor-workflowbranch.md) — branch workflow, pending releases
- [64-refactor-workflow.md](64-refactor-workflow.md) — main workflow, validation
- [65-refactor-assets.md](65-refactor-assets.md) — cross-compilation, build helpers
- [78-refactor-compress.md](78-refactor-compress.md) — zip and tar.gz compression

**Related `cmd/` refactors:**
- [70-refactor-listreleases.md](70-refactor-listreleases.md) — release listing
- [71-refactor-listversions.md](71-refactor-listversions.md) — version listing
