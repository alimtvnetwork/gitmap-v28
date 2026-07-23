# Refactor: release/ziparchive.go

## Problem

`ziparchive.go` is 362 lines — well over the 200-line limit. It contains
three distinct responsibilities mixed into one file:

1. **Group/ad-hoc orchestration** — resolving groups from DB, routing
   ad-hoc paths, building archive lists
2. **Low-level ZIP I/O** — creating max-compression archives, adding
   files/folders, SHA-1 hashing, summary logging
3. **Dry-run preview** — printing what archives would be produced

---

## Target Layout

Split into three files. Each stays under 200 lines.

### 1. `ziparchive.go` — Orchestration + ad-hoc routing (≈130 lines)

Keeps the public API and group/ad-hoc coordination.

**Retains:**
- `type ZipGroupArchive` (lines 22–26)
- `BuildZipGroupArchives()` (lines 30–45)
- `buildOneZipGroup()` (lines 48–91)
- `resolveArchiveName()` (lines 94–100)
- `BuildAdHocArchive()` (lines 105–111)
- `buildAdHocBundle()` (lines 114–128)
- `buildAdHocIndividual()` (lines 131–153)
- `pathsToItems()` (lines 156–175)

**Estimated:** ~130 lines (including imports and spacing).

---

### 2. `zipio.go` — Low-level ZIP creation + hashing (≈110 lines)

All file I/O, compression, and checksum logic.

**Moves:**
- `createMaxCompressZip()` (lines 178–222)
- `logArchiveSummary()` (lines 225–240)
- `sha1File()` (lines 243–258)
- `addSingleFileToZip()` (lines 261–289)
- `addFolderToZip()` (lines 292–313)

**Imports needed:** `archive/zip`, `crypto/sha1`, `encoding/hex`, `fmt`,
`io`, `os`, `path/filepath`, `model`, `verbose`

**Estimated:** ~110 lines.

---

### 3. `zipdryrun.go` — Dry-run preview (≈55 lines)

Preview output for zip groups and ad-hoc archives.

**Moves:**
- `DryRunZipGroups()` (lines 316–342)
- `DryRunAdHoc()` (lines 345–362)

**Imports needed:** `fmt`, `path/filepath`, `strings`, `constants`,
`store`

**Estimated:** ~55 lines.

---

## Migration Rules

1. **No behavior changes.** Pure file-level extraction — no renames,
   no signature changes, no logic modifications.
2. **Package stays `release`.** All files remain in `gitmap-v27/release/`.
3. **Shared functions** (`resolveArchiveName`) stay in `ziparchive.go`
   since both orchestration and dry-run reference it.
4. **Import deduplication.** Each new file declares only its own imports.
5. **Blank line before return** rule applies to all moved functions.

---

## Dependency Graph

```
ziparchive.go (orchestrator)
  ├── calls createMaxCompressZip()  → zipio.go
  ├── calls pathsToItems()          → (self)
  └── calls resolveArchiveName()    → (self)

zipio.go (I/O layer)
  ├── calls addSingleFileToZip()    → (self)
  ├── calls addFolderToZip()        → (self)
  ├── calls logArchiveSummary()     → (self)
  └── calls sha1File()              → (self)

zipdryrun.go (preview)
  ├── calls resolveArchiveName()    → ziparchive.go
  └── reads DB via store.DB         → (external)

workflowzip.go
  ├── calls BuildZipGroupArchives() → ziparchive.go
  ├── calls BuildAdHocArchive()     → ziparchive.go
  └── calls DryRunZipGroups()       → zipdryrun.go

workflowdryrun.go
  ├── calls DryRunZipGroups()       → zipdryrun.go
  └── calls DryRunAdHoc()           → zipdryrun.go
```

No circular dependencies. All calls are within the same package.

---

## Verification

1. `go build ./...` compiles without errors.
2. `go vet ./release/` passes.
3. All existing tests pass: `go test ./...`
4. No file in `release/` exceeds 200 lines after the split.
5. `wc -l release/zip*.go` confirms each file is within budget.

---

## Acceptance Criteria

1. `ziparchive.go` ≤ 200 lines.
2. Two new files created (`zipio.go`, `zipdryrun.go`), each ≤ 200 lines.
3. Zero functional changes — identical binary output.
4. All imports minimal and correctly scoped per file.
5. `go test ./...` passes.

---

## See Also

**Same package (`release/`) refactors:**

- [58-refactor-workflowfinalize.md](58-refactor-workflowfinalize.md) — pipeline, metadata, zip, GitHub upload
- [61-refactor-autocommit.md](61-refactor-autocommit.md) — auto-commit, git operations
- [63-refactor-workflowbranch.md](63-refactor-workflowbranch.md) — branch workflow, pending releases
- [64-refactor-workflow.md](64-refactor-workflow.md) — main workflow, validation
- [65-refactor-assets.md](65-refactor-assets.md) — cross-compilation, build helpers
- [78-refactor-compress.md](78-refactor-compress.md) — zip and tar.gz compression

**Related `cmd/` refactors:**
- [70-refactor-listreleases.md](70-refactor-listreleases.md) — release listing
- [71-refactor-listversions.md](71-refactor-listversions.md) — version listing
