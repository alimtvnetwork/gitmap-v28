# Component Technical Specification: Pull Pipeline Deduplication, Ignore Auditing Deduplication, Summary Render Deduplication, and Worker Safety

- **Spec Document:** `02-spec/21-app/66-fix-gitmap-pa-duplicate-repos/02-component-spec.md`
- **Spec ID:** 66-02
- **Status:** Approved / Ready for Implementation
- **Author:** Spec Author 02
- **Parent Plan:** [.ai-memory/plans/pending/66-fix-gitmap-pa-duplicate-repos.md](../../../.ai-memory/plans/pending/66-fix-gitmap-pa-duplicate-repos.md)
- **Target Subtasks:**
  - [Subtask 66.2: Pull Pipeline Deduplication, Ignore Auditing, Render Deduplication & Worker Safety](../../../.ai-memory/plans/subtasks/66-fix-gitmap-pa-duplicate-repos/02-pull-pipeline-deduplication.md)

---

## 1. Executive Summary & Problem Analysis

### 1.1 Context & Incident Description
When executing batch pull operations (`gitmap pa` / `gitmap pull --all`), users on Windows systems observed high rates of duplicated repository processing, misleading progress tallies, and unexpected git fast-forward failures:

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
```

In the target workspace, there were only **78 physical repositories**. However, GitMap reported **139 repositories pulled**, 37 failures citing `fatal: Cannot fast-forward to multiple branches.`, and 32 repositories with `.gitignore` issues.

### 1.2 Root Cause Analysis
1. **Unsanitized Pipeline Ingestion:**
   Even though database migrations address primary key collation in `gitmap.db`, `loadAllRecordsDB()` and `resolvePullTargets()` in `cli/cmdpull/` ingested raw records directly into the pull queue. Any existing case discrepancy (such as `D:\work\repo` vs `d:\work\repo` or multiple database entries for identical repo slugs) was fed straight into the concurrency worker pool.
2. **Concurrent Worker Collision on the Same Working Tree:**
   When both `D:\work\repo` and `d:\work\repo` were enqueued simultaneously, two independent worker goroutines in `cli/cmdpull/` spawned `git pull --progress --ff-only --autostash` on the exact same `.git` directory at the same time. The concurrent fetch operations wrote duplicate branch references into `.git/FETCH_HEAD`, causing Git to emit `fatal: Cannot fast-forward to multiple branches.` and abort the fast-forward.
3. **Multiplied Ignore Scanning:**
   The asynchronous background ignore scanner in `cli/cmdpull/pull_concurrency.go` received the same un-deduplicated record slice. Consequently, it audited the same physical directories twice, generating duplicate `IgnoreRepoIssue` entries and reporting inflated issue tallies.
4. **Duplicate Terminal Summary Rendering:**
   `RenderConciseActiveResultsTo` in `cli/cmdpull/pull_efficient_render.go` iterated through raw `states []*PullRepoState`. Multiple worker states for the same repository accumulated in the `failed`, `dirty`, and `updated` slices, inflating `len(failed)` to 37 and printing duplicate bullet rows for identical repositories.
5. **Uncaught Git Error in Worker Safety:**
   `isDivergedOutput` in `cli/cloner/safe_pull.go` only checked for `"Not possible to fast-forward"`, `"diverged"`, and `"non-fast-forward"`. It did not recognize `"Cannot fast-forward to multiple branches."`, preventing auto-merge recovery or diagnostic retry from activating.

### 1.3 Architectural Solution
This component specification establishes an end-to-end in-memory defense-in-depth architecture across four critical layers:
1. **In-Memory Target Deduplication:** Normalize repository paths canonically (`filepath.Clean`, forward slashes, lowercasing Windows drive/path) and deduplicate by canonical path and slug before workers are launched.
2. **Background Ignore Scan Deduplication:** Enforce uniqueness in `scanColdRecordsSequentially` and `collectAndRemediateIgnoreIssues` so each repository directory is inspected and reported at most once.
3. **Summary Render State Deduplication:** Deduplicate `PullRepoState` slices in `pull_efficient_render.go` so `len(failed)`, `len(dirty)`, and `len(updated)` represent physical unique repositories.
4. **Worker Concurrency Safety:** Update `cli/cloner/safe_pull.go` to recognize `Cannot fast-forward to multiple branches.` as a diverged condition, enable retry backoff, and enhance diagnostic hints in `cli/cloner/pulldiag.go`.

---

## 2. Component 1: In-Memory Target Deduplication Pipeline

### 2.1 File Location & Responsibilities
- **Primary Source:** [pull.go](file:///d:/work/gitmap/cli/cmdpull/pull.go)
- **Helpers:** [helpers.go](file:///d:/work/gitmap/cli/cmdpull/helpers.go)

### 2.2 Canonical Path Normalization
Windows paths can vary by drive letter casing (`D:\` vs `d:\`), slash direction (`\` vs `/`), trailing slashes (`D:\work\` vs `D:\work`), and dot segments (`D:\work\.\repo`). Canonicalization maps any path variant of a physical repository to a single deterministic key.

#### Implementation Contract
```go
// CanonicalRepoPathKey produces a normalized, lowercase forward-slash path key.
// It resolves relative segments, eliminates trailing slashes, and handles Windows drive case-insensitivity.
func CanonicalRepoPathKey(path string) string {
	if strings.TrimSpace(path) == "" {
		return ""
	}
	clean := filepath.Clean(path)
	slashed := filepath.ToSlash(clean)
	return strings.ToLower(slashed)
}
```

#### Normalization Examples
| Raw Input Path | `filepath.Clean` | `filepath.ToSlash` | `CanonicalRepoPathKey` |
| :--- | :--- | :--- | :--- |
| `D:\work\gitmap` | `D:\work\gitmap` | `D:/work/gitmap` | `d:/work/gitmap` |
| `d:\work\gitmap\` | `d:\work\gitmap` | `d:/work/gitmap` | `d:/work/gitmap` |
| `D:/work/./gitmap` | `D:\work\gitmap` | `D:/work/gitmap` | `d:/work/gitmap` |
| `d:/WORK/gitmap` | `d:\WORK\gitmap` | `d:/WORK/gitmap` | `d:/work/gitmap` |

### 2.3 Record Deduplication Algorithm
Records must be deduplicated across two dimensions:
1. **Canonical Path:** Avoid pulling the same physical directory on disk.
2. **Repository Slug:** Avoid pulling identical repo entities registered under differing aliases or paths.

#### Algorithm Specification
```go
// deduplicatePullRecords filters duplicates from a slice of ScanRecord.
// It ensures each repository canonical path and slug is enqueued exactly once.
func deduplicatePullRecords(records []model.ScanRecord) []model.ScanRecord {
	if len(records) <= 1 {
		return records
	}
	seenPath := make(map[string]bool, len(records))
	seenSlug := make(map[string]bool, len(records))
	unique := make([]model.ScanRecord, 0, len(records))

	for _, rec := range records {
		pathKey := CanonicalRepoPathKey(rec.AbsolutePath)
		slugKey := strings.ToLower(strings.TrimSpace(rec.Slug))
		if slugKey == "" {
			slugKey = strings.ToLower(strings.TrimSpace(rec.RepoName))
		}

		if isRecordDuplicate(pathKey, slugKey, seenPath, seenSlug) {
			continue
		}
		markRecordSeen(pathKey, slugKey, seenPath, seenSlug)
		unique = append(unique, rec)
	}

	return unique
}

