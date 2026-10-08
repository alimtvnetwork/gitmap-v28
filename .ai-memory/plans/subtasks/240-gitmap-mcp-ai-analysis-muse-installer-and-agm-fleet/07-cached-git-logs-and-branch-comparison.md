# Subtask 07: High-Speed Cached Git Logs & Branch Comparison

> **Subtask ID:** Subtask-07  
> **Parent Plan:** `.ai-memory/plans/240-gitmap-mcp-ai-analysis-muse-installer-and-agm-fleet.md`  
> **Target Subsystems:** `cli/cmdlog/`, `cli/cmdgit/`, `cli/cmd/`  
> **Owned Files:**  
> - `cli/cmdlog/log_types.go`  
> - `cli/cmdlog/log_db.go`  
> - `cli/cmdlog/log_sync.go`  
> - `cli/cmdlog/log_cmd.go`  
> - `cli/cmdlog/log_test.go`  
> - `cli/cmdgit/branch_compare.go`  
> - `cli/cmdgit/branch_compare_test.go`  
> - `cli/cmd/branch.go`  
> - `cli/cmd/roottooling.go`  

---

## 1. Concrete Objectives

1. **Split-DB Cached Commit Storage:**
   - Establish `.gitmap/data/repodb/git_logs.db` storing structured commit history in table `CachedCommits` and tracking tip hashes in `RefPointers`.
   - Store hash, parent hashes, author, committer, dates, subject, body, changed file count, insertions, deletions, and JSON stats.
2. **Sub-5ms Query Latency & Incremental Synchronization:**
   - Implement `SyncCommitCache(repoPath string, isForced bool)`:
     - Read `.git/HEAD` and current branch ref tip hash in < 1ms.
     - If the tip hash matches `RefPointers.TipCommitHash`, serve cached commits directly from SQLite (benchmark < 5ms).
     - If ref has advanced, execute incremental delta extraction (`git log <cached_tip>..HEAD --format=...`), parse records, bulk-insert into `CachedCommits`, and update the ref pointer.
3. **High-Speed Cached Log CLI (`gitmap log`):**
   - Provide command `gitmap log` displaying recent commits in a clean ANSI format with relative timestamps, author badges, and commit hashes.
   - Support flags: `--limit` / `-n` (default 20), `--author`, `--since`, `--stat`, `--graph`, `--refresh` / `-r`, and `--json`.
4. **Fast Branch Comparison Subsystem (`gitmap diff-branch`, `gitmap branch compare`):**
   - Compare two branches (e.g., `main` vs `feature/branch`, or current branch vs `origin/main`):
     - Compute divergence merge-base hash (`git merge-base branchA branchB`).
     - Query commits ahead in branch A and commits behind in branch B.
     - Extract file difference summary: counts of added, modified, deleted files, total insertions, and deletions.
   - If both branch tips are cached, calculate divergence sets in SQLite without triggering Git subprocesses.
5. **Command Dispatcher Wiring:**
   - Register `gitmap log` in `cli/cmd/roottooling.go`.
   - Register `gitmap diff-branch` in `cli/cmd/roottooling.go`.
   - Add `compare` subcommand to `gitmap branch` in `cli/cmd/branch.go`.

---

## 2. Core Domain Types & Structs

```go
package cmdlog

import (
	"time"
)

// CachedCommit represents a commit row in SQLite.
type CachedCommit struct {
	CommitHash        string    `json:"commit_hash"`
	ParentHashes      []string  `json:"parent_hashes"`
	AuthorName        string    `json:"author_name"`
	AuthorEmail       string    `json:"author_email"`
	AuthorDate        time.Time `json:"author_date"`
	CommitterName     string    `json:"committer_name"`
	CommitterEmail    string    `json:"committer_email"`
	CommitDate        time.Time `json:"commit_date"`
	Subject           string    `json:"subject"`
	Body              string    `json:"body,omitempty"`
	FilesChangedCount int       `json:"files_changed_count"`
	LinesInserted     int       `json:"lines_inserted"`
	LinesDeleted      int       `json:"lines_deleted"`
	StatsJSON         string    `json:"stats_json,omitempty"`
	CachedAt          time.Time `json:"cached_at"`
}

// LogQueryOptions holds CLI query parameters for gitmap log.
type LogQueryOptions struct {
	Limit      int       `json:"limit"`
	Author     string    `json:"author,omitempty"`
	Since      time.Time `json:"since,omitempty"`
	FileFilter string    `json:"file_filter,omitempty"`
	IsRefresh  bool      `json:"is_refresh"`
	IsStat     bool      `json:"is_stat"`
	IsGraph    bool      `json:"is_graph"`
	IsJSON     bool      `json:"is_json"`
}

// BranchDiffResult encapsulates divergence telemetry between two branch references.
type BranchDiffResult struct {
	BranchA         string         `json:"branch_a"`
	BranchB         string         `json:"branch_b"`
	MergeBaseHash   string         `json:"merge_base_hash"`
	AheadCount      int            `json:"ahead_count"`
	BehindCount     int            `json:"behind_count"`
	AheadCommits    []CachedCommit `json:"ahead_commits"`
	BehindCommits   []CachedCommit `json:"behind_commits"`
	FilesAdded      int            `json:"files_added"`
	FilesModified   int            `json:"files_modified"`
	FilesDeleted    int            `json:"files_deleted"`
	TotalInsertions int            `json:"total_insertions"`
	TotalDeletions  int            `json:"total_deletions"`
	DurationMs      int64          `json:"duration_ms"`
}
```

