# Refactor: release/workflow.go

## Problem

`workflow.go` is 291 lines — over the 200-line limit. It contains
three distinct responsibilities:

1. **Public API + orchestration** — `Execute`, `performRelease`,
   `executeSteps`, tag message resolution
2. **Version resolution** — `resolveVersion`, `resolveBump`,
   `resolveLatestVersion`, `resolveFromFile`
3. **Duplicate detection** — `checkDuplicate`, `handleOrphanedMeta`

---

## Target Layout

Split into two files. Each stays under 200 lines.

### 1. `workflow.go` — Orchestration + version resolution (≈170 lines)

Keeps the public entry point and version logic.

**Retains:**
- `type Options` (lines 20–40)
- `type Result` (lines 43–50)
- `Execute()` (lines 53–74)
- `resolveVersion()` (lines 77–107)
- `resolveBump()` (lines 110–128)
- `resolveLatestVersion()` (lines 131–148)
- `resolveFromFile()` (lines 151–160)
- `performRelease()` (lines 213–264)
- `executeSteps()` (lines 267–281)
- `resolveTagMessage()` (lines 284–291)

**Estimated:** ~170 lines (including imports and spacing).

---

### 2. `workflowvalidate.go` — Duplicate detection + orphan handling (≈60 lines)

Version existence checks and orphaned metadata cleanup.

**Moves:**
- `checkDuplicate()` (lines 164–181)
- `handleOrphanedMeta()` (lines 185–210)

**Imports needed:** `bufio`, `fmt`, `os`, `path/filepath`, `strings`,
`constants`

**Estimated:** ~60 lines.

---

## Migration Rules

1. **No behavior changes.** Pure file-level extraction — no renames,
   no signature changes, no logic modifications.
2. **Package stays `release`.** All files remain in `gitmap-v27/release/`.
3. **Import deduplication.** Each new file declares only its own imports.
4. **Blank line before return** rule applies to all moved functions.

---

## Dependency Graph

```
workflow.go (orchestrator)
  ├── calls resolveVersion()        → (self)
  ├── calls checkDuplicate()        → workflowvalidate.go
  ├── calls performRelease()        → (self)
  ├── calls executeSteps()          → (self)
  ├── calls pushAndFinalize()       → workflowfinalize.go
  └── calls writeMetadata()         → workflowfinalize.go

workflowvalidate.go (validation)
  ├── calls handleOrphanedMeta()    → (self)
  ├── calls ReleaseExists()         → metadata.go
  ├── calls TagExistsLocally()      → gitops.go
  └── calls TagExistsRemote()       → gitops.go
```

No circular dependencies. All calls are within the same package.

---

## Verification

1. `go build ./...` compiles without errors.
2. `go vet ./release/` passes.
3. All existing tests pass: `go test ./...`
4. No file in `release/` exceeds 200 lines after the split.
5. `wc -l release/workflow.go release/workflowvalidate.go`
   confirms each file is within budget.

---

## Acceptance Criteria

1. `workflow.go` ≤ 200 lines.
2. One new file created (`workflowvalidate.go`), ≤ 200 lines.
3. Zero functional changes — identical binary output.
4. All imports minimal and correctly scoped per file.
5. `go test ./...` passes.

---

## See Also

**Same package (`release/`) refactors:**

- [58-refactor-workflowfinalize.md](58-refactor-workflowfinalize.md) — pipeline, metadata, zip, GitHub upload
- [91-refactor-ziparchive.md](91-refactor-ziparchive.md) — zip I/O, dry-run, archive building
- [61-refactor-autocommit.md](61-refactor-autocommit.md) — auto-commit, git operations
- [63-refactor-workflowbranch.md](63-refactor-workflowbranch.md) — branch workflow, pending releases
- [65-refactor-assets.md](65-refactor-assets.md) — cross-compilation, build helpers
- [78-refactor-compress.md](78-refactor-compress.md) — zip and tar.gz compression

**Related `cmd/` refactors:**
- [70-refactor-listreleases.md](70-refactor-listreleases.md) — release listing
- [71-refactor-listversions.md](71-refactor-listversions.md) — version listing
