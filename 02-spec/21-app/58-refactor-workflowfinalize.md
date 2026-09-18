# Refactor: release/workflowfinalize.go

## Problem

`workflowfinalize.go` is 498 lines — well over the 200-line limit. It
contains four distinct responsibilities mixed into one file:

1. **Finalization pipeline** — push, asset collection, compress, checksums
2. **Dry-run preview** — printing what would happen
3. **Zip/ad-hoc archive building** — DB-backed and ad-hoc zip groups
4. **GitHub upload** — token detection, release creation, asset upload
5. **Metadata persistence** — writing JSON, updating latest.json

---

## Target Layout

Split into four files. Each file stays under 200 lines.

### 1. `workflowfinalize.go` — Core pipeline + metadata (≈150 lines)

Keeps the orchestration entry point and metadata functions.

**Retains:**
- `var LastMeta *ReleaseMeta` (line 15)
- `var lastZipChecksums map[string]string` (line 20)
- `pushAndFinalize()` (lines 22–75)
- `writeMetadata()` (lines 172–192)
- `buildReleaseMeta()` (lines 194–225)
- `collectZipGroupNames()` (lines 227–241)
- `loadChangelogNotes()` (lines 243–256)
- `updateLatestIfStable()` (lines 258–286)
- `returnToBranch()` (lines 373–387)

**Estimated:** ~145 lines (including imports and spacing).

---

### 2. `workflowdryrun.go` — Dry-run preview (≈120 lines)

All `printDryRun*` functions.

**Moves:**
- `printDryRun()` (lines 288–298)
- `printDryRunGoAssets()` (lines 300–326)
- `printDryRunSteps()` (lines 328–342)
- `printDryRunAssets()` (lines 344–362)
- `printDryRunMeta()` (lines 364–371)
- `printDryRunZipGroups()` (lines 486–498)

**Imports needed:** `fmt`, `constants`, `store`

**Estimated:** ~115 lines.

---

### 3. `workflowzip.go` — Zip group and ad-hoc archive building (≈120 lines)

All zip-related asset building and checksum collection.

**Moves:**
- `buildZipGroupAssets()` (lines 389–431)
- `buildAdHocZipAssets()` (lines 433–464)
- `collectZipChecksums()` (lines 466–484)

**Imports needed:** `fmt`, `os`, `path/filepath`, `constants`, `store`, `verbose`

**Estimated:** ~105 lines.

---

### 4. `workflowgithub.go` — GitHub release + Go asset building (≈130 lines)

GitHub upload logic and Go cross-compilation.

**Moves:**
- `uploadToGitHub()` (lines 127–170)
- `buildGoAssetsIfApplicable()` (lines 77–125)

**Imports needed:** `fmt`, `os`, `constants`, `verbose`

**Estimated:** ~105 lines.

---

## Migration Rules

1. **No behavior changes.** Pure file-level extraction — no renames,
   no signature changes, no logic modifications.
2. **Package stays `release`.** All files remain in `gitmap-v28/release/`.
3. **Shared state** (`LastMeta`, `lastZipChecksums`) stays in
   `workflowfinalize.go` since it's the orchestration hub. Other files
   reference these package-level vars directly.
4. **Import deduplication.** Each new file declares only its own imports.
5. **Blank line before return** rule applies to all moved functions.

---

## Dependency Graph

```
workflowfinalize.go (orchestrator)
  ├── calls buildGoAssetsIfApplicable()  → workflowgithub.go
  ├── calls buildZipGroupAssets()        → workflowzip.go
  ├── calls buildAdHocZipAssets()        → workflowzip.go
  ├── calls uploadToGitHub()             → workflowgithub.go
  └── calls writeMetadata()              → (self)

workflow.go
  ├── calls pushAndFinalize()            → workflowfinalize.go
  ├── calls writeMetadata()              → workflowfinalize.go
  ├── calls printDryRun()                → workflowdryrun.go
  └── calls returnToBranch()             → workflowfinalize.go

workflowbranch.go
  ├── calls pushAndFinalize()            → workflowfinalize.go
  └── calls printDryRun()                → workflowdryrun.go
```

No circular dependencies. All calls are within the same package.

---

## Verification

1. `go build ./...` compiles without errors.
2. `go vet ./release/` passes.
3. All existing tests pass: `go test ./...`
4. No file in `release/` exceeds 200 lines after the split.
5. `wc -l release/workflow*.go` confirms each file is within budget.

---

## Acceptance Criteria

1. `workflowfinalize.go` ≤ 200 lines.
2. Three new files created, each ≤ 200 lines.
3. Zero functional changes — identical binary output.
4. All imports minimal and correctly scoped per file.
5. `go test ./...` passes.

---

## See Also

**Same package (`release/`) refactors:**

- [91-refactor-ziparchive.md](91-refactor-ziparchive.md) — zip I/O, dry-run, archive building
- [61-refactor-autocommit.md](61-refactor-autocommit.md) — auto-commit, git operations
- [63-refactor-workflowbranch.md](63-refactor-workflowbranch.md) — branch workflow, pending releases
- [64-refactor-workflow.md](64-refactor-workflow.md) — main workflow, validation
- [65-refactor-assets.md](65-refactor-assets.md) — cross-compilation, build helpers
- [78-refactor-compress.md](78-refactor-compress.md) — zip and tar.gz compression

**Related `cmd/` refactors:**
- [70-refactor-listreleases.md](70-refactor-listreleases.md) — release listing
- [71-refactor-listversions.md](71-refactor-listversions.md) — version listing