---

## 3. Implementation Checklist

- [ ] **1. SQLite Split-DB Schema & Connection (`cli/cmdlog/log_db.go`):**
  - Implement `OpenRepoLogDB(repoPath string) (*sql.DB, error)`.
  - Create table `CachedCommits` with indexes on `CommitDate` and `AuthorName`.
  - Create table `RefPointers` mapping branch refs to tip commit hashes.
  - Create table `BranchDiffCache` to memoize expensive merge-base computations.
- [ ] **2. Incremental Sync Engine (`cli/cmdlog/log_sync.go`):**
  - Implement `GetRepoRefTip(repoPath string, refName string) (string, error)`.
  - Implement `SyncCommitCache(repoPath string, isForced bool) error`.
  - Check `RefPointers` for existing tip. If match found and not forced, return immediately.
  - If no tip exists (cold cache), run initial batch import: `git log -n 1000 --format="%H%x1f%P%x1f%an%x1f%ae%x1f%aI%x1f%s%x1e"`.
  - If tip exists but ref moved, run incremental sync: `git log <cached_tip>..HEAD --format=...`.
  - Bulk-insert commits inside an atomic SQLite transaction.
  - Update `RefPointers` with the new tip commit hash.
- [ ] **3. Log Query & Formatting Engine (`cli/cmdlog/log_cmd.go`):**
  - Implement `RunLogCommand(args []string) error`.
  - Parse CLI flags (`-n`, `--limit`, `--author`, `--since`, `--stat`, `--refresh`, `--json`).
  - Query `CachedCommits` using parameterized SQL queries.
  - Render compact 1-line ANSI output per commit:
    `[<short_hash>] <subject> (<relative_time>) - <author>`
  - If `--stat` is enabled, render lines inserted and deleted badges.
  - If `--json` is enabled, emit structured JSON.
- [ ] **4. Branch Comparison Engine (`cli/cmdgit/branch_compare.go`):**
  - Implement `RunDiffBranchCommand(args []string) error`.
  - Resolve branch names: if branch A omitted, use current branch; if branch B omitted, resolve repository default branch (`main`).
  - Compute merge base commit via `git merge-base <branchA> <branchB>`.
  - Extract unique commits in branch A and branch B using SQLite graph queries or `git log`.
  - Compute file diff summary via `git diff --shortstat <mergeBase>..<branchA>`.
  - Format ANSI side-by-side comparison summary.
- [ ] **5. CLI Wiring & Subcommand Integration (`cli/cmd/roottooling.go`, `cli/cmd/branch.go`):**
  - Register `log` and `diff-branch` in `cli/cmd/roottooling.go`.
  - Add `compare` subcommand in `cli/cmd/branch.go` delegating to `cmdgit.RunDiffBranchCommand`.
- [ ] **6. Performance & Integrity Tests (`cli/cmdlog/log_test.go`, `cli/cmdgit/branch_compare_test.go`):**
  - Benchmark query latency: assert cache hit returns within 5ms.
  - Test incremental sync on commit generation.
  - Test branch comparison ahead/behind counts against known commits.

---

## 4. Acceptance Criteria

1. Running `gitmap log` returns the most recent 20 commits in under 15ms on warm cache.
2. Making a new commit advances HEAD and triggers an incremental delta sync without re-parsing full history.
3. Running `gitmap diff-branch` accurately reports ahead/behind commit counts and file mutation statistics.
4. Running `gitmap branch compare feature/x main` provides identical results to `gitmap diff-branch feature/x main`.
5. All outputs strictly adhere to positive boolean naming conventions and relative path formatting.

---

## 5. Verification Commands

```bash
# Run unit tests
go test -v ./cli/cmdlog/...
go test -v ./cli/cmdgit/...

# Verify compilation
go build -v ./cli/...

# Benchmark cached log performance
gitmap log -n 20
gitmap log --refresh -n 20

# Test branch comparison
gitmap diff-branch
```
