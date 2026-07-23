# Refactor: release/workflowbranch.go

## Problem

`workflowbranch.go` is 310 lines — over the 200-line limit. It contains
three distinct responsibilities:

1. **Branch release** — `ExecuteFromBranch`, branch validation,
   checkout and tag/push completion
2. **Pending discovery** — `ExecutePending`, branch listing, filtering,
   metadata-based pending detection
3. **Metadata-based release** — creating branch+tag from stored commit
   SHAs, dry-run output for metadata releases

---

## Target Layout

Split into two files. Each stays under 200 lines.

### 1. `workflowbranch.go` — Branch release + pending orchestration (≈160 lines)

Keeps the public API for branch-based and pending releases.

**Retains:**
- `ExecuteFromBranch()` (lines 12–32)
- `extractVersionFromBranch()` (lines 35–44)
- `validateExistingBranch()` (lines 47–53)
- `completeBranchRelease()` (lines 56–94)
- `ExecutePending()` (lines 98–126)
- `releasePendingBranches()` (lines 236–246)
- `listReleaseBranches()` (lines 249–257)
- `parseBranchLines()` (lines 260–273)
- `filterPendingBranches()` (lines 276–286)
- `isPendingBranch()` (lines 289–298)
- `tagIsMissing()` (lines 301–310)

**Estimated:** ~160 lines (including imports and spacing).

---

### 2. `workflowpending.go` — Metadata-based pending release (≈110 lines)

Metadata discovery and commit-SHA-based release creation.

**Moves:**
- `discoverMetadataPending()` (lines 130–150)
- `isMetaPending()` (lines 154–171)
- `releasePendingFromMetadata()` (lines 174–184)
- `releaseFromMetadata()` (lines 187–233)

**Imports needed:** `fmt`, `constants`

**Estimated:** ~110 lines.

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
workflowbranch.go (orchestrator)
  ├── calls extractVersionFromBranch()   → (self)
  ├── calls validateExistingBranch()     → (self)
  ├── calls completeBranchRelease()      → (self)
  ├── calls filterPendingBranches()      → (self)
  ├── calls discoverMetadataPending()    → workflowpending.go
  ├── calls releasePendingBranches()     → (self)
  ├── calls releasePendingFromMetadata() → workflowpending.go
  └── calls pushAndFinalize()            → workflowfinalize.go

workflowpending.go (metadata releases)
  ├── calls isMetaPending()              → (self)
  ├── calls releaseFromMetadata()        → (self)
  └── calls pushAndFinalize()            → workflowfinalize.go
```

No circular dependencies. All calls are within the same package.

---

## Verification

1. `go build ./...` compiles without errors.
2. `go vet ./release/` passes.
3. All existing tests pass: `go test ./...`
4. No file in `release/` exceeds 200 lines after the split.
5. `wc -l release/workflowbranch.go release/workflowpending.go`
   confirms each file is within budget.

---

## Acceptance Criteria

1. `workflowbranch.go` ≤ 200 lines.
2. One new file created (`workflowpending.go`), ≤ 200 lines.
3. Zero functional changes — identical binary output.
4. All imports minimal and correctly scoped per file.
5. `go test ./...` passes.

---

## See Also

**Same package (`release/`) refactors:**

- [58-refactor-workflowfinalize.md](58-refactor-workflowfinalize.md) — pipeline, metadata, zip, GitHub upload
- [91-refactor-ziparchive.md](91-refactor-ziparchive.md) — zip I/O, dry-run, archive building
- [61-refactor-autocommit.md](61-refactor-autocommit.md) — auto-commit, git operations
- [64-refactor-workflow.md](64-refactor-workflow.md) — main workflow, validation
- [65-refactor-assets.md](65-refactor-assets.md) — cross-compilation, build helpers
- [78-refactor-compress.md](78-refactor-compress.md) — zip and tar.gz compression

**Related `cmd/` refactors:**
- [70-refactor-listreleases.md](70-refactor-listreleases.md) — release listing
- [71-refactor-listversions.md](71-refactor-listversions.md) — version listing