func isRecordDuplicate(pathKey, slugKey string, seenPath, seenSlug map[string]bool) bool {
	if pathKey != "" && seenPath[pathKey] {
		return true
	}
	if slugKey != "" && seenSlug[slugKey] {
		return true
	}
	return false
}

func markRecordSeen(pathKey, slugKey string, seenPath, seenSlug map[string]bool) {
	if pathKey != "" {
		seenPath[pathKey] = true
	}
	if slugKey != "" {
		seenSlug[slugKey] = true
	}
}
```

### 2.4 Integration Points in `cmdpull`
The deduplication filter must be injected into all target ingestion pathways:

1. **`loadAllRecordsDB()` in `cli/cmdpull/helpers.go`:**
   ```go
   func loadAllRecordsDB() []model.ScanRecord {
       if LoadAllRecordsDBFn != nil {
           return deduplicatePullRecords(LoadAllRecordsDBFn())
       }
       return nil
   }
   ```
2. **`resolvePullTargets()` in `cli/cmdpull/pull.go`:**
   ```go
   func resolvePullTargets(slug, groupName string, all bool) []model.ScanRecord {
       raw := resolveRawPullTargets(slug, groupName, all)
       return deduplicatePullRecords(raw)
   }

   func resolveRawPullTargets(slug, groupName string, all bool) []model.ScanRecord {
       if HasAlias() {
           return resolveAliasRecord()
       }
       if len(groupName) > 0 {
           return loadRecordsByGroup(groupName)
       }
       if all {
           return loadAllRecordsDB()
       }
       return resolveSlugTarget(slug)
   }
   ```
3. **`findChildrenOfCWD()` in `cli/cmdpull/pull.go`:**
   Ensure children matched from `loadAllRecordsDB()` are filtered through `deduplicatePullRecords`.

---

## 3. Component 2: Background `.gitignore` Scanner Deduplication

### 3.1 File Location & Responsibilities
- **Concurrency Orchestration:** [pull_concurrency.go](file:///d:/work/gitmap/cli/cmdpull/pull_concurrency.go)
- **Issue Aggregation & Remediation:** [pull.go](file:///d:/work/gitmap/cli/cmdpull/pull.go)

### 3.2 Problem in Existing Scanner
When `StartThrottledAsyncIgnoreScan` executes, background worker thread scans cold records sequentially via `scanColdRecordsSequentially(records)`. If `records` contains duplicates, `inspectAndCacheRepoIgnore` runs multiple times for the same directory, producing duplicate entries in `issues []IgnoreRepoIssue`.

At the end of the pull batch:
`printIgnoreIssuesReport(issues)` prints `Detected .gitignore issues in %d repository(ies)`, which reflected 32 repositories instead of the true count of affected unique repositories.

### 3.3 Deduplication in `pull_concurrency.go`
In `cli/cmdpull/pull_concurrency.go`:

```go
func scanColdRecordsSequentially(records []model.ScanRecord) []IgnoreRepoIssue {
	var issues []IgnoreRepoIssue
	seenPath := make(map[string]bool, len(records))

	for i, rec := range records {
		pathKey := CanonicalRepoPathKey(rec.AbsolutePath)
		if pathKey != "" && seenPath[pathKey] {
			continue
		}
		if pathKey != "" {
			seenPath[pathKey] = true
		}

		yieldBetweenIgnoreAudits(i)
		issue := inspectAndCacheRepoIgnore(rec)
		if issue.HasIssues() {
			issues = append(issues, issue)
		}
	}
	return sortIgnoreIssues(issues)
}
```

### 3.4 Defensive Deduplication in `pull.go`
In `cli/cmdpull/pull.go`:

```go
// DeduplicateIgnoreIssues eliminates duplicate issues by canonical repo path.
func DeduplicateIgnoreIssues(issues []IgnoreRepoIssue) []IgnoreRepoIssue {
	if len(issues) <= 1 {
		return issues
	}
	seen := make(map[string]bool, len(issues))
	unique := make([]IgnoreRepoIssue, 0, len(issues))

	for _, issue := range issues {
		key := CanonicalRepoPathKey(issue.RepoPath)
		if key == "" {
			key = strings.ToLower(strings.TrimSpace(issue.RepoName))
		}
		if seen[key] {
			continue
		}
		seen[key] = true
		unique = append(unique, issue)
	}
	return unique
}

