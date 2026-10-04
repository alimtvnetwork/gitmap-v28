# Completed Plan: 66-fix-gitmap-pa-duplicate-repos

## User Request (Verbatim)
most of the packages seems like repeated when did the pull-all or pa 

gitmap pa

```
Failed Repositories (37):
    • ai-empathy-prompt-tuner-v1             failed
        ↳ Reason: fatal: Cannot fast-forward to multiple branches.
...
    • wp-onboarding-v17                      failed
        ↳ Reason: fatal: Cannot fast-forward to multiple branches.

  ✓ Pull all complete: 139 pulled (97 active, 42 up-to-date) (55.5s)

  ⚠ Detected .gitignore issues in 32 repository(ies):
    • cat-my-v12
...
```
Can you please find the root cause and try to fix it as well

---

## 1. Root Cause Summary
- **Primary Root Cause (Database Case-Sensitivity)**: In `gitmap.db`, `IdxRepo_AbsolutePath` on `Repo(AbsolutePath)` and `IdxScanFolder_AbsolutePath` on `ScanFolder(AbsolutePath)` lacked `COLLATE NOCASE`. On Windows, drive letter casing variants (`D:\work` vs `d:\work`) were stored as distinct rows, accumulating 139 rows in `Repo` for only 78 physical repositories (61 repos duplicated).
- **Secondary Root Cause (Missing In-Memory Deduplication)**: `loadAllRecordsDB()` and `resolvePullTargets()` directly queried `SELECT ... FROM Repo ORDER BY Slug` via `db.ListRepos()` and passed all 139 records to `pull-all` worker channels and the `.gitignore` background scanner without canonical path or slug deduplication.
- **Tertiary Root Cause (Worker Concurrency Collision)**: Parallel workers dispatched `./X` and `./X` simultaneously on the same `.git` folder, triggering concurrent `git fetch` operations that wrote competing branch heads into `.git/FETCH_HEAD`, causing `fatal: Cannot fast-forward to multiple branches.`.
- **Quaternary Root Cause (Unfiltered Reporting)**: Both the pull failed summary (`renderFailedGroup`) and `.gitignore` audit (`printIgnoreIssuesReport`) printed every worker state without deduplication, inflating the failure tally to 37 and listing duplicate bullet points.

---

## 2. Remediations Implemented

### Task-01: Database Schema Collation & Deduplication Migration (Worker 01)
- **`cli/constants/constants_store.go`**:
  - Updated `SQLCreateAbsPathIndex` to `CREATE UNIQUE INDEX IF NOT EXISTS IdxRepo_AbsolutePath ON Repo(AbsolutePath COLLATE NOCASE)`.
  - Added `SQLDropRepoAbsPathIndex` and `SQLDeduplicateRepos`.
- **`cli/constants/constants_scan_folder.go`**:
  - Updated `SQLCreateScanFolderPathIndex` to `CREATE UNIQUE INDEX IF NOT EXISTS IdxScanFolder_AbsolutePath ON ScanFolder(AbsolutePath COLLATE NOCASE)`.
  - Added `SQLDropScanFolderPathIndex` and `SQLDeduplicateScanFolders`.
- **`cli/store/migrations.go`**:
  - Implemented migration `Migration_AddRepoAbsolutePathCollateNoCase` (Schema Version 33):
    - Remaps child foreign keys in `Release`, `GroupRepo`, and `VersionProbe` to the surviving `MIN(RepoId)`.
    - Deletes duplicate `Repo` rows grouped by `LOWER(AbsolutePath)`.
    - Drops and recreates `IdxRepo_AbsolutePath` with `COLLATE NOCASE`.
    - Remaps `Repo.ScanFolderId` pointers and deduplicates `ScanFolder` rows.
    - Recreates `IdxScanFolder_AbsolutePath` with `COLLATE NOCASE`.
    - Bumps schema version to 33.
- **`cli/store/repo.go`**:
  - Added `NormalizeStoragePath` for canonical write-time path formatting with capitalized Windows drive letters.
  - Applied `NormalizeStoragePath` in `upsertOneRepo`, `DeleteByPath`, and `FindByPath`.
- **`cli/cmd/reconcile_db.go`**:
  - Replaced flawed substring check in `isSubPath` with boundary-safe, separator-aware, case-insensitive logic.
  - Used case-folded normalized paths in `validPaths` set.

### Task-02: Pull Pipeline In-Memory Deduplication & Worker Safety (Worker 02)
- **`cli/cmdpull/pull.go` & `cli/cmdpull/helpers.go`**:
  - Implemented `CanonicalRepoPathKey(path string) string` converting paths via `filepath.Clean()`, `filepath.ToSlash()`, and lowercasing.
  - Implemented `deduplicatePullRecords(records []model.ScanRecord) []model.ScanRecord` to deduplicate by canonical path and repository slug.
  - Wrapped `loadAllRecordsDB()`, `resolvePullTargets()`, `findChildrenOfCWD()`, and `resolvePullBatchRecords()`.
  - Implemented `DeduplicateIgnoreIssues(issues []IgnoreRepoIssue) []IgnoreRepoIssue` to eliminate duplicate `.gitignore` issue reporting.
- **`cli/cmdpull/pull_concurrency.go`**:
  - Added canonical path tracking in `scanColdRecordsSequentially` to skip checking identical repositories across workers.
- **`cli/cmdpull/pull_efficient_render.go`**:
  - Implemented `DeduplicateRepoStates(states []*PullRepoState) []*PullRepoState` with state prioritization.
  - Deduplicated `failed` slice in `renderFailedGroup` so `len(failed)` reflects unique repositories.
- **`cli/cloner/safe_pull.go` & `cli/cloner/pulldiag.go`**:
  - Added detection for `Cannot fast-forward to multiple branches.` in `isDivergedOutput` to trigger auto-merge fallback.
  - Added diagnosis recognition for multi-branch `FETCH_HEAD` conflicts in `pulldiag.go`.

---

## 3. Specifications & Verification Evidence
- **Specs**:
  - `02-spec/21-app/66-fix-gitmap-pa-duplicate-repos/01-architecture-spec.md`
  - `02-spec/21-app/66-fix-gitmap-pa-duplicate-repos/02-component-spec.md`
- **Verification Evidence**:
  - `python 03-ai-scripts/05-guideline-autofixer.py cli/store --check-only --ext .go`: 193 files PASS
  - `python 03-ai-scripts/05-guideline-autofixer.py cli/cmdpull --check-only --ext .go`: 57 files PASS
  - `python 03-ai-scripts/05-guideline-autofixer.py cli/cloner --check-only --ext .go`: 24 files PASS
  - `python 03-ai-scripts/05-guideline-autofixer.py cli/constants --check-only --ext .go`: 142 files PASS
  - `python 03-ai-scripts/05-guideline-autofixer.py cli/cmd --check-only --ext .go`: 829 files PASS
  - `python linter-scripts/check-relative-paths.py`: PASS across 8,721 files
  - `python linter-scripts/check-forbidden-strings.py`: PASS
  - `python linter-scripts/check-prompts-loaded.py`: PASS
  - Secrets Gate: PASS (0 hits across 3,874 files)
