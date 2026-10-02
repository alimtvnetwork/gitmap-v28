# App Issue 62: Pull-All Duplicate Repositories, Database Case-Sensitivity & Concurrent Fetch Failures RCA

**Issue ID:** 62  
**Date:** 2026-10-02  
**Status:** Resolved  
**Affected Subsystem:** `cli/cmdpull`, `cli/store`, `cli/cloner`, `cli/constants`

---

## 1. Reproduction

Running `gitmap pa` (`gitmap pull-all`) on a Windows machine where repository paths were registered across different terminal sessions:

```text
gitmap pa
Failed Repositories (37):
    • ai-empathy-prompt-tuner-v1             failed
        ↳ Reason: fatal: Cannot fast-forward to multiple branches.
    • ai-empathy-prompt-tuner-v1             failed
        ↳ Reason: fatal: Cannot fast-forward to multiple branches.
    • prompt-architect-v2                    failed
        ↳ Reason: fatal: Cannot fast-forward to multiple branches.
...
  ✓ Pull all complete: 139 pulled (97 active, 42 up-to-date) (55.5s)

  ⚠ Detected .gitignore issues in 32 repository(ies):
    • cat-my-v12
      - Duplicate pattern in .gitignore: /scratch/
    • cat-my-v12
      - Duplicate pattern in .gitignore: /scratch/
```

### Symptoms:
1. Every failed repository was listed at least twice (some up to 4 times).
2. The pull failure reason was `fatal: Cannot fast-forward to multiple branches.`.
3. Total pulled count was inflated to 139 for only 78 physical repositories on disk.
4. `.gitignore` issues were reported redundantly for identical repositories.

---

## 2. Root Cause Analysis (4-Part RCA)

### Symptom & Discovery:
`gitmap pa` dispatched 139 pull jobs and reported 37 failures containing repeated entries.

### Direct Cause:
1. **SQLite Default Binary Collation**: `IdxRepo_AbsolutePath` on `Repo(AbsolutePath)` lacked `COLLATE NOCASE`. Under SQLite default binary collation, `D:\work\...` and `d:\work\...` were treated as distinct strings.
2. **Missing In-Memory Deduplication**: `loadAllRecordsDB()` and `resolvePullTargets()` queried `SELECT ... FROM Repo ORDER BY Slug` via `db.ListRepos()` and passed all 139 rows to `pull-all` worker channels and the `.gitignore` background scanner without canonical path or slug deduplication.

### Compound Failure:
When parallel workers consumed the job queue, Worker A pulled `D:\work\gitmap` and Worker B pulled `d:\work\gitmap` simultaneously on the same `.git` directory. Concurrent `git fetch` operations wrote multiple branch heads into `.git/FETCH_HEAD`. Git's merge fast-forward logic (`builtin/pull.c`) saw multiple merge heads (`merge_heads.nr > 1`) and aborted with `fatal: Cannot fast-forward to multiple branches.`.

### Reporting Flaw:
The concise output renderer and ignore scanner had no deduplication on `failed` states or `issues`, multiplying terminal output and inflating failure metrics.

---

## 3. Fix & Remediation

1. **Database Schema & Migration (Schema Version 33)**:
   - Updated `SQLCreateAbsPathIndex` in `cli/constants/constants_store.go` to use `ON Repo(AbsolutePath COLLATE NOCASE)`.
   - Updated `SQLCreateScanFolderPathIndex` in `cli/constants/constants_scan_folder.go` to use `ON ScanFolder(AbsolutePath COLLATE NOCASE)`.
   - Implemented `Migration_AddRepoAbsolutePathCollateNoCase` in `cli/store/migrations.go` to deduplicate `Repo` and `ScanFolder` rows keeping `MIN(RepoId)` and recreate indexes with `COLLATE NOCASE`.
   - Added `NormalizeStoragePath` in `cli/store/repo.go` to enforce uniform drive-letter casing.
   - Refactored `isSubPath` in `cli/cmd/reconcile_db.go` to be case-insensitive on Windows.
2. **In-Memory Target Deduplication**:
   - Added `CanonicalRepoPathKey` and `deduplicatePullRecords` in `cli/cmdpull/pull.go` to filter loaded targets by canonical path and slug.
   - Wrapped `loadAllRecordsDB`, `resolvePullTargets`, `findChildrenOfCWD`, and `resolvePullBatchRecords`.
3. **Ignore & Render Deduplication**:
   - Added canonical path tracking in `scanColdRecordsSequentially` in `cli/cmdpull/pull_concurrency.go`.
   - Added `DeduplicateIgnoreIssues` in `cli/cmdpull/pull.go`.
   - Added `DeduplicateRepoStates` in `cli/cmdpull/pull_efficient_render.go` so `renderFailedGroup` and `len(failed)` count distinct repositories.
4. **Worker Safety**:
   - Updated `isDivergedOutput` in `cli/cloner/safe_pull.go` to detect `Cannot fast-forward to multiple branches.` and trigger auto-merge fallback.
   - Added diagnostic hint for multiple branch conflicts in `cli/cloner/pulldiag.go`.

---

## 4. Prevention

1. **Collation Guideline**: All path-based unique indexes in SQLite databases must explicitly declare `COLLATE NOCASE` when indexing Windows filesystem paths.
2. **Input Normalization**: Ingestion points must clean paths using `filepath.Clean` and normalize drive letter casing before persistence.
3. **Defense-in-Depth Deduplication**: Batch execution pipelines (`pull-all`, `status-all`, `clone-all`) must always pass targets through in-memory deduplication sets to guard against database anomalies.