func sortIgnoreIssues(issues []IgnoreRepoIssue) []IgnoreRepoIssue {
	deduped := DeduplicateIgnoreIssues(issues)
	sorted := make([]IgnoreRepoIssue, len(deduped))
	copy(sorted, deduped)
	sort.Slice(sorted, func(i, j int) bool {
		return strings.ToLower(sorted[i].RepoName) < strings.ToLower(sorted[j].RepoName)
	})
	return sorted
}

func collectAndRemediateIgnoreIssues(handle *IgnoreScanHandle, opts pullOptions) {
	if handle == nil {
		return
	}
	issues := handle.Collect()
	if len(issues) == 0 {
		return
	}
	deduped := DeduplicateIgnoreIssues(issues)
	isAutoYes := opts.yes || opts.autoFix
	handleIgnoreRemediation(deduped, isAutoYes, opts.isJSON)
}
```

---

## 4. Component 3: Summary Render Deduplication

### 4.1 File Location & Responsibilities
- **Summary Render Pipeline:** [pull_efficient_render.go](file:///d:/work/gitmap/cli/cmdpull/pull_efficient_render.go)
- **Batch Output Coordinator:** [pull.go](file:///d:/work/gitmap/cli/cmdpull/pull.go)

### 4.2 State Deduplication Logic
The function `RenderConciseActiveResultsTo(w io.Writer, states []*PullRepoState, allRecords ...[]model.ScanRecord)` groups states into `updated`, `dirty`, and `failed`.

To ensure that each physical repository is rendered exactly once and that `len(failed)`, `len(dirty)`, and `len(updated)` represent unique counts:

```go
// DeduplicateRepoStates ensures each repository appears at most once in output states.
// When collisions occur, error/failure states take precedence over up-to-date states.
func DeduplicateRepoStates(states []*PullRepoState) []*PullRepoState {
	if len(states) <= 1 {
		return states
	}
	seen := make(map[string]int, len(states))
	unique := make([]*PullRepoState, 0, len(states))

	for _, s := range states {
		if s == nil {
			continue
		}
		key := CanonicalRepoPathKey(s.RepoPath)
		if key == "" {
			key = strings.ToLower(strings.TrimSpace(s.RepoName))
		}

		if idx, exists := seen[key]; exists {
			mergeRepoStateIfPrioritized(unique, idx, s)
			continue
		}
		seen[key] = len(unique)
		unique = append(unique, s)
	}

	return unique
}

