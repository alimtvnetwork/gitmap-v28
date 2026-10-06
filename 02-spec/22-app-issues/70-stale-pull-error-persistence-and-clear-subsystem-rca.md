# Issue 70: Stale Pull Error Persistence, Missing Eviction, and Clear Subsystem RCA

> **Issue ID:** `ISSUE-70`
> **Target Subsystem:** `cli/cmdpull`, `cli/cmdpullerror`, `cli/store` (Pull Split-DB)
> **Affected Repository:** `gitlogger-new-v2` (`d:\work\gitlogger-new`)
> **Status:** Resolved
> **Date:** 2026-10-06

---

## 1. Problem Reproduction & Symptom Analysis

### User Invocation & Observed Error
When running `gitmap pull-error gitlogger-new-v2` (or alias `gitmap pulle gitlogger-new-v2`):
```text
safe-pull failed for gitlogger-new-v2 (gitlogger-new) (attempt 1/4):
repo="D:\work\gitlogger-new" branch="main": [E_MERGE_CONFLICT:EXECUTION] pull_merge:
git merge conflict detected: auto-merge pull aborted (at=cloner/safe_pull.go:28) (creator=safe_pull)

Auto-merging .ai-memory/plans/01-index.md
CONFLICT (content): Merge conflict in .ai-memory/plans/01-index.md
Auto-merging .gitignore
CONFLICT (content): Merge conflict in .gitignore
Auto-merging src/components/slides/edit/SelectionOverlay.tsx
CONFLICT (add/add): Merge conflict in src/components/slides/edit/SelectionOverlay.tsx
Auto-merging src/components/slides/edit/applyEdit.ts
CONFLICT (add/add): Merge conflict in src/components/slides/edit/applyEdit.ts
...
Auto-merging src/lib/db/index.ts
CONFLICT (add/add): Merge conflict in src/lib/db/index.t...
Diagnosis: non-unlink git pull failure (check auth/merge or run pull manually for full output)
Actionable Remediation Command:
  gitmap fix gitlogger-new-v2 or gitmap stash
```

However, checking `d:\work\gitlogger-new` showed:
- Branch `main` was clean and fully synchronized with upstream `origin/main` (`Already up to date.`).
- Executing `gitmap pull gitlogger-new-v2` reported 100% success (`✔ active`).
- Despite the repository being completely healthy and up-to-date, `gitmap pull-error gitlogger-new-v2` continuously dumped the stale historical merge conflict diagnostic.

---

## 2. Root Cause Analysis (4-Part RCA)

### 1. Symptom
`gitmap pull-error gitlogger-new-v2` continuously reported historical merge conflict failures `[E_MERGE_CONFLICT:EXECUTION]` even though the repository branch divergence had already been resolved and subsequent pulls succeeded cleanly.

### 2. Root Cause
GitMap's pull telemetry synchronization (`cli/cmdpull/pull_db_sync.go` and `cli/cmdpull/pull_efficient.go`) only wrote failure records to SQLite `pull_errors` without evicting or resolving previous records on successful pull operations, and `cli/cmdpullerror/pull_error_cmd.go` lacked eviction logic and an explicit `clear` subcommand or `--clear` flag.

### 3. Resolution
1. **Reconciled Repository Validation:** Verified `d:\work\gitlogger-new` has resolved all merge conflicts, merged upstream commit `7017153`, and committed documentation (`6dd9cf7`).
2. **Database Eviction Engine:** Added `ClearPullErrors(repoSlug string)` and `ClearPullErrorsForRepo(repoSlug, repoPath string)` to `*PullSplitDB` in [pull_split_db_errors_clear.go](file:///d:/work/gitmap/cli/store/pull_split_db_errors_clear.go).
3. **Automatic Lifecycle Resolution:** Updated [pull_db_sync.go](file:///d:/work/gitmap/cli/cmdpull/pull_db_sync.go) and [pull_efficient.go](file:///d:/work/gitmap/cli/cmdpull/pull_efficient.go) to automatically evict resolved repository error entries whenever a repository pulls successfully (`!isStateFailure`).
4. **Manual Administrative Clearing:** Implemented `gitmap pull-error clear [target]` and `gitmap pull-error --clear` / `-c` support in [pull_error_cmd.go](file:///d:/work/gitmap/cli/cmdpullerror/pull_error_cmd.go) and [pull_error_clear.go](file:///d:/work/gitmap/cli/cmdpullerror/pull_error_clear.go), with remote fleet forwarding for SSH clusters.
5. **Telemetry Sanitization:** Purged the 7 stale historical error rows for `gitlogger-new-v2` from `gitmap-pull.db`, verifying `gitmap pull-error gitlogger-new-v2` returns clean zero-error status (`✔ No pull errors recorded for gitlogger-new-v2.`).

### 4. Prevention & Learnings
Diagnostic and error reporting stores must implement a symmetrical lifecycle: every error recording event on failure must have a corresponding eviction or resolution transition on success. Furthermore, inspection CLIs must always provide first-class administrative verbs (`clear` / `--clear`) to allow operators and agents to explicitly purge resolved diagnostics.

---

## 3. Verification Evidence

1. `gitmap pull gitlogger-new-v2`:
   ```text
   [====================] 100% (1/1 repos) | [=] gitlogger-new-v2: up-to-date
   STATUS: ✔ active
   ```
2. `gitmap pull-error gitlogger-new-v2`:
   ```text
   ✔ No pull errors recorded for gitlogger-new-v2.
   ```
3. `gitmap pull-error gitlogger-new-v2 --json`:
   ```json
   []
   ```
4. Guideline compliance:
   - `python 03-ai-scripts/05-guideline-autofixer.py --check-only cli/store` -> 0 violations.
   - `python 03-ai-scripts/05-guideline-autofixer.py --check-only cli/cmdpull` -> 0 violations.
   - `python 03-ai-scripts/05-guideline-autofixer.py --check-only cli/cmdpullerror` -> 0 violations.