func mergeRepoStateIfPrioritized(unique []*PullRepoState, existingIdx int, incoming *PullRepoState) {
	existing := unique[existingIdx]
	// Prefer error/dirty state over clean/up-to-date state for accurate diagnostic display
	isIncomingActionable := incoming.IsDirty || incoming.ErrorMsg != "" || incoming.Changes == "failed"
	isExistingPassive := !existing.IsDirty && existing.ErrorMsg == "" && existing.Changes != "failed"

	if isIncomingActionable && isExistingPassive {
		unique[existingIdx] = incoming
	}
}
```

### 4.3 Integration in `RenderConciseActiveResultsTo`
In `cli/cmdpull/pull_efficient_render.go`:

```go
func RenderConciseActiveResultsTo(w io.Writer, states []*PullRepoState, allRecords ...[]model.ScanRecord) {
	dedupedStates := DeduplicateRepoStates(states)
	colWidth := ResolveConciseRepoColWidth(dedupedStates, allRecords...)
	var updated, dirty, failed []*PullRepoState

	for _, s := range dedupedStates {
		if s.IsDirty || s.Changes == "dirty" {
			dirty = append(dirty, s)
		} else if isFailedRepoState(s) {
			failed = append(failed, s)
		} else if isUpdatedRepoState(s) {
			updated = append(updated, s)
		}
	}
    // ... renderUpdatedGroup, renderDirtyGroup, renderFailedGroup ...
}

func isFailedRepoState(s *PullRepoState) bool {
	return s.ErrorMsg != "" || s.Step == PullStepTypeError || s.Step == PullStepTypeConflict || s.Changes == "failed"
}
```

### 4.4 Consistency with `pull.go` Counters
In `cli/cmdpull/pull.go`:
- `runPullBatchExecution` returns `sortedStates`:
  ```go
  func runPullBatchExecution(records []model.ScanRecord, opts pullOptions) (*PullProgressBar, []*PullRepoState, time.Duration) {
      // ...
      dedupedStates := DeduplicateRepoStates(bar.States())
      return bar, sortStatesAlphabetically(dedupedStates), dur
  }
  ```
- This ensures `countActiveStates(sortedStates)` and `printPullAllFastSummary(len(sortedStates), ...)` output the true physical repository count (e.g. 78 repositories instead of 139).

---

## 5. Component 4: Worker Concurrency Safety & Fast-Forward Guard

### 5.1 File Location & Responsibilities
- **Safe Pull Execution:** [safe_pull.go](file:///d:/work/gitmap/cli/cloner/safe_pull.go)
- **Pull Diagnostics:** [pulldiag.go](file:///d:/work/gitmap/cli/cloner/pulldiag.go)

### 5.2 Handling `Cannot fast-forward to multiple branches.`
When Git executes `git pull --progress --ff-only --autostash` and `.git/FETCH_HEAD` contains multiple entries, Git outputs:
```
fatal: Cannot fast-forward to multiple branches.
```
In `cli/cloner/safe_pull.go`, `runGitPullWithProgress` checks `isDivergedOutput(out)`. Currently, this function returns `false` for multi-branch output, treating it as an unrecoverable failure instead of attempting auto-merge or diagnosis.

#### Enhanced `isDivergedOutput`
```go
func isDivergedOutput(output string) bool {
	return strings.Contains(output, "Not possible to fast-forward") ||
		strings.Contains(output, "diverged") ||
		strings.Contains(output, "non-fast-forward") ||
		strings.Contains(output, "Cannot fast-forward to multiple branches") ||
		strings.Contains(output, "cannot fast-forward to multiple branches")
}
```

### 5.3 Diagnostic Hinting in `pulldiag.go`
In `cli/cloner/pulldiag.go`, update `collectDiagnosisHints`:

```go
func collectDiagnosisHints(repoDir, output string) []string {
	hints := make([]string, 0, 5)
	if hasUnlinkFailure(output) {
		hints = append(hints, "file lock/read-only attribute blocked replacing old files")
	}
	if hasUnmergedFailure(output) {
		hints = append(hints, "unresolved merge conflict detected; run 'gitmap fix-git' or 'git merge --abort'")
	}
	if hasUntrackedOverwriteFailure(output) {
		hints = append(hints, "untracked files conflict with incoming commits; run 'gitmap fix-git' to backup & pull")
	}
	if hasMultipleBranchesFailure(output) {
		hints = append(hints, "multiple upstream branches or concurrent fetch conflict in FETCH_HEAD; run pull targeting specific branch")
	}
	if hasPathLengthRisk(repoDir, output) {
		hints = append(hints, "Windows path length risk detected; use a shorter base path like C:\\src")
	}
	if strings.Contains(strings.ToLower(repoDir), "onedrive") {
		hints = append(hints, "repo is under a synced folder (OneDrive), which often locks files")
	}
	return hints
}

func hasMultipleBranchesFailure(output string) bool {
	lower := strings.ToLower(output)
	return strings.Contains(lower, "cannot fast-forward to multiple branches")
}
```

### 5.4 Concurrency Isolation Guard
By pairing:
1. **Upstream In-Memory Deduplication** (prevents two goroutines from ever running `safe_pull` on the same directory simultaneously), and
2. **`isDivergedOutput` + `attemptAutoMergePull` fallback** (recovers repositories that have pre-existing multi-branch FETCH_HEAD pollution),

Worker threads are completely isolated from cross-process race conditions.

---

## 6. Coding Guidelines & Quality Compliance Matrix

| Rule | Requirement | Component Implementation Strategy |
| :--- | :--- | :--- |
| **R1: Function Length** | Max 15 lines per function | Decompose deduplication logic into small helpers (`isRecordDuplicate`, `markRecordSeen`, `isFailedRepoState`, `hasMultipleBranchesFailure`). |
| **R2: Positive Booleans** | Positive prefix (`is`, `has`), no negative booleans | Use `isRecordDuplicate`, `hasIssues`, `isIncomingActionable`, `isExistingPassive`, avoiding `noDuplicates` or `isNotUnique`. |
| **R3: Error Management** | Wrap errors with `apperror` | Wrap any path resolution or file inspect errors using structured `apperror.WrapSimple`. |
| **R4: Immutability** | Avoid mutating shared slices in-place | Create new slices `unique := make([]model.ScanRecord, 0, len(records))` and return pure deduplicated arrays. |
| **R5: No Magic Literals** | Use existing constants | Reuse `constants.SafePullRetryAttempts`, `constants.ColorRed`, etc. |

---

## 7. Verification & Acceptance Criteria

### 7.1 Automated Unit Tests
1. **`CanonicalRepoPathKey` Tests:**
   - Verify `D:\work\gitmap`, `d:\work\gitmap`, `D:/work/gitmap`, and `d:\work\gitmap\` all resolve to identical string `"d:/work/gitmap"`.
2. **`deduplicatePullRecords` Tests:**
   - Given a slice of 4 records containing 2 duplicates with varying path case and matching slugs, output slice has length exactly 2.
3. **`DeduplicateIgnoreIssues` Tests:**
   - Given 2 `IgnoreRepoIssue` structs pointing to `D:\work\repo` and `d:\work\repo`, output slice has length 1.
4. **`DeduplicateRepoStates` Tests:**
   - Given multiple states for the same repo (one failed, one up-to-date), output retains the failed state for diagnostic visibility, length is 1.
5. **`isDivergedOutput` Tests:**
   - Output containing `fatal: Cannot fast-forward to multiple branches.` returns `true`.

### 7.2 Manual & End-to-End Verification
- Run `gitmap pa` on workspace with duplicate SQLite records.
- Expected outcome:
  - Pulled count strictly matches unique physical repository count (78 repos).
  - Failed repository list contains 0 duplicates.
  - Zero repositories fail with `Cannot fast-forward to multiple branches.`.
  - `.gitignore` issue detection reports unique physical repositories without repeats.
