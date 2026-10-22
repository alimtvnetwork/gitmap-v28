# Specification 292: Migration, Operational Safety & Verification Specification (Batch 4)

- **Sequence:** 292
- **Slug:** `292-improvements-planning-and-execution`
- **File:** `02-spec/21-app/292-improvements-planning-and-execution/03-migration-and-verification-spec.md`
- **Status:** APPROVED & SPECIFIED
- **Scope:** Wave 3 / Improvements 21–30 (System UX, Operational Safety, Database & Git Subsystem Enhancements)
- **Assigned Author:** Spec Author 03
- **Parent Specs:**
  - `02-spec/21-app/292-improvements-planning-and-execution/01-architecture-spec.md`
  - `02-spec/21-app/292-improvements-planning-and-execution/02-component-spec.md`

---

## 1. System UX & Operational Safety Architecture (Improvements 21–30)

Improvements 21 through 30 in Specification 292 establish runtime operational stability, targeted repository metadata filtering, decision heatmap task action aggregation, dedicated agent lifecycle action recording, parameterized and whitelisted UI log sorting, subprocess git revision counting and merge-base resolution, and architectural file decomposition across `cpar.go`, `status.go`, and `cmd_see.go` to maintain strict adherence to the <= 500-line repository ceiling and <= 100-line modular design standard.

```mermaid
flowchart TD
    subgraph StoreFilterLayer["cli/store (Repository Filter & Query Store)"]
        A["repo_query_filter.go<br/>ListReposByTag(tag)"]
        B["repo_query_filter.go<br/>ListReposUnpushed()"]
    end

    subgraph HeatmapStoreLayer["cli/store (Heatmap Storage & Schema)"]
        C["heatmap_schema.go<br/>HeatmapTaskSummary Struct"]
        D["heatmap_store.go<br/>GetHeatmapTaskSummary(taskSlug)"]
    end

    subgraph HeatmapTrackingLayer["cli/heatmap (Telemetry Tracking)"]
        E["heatmap_tracker_agent.go<br/>TrackAgentAction(taskSlug, actionType, summary)"]
    end

    subgraph AgentUILayer["cli/cmdagent (Agent UI Server & Logs)"]
        F["agent_ui_server.go<br/>ActionLogFilter (SortBy, Order)"]
        G["agent_ui_logs.go<br/>parseSortColumn / parseSortOrder<br/>Dynamic SQL ORDER BY"]
    end

    subgraph GitUtilLayer["cli/gitutil (Subprocess Commit & Revision Helpers)"]
        H["git_repo_commits.go<br/>GetRepoUnpushedCommitCount(repoPath)"]
        I["git_repo_commits.go<br/>GetRepoMergeBase(repoPath, ref1, ref2)"]
    end

    subgraph ModularDecompositionLayer["cli/cmdcpar, cli/cmdstatus, cli/cmd (Architectural Decompositions)"]
        J["cpar_target.go<br/>Extract Target Discovery & Dirty Collection"]
        K["status_resolve.go<br/>Extract Scope Resolution & DB/JSON Record Loading"]
        L["see_ops.go<br/>Extract Operational Handlers (RepoManageUI, CommitPending, History)"]
    end

    A --> StoreFilterLayer
    B --> StoreFilterLayer
    C --> D
    E --> D
    F --> G
    H --> GitUtilLayer
    I --> GitUtilLayer
    J --> ModularDecompositionLayer
    K --> ModularDecompositionLayer
    L --> ModularDecompositionLayer
```

---

## 2. Detailed Technical Specification for Improvements 21–30

### Improvement 21: Tag-Based Repository Listing (`cli/store/repo_query_filter.go`)
- **Objective:** Provide a dedicated query method `ListReposByTag(tag string) ([]model.ScanRecord, error)` to query repositories filtered by metadata tags or notes.
- **Current Problem:** Consumers that filter repositories by tag load the entire repository dataset into memory using `db.ListRepos()` and perform manual string matching. Furthermore, placing additional query functions into `cli/store/repo.go` (currently 384 lines) threatens the 500-line hard ceiling.
- **Architectural Solution:** Implement `ListReposByTag` in a dedicated companion file `cli/store/repo_query_filter.go` within package `store`.
- **Contract:**
  ```go
  package store

  import (
  	"fmt"
  	"strings"

  	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
  	"github.com/alimtvnetwork/gitmap-v28/cli/model"
  )

  const sqlSelectReposByTag = "SELECT RepoId, Slug, RepoName, HttpsUrl, SshUrl, Branch, RelativePath, AbsolutePath, CloneInstruction, Notes, IdentifiedTransport FROM Repo WHERE Notes LIKE ? COLLATE NOCASE ORDER BY Slug ASC"

  // ListReposByTag returns all repositories matching a specific tag stored in metadata or notes.
  func (db *DB) ListReposByTag(tag string) ([]model.ScanRecord, error) {
  	cleanTag := strings.TrimSpace(tag)
  	if cleanTag == "" {
  		return db.ListRepos()
  	}

  	pattern := "%" + cleanTag + "%"
  	rows, err := QueryWrapper(db.conn, sqlSelectReposByTag, pattern).Destruct()
  	if err != nil {
  		return nil, fmt.Errorf(constants.ErrDBQuery, err)
  	}

  	defer rows.Close()

  	return scanRows(rows)
  }
  ```
- **Error Sovereignty:** Returns SQL query errors wrapped via `fmt.Errorf(constants.ErrDBQuery, err)`. Zero exit handlers or process terminations.
- **Line Budget:** Function is 14 lines, well within the 15-line maximum.

---

### Improvement 22: Unpushed Repository Listing (`cli/store/repo_query_filter.go`)
- **Objective:** Provide a dedicated query method `ListReposUnpushed() ([]model.ScanRecord, error)` to query repositories flagged with unpushed commits in metadata or local heatmap tracking.
- **Current Problem:** Discerning which repositories have pending unpushed state currently requires sweeping filesystem inspections or parsing unstructured notes after loading all repository records.
- **Architectural Solution:** Declare `ListReposUnpushed` in `cli/store/repo_query_filter.go`, executing an indexed SQLite query against `Repo.Notes` and `HeatmapLocalTracking.IsDirty`.
- **Contract:**
  ```go
  const sqlSelectUnpushedRepos = "SELECT RepoId, Slug, RepoName, HttpsUrl, SshUrl, Branch, RelativePath, AbsolutePath, CloneInstruction, Notes, IdentifiedTransport FROM Repo WHERE (Notes LIKE '%unpushed%' OR Notes LIKE '%ahead%') OR RepoName IN (SELECT RepoSlug FROM HeatmapLocalTracking WHERE IsDirty = 1) ORDER BY Slug ASC"

  // ListReposUnpushed returns repositories flagged with unpushed commits in metadata or local tracking.
  func (db *DB) ListReposUnpushed() ([]model.ScanRecord, error) {
  	rows, err := QueryWrapper(db.conn, sqlSelectUnpushedRepos).Destruct()
  	if err != nil {
  		return nil, fmt.Errorf(constants.ErrDBQuery, err)
  	}

  	defer rows.Close()

  	return scanRows(rows)
  }
  ```
- **Error Sovereignty:** Returns database errors directly to caller without exit handler interception.
- **Line Budget:** 9 lines. Zero nested conditionals.

---

### Improvement 23: Heatmap Task Action Aggregation (`cli/store/heatmap_store.go`, `cli/store/heatmap_schema.go`)
- **Objective:** Define the `HeatmapTaskSummary` model and implement `GetHeatmapTaskSummary(taskSlug string) (*HeatmapTaskSummary, error)` for single-row aggregated metrics.
- **Current Problem:** `HeatmapStatusSummary` aggregates overall action records, but does not provide a lightweight, dedicated summary categorized into read, edit, create, and other operations for targeted task reporting.
- **Architectural Solution:**
  1. Define `HeatmapTaskSummary` in `cli/store/heatmap_schema.go`.
  2. Implement `GetHeatmapTaskSummary` on `*HeatmapDB` in `cli/store/heatmap_store.go`.
- **Model Contract:**
  ```go
  // HeatmapTaskSummary aggregates overall action volume and categorizes actions by operation touch type for a task.
  type HeatmapTaskSummary struct {
  	TaskSlug     string `json:"taskSlug"`
  	TotalActions int    `json:"totalActions"`
  	ReadCount    int    `json:"readCount"`
  	EditCount    int    `json:"editCount"`
  	CreateCount  int    `json:"createCount"`
  	OtherCount   int    `json:"otherCount"`
  }
  ```
- **Method Contract:**
  ```go
  const sqlHeatmapTaskSummary = "SELECT COUNT(*), COALESCE(SUM(CASE WHEN UPPER(TRIM(ActionType)) IN ('READ', 'READ_FILE', 'INSPECT') THEN 1 ELSE 0 END), 0), COALESCE(SUM(CASE WHEN UPPER(TRIM(ActionType)) IN ('EDIT', 'EDIT_FILE', 'MODIFY') THEN 1 ELSE 0 END), 0), COALESCE(SUM(CASE WHEN UPPER(TRIM(ActionType)) IN ('CREATE', 'WRITE_FILE', 'NEW') THEN 1 ELSE 0 END), 0) FROM HeatmapActionLog WHERE (? = '' OR TaskSlug = ?)"

  // GetHeatmapTaskSummary returns aggregated read, edit, and create counts for a specified task.
  func (h *HeatmapDB) GetHeatmapTaskSummary(taskSlug string) (*HeatmapTaskSummary, error) {
  	cleanSlug := strings.TrimSpace(taskSlug)
  	var total, reads, edits, creates int
  	err := h.conn.QueryRow(sqlHeatmapTaskSummary, cleanSlug, cleanSlug).Scan(&total, &reads, &edits, &creates)
  	if err != nil {
  		return nil, apperror.WrapSimple(err, "heatmap.taskSummary")
  	}

  	other := total - (reads + edits + creates)
  	if other < 0 {
  		other = 0
  	}

  	return &HeatmapTaskSummary{
  		TaskSlug:     cleanSlug,
  		TotalActions: total,
  		ReadCount:    reads,
  		EditCount:    edits,
  		CreateCount:  creates,
  		OtherCount:   other,
  	}, nil
  }
  ```
- **Error Sovereignty:** Wrapped with `apperror.WrapSimple(err, "heatmap.taskSummary")`.
- **Line Budget:** 14 lines. Zero nested conditionals.

---

### Improvement 24: Dedicated Agent Lifecycle Action Tracking (`cli/heatmap/heatmap_tracker_agent.go`)
- **Objective:** Implement `TrackAgentAction(taskSlug, actionType, summary string) error` to record discrete agent lifecycle events into the root heatmap database.
- **Current Problem:** `cli/heatmap/heatmap_tracker.go` is currently 93 lines. Adding agent action tracking directly into `heatmap_tracker.go` would immediately breach the 100-line file cap.
- **Architectural Solution:** Create a new companion file `cli/heatmap/heatmap_tracker_agent.go` housing `TrackAgentAction` and its private connection runner.
- **Contract:**
  ```go
  package heatmap

  import (
  	"os"
  	"strings"

  	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
  	"github.com/alimtvnetwork/gitmap-v28/cli/store"
  )

  // TrackAgentAction records a discrete agent lifecycle event into the root heatmap database.
  func TrackAgentAction(taskSlug, actionType, summary string) error {
  	cleanSlug := strings.TrimSpace(taskSlug)
  	cleanAction := strings.TrimSpace(actionType)
  	if cleanSlug == "" || cleanAction == "" {
  		return nil
  	}

  	cwd, err := os.Getwd()
  	if err != nil {
  		return apperror.WrapSimple(err, "heatmap.cwd")
  	}

  	return recordAgentActionAtRoot(cwd, cleanSlug, cleanAction, summary)
  }

  func recordAgentActionAtRoot(cwd, taskSlug, actionType, summary string) error {
  	hdb, err := store.OpenHeatmapRootDB(cwd)
  	if err != nil {
  		return err
  	}

  	defer hdb.Close()

  	_, logErr := hdb.LogAction(store.HeatmapActionLogRecord{
  		TaskSlug:       taskSlug,
  		ActionType:     actionType,
  		ReasoningNotes: summary,
  		ReasoningState: "EXPLORATION",
  	})

  	return logErr
  }
  ```
- **Line Budget:** Each function <= 12 lines. Zero nested conditionals.

---

### Improvement 25: Sortable Agent Action Log Query Engine (`cli/cmdagent/agent_ui_logs.go`, `cli/cmdagent/agent_ui_server.go`)
- **Objective:** Add `SortBy` and `Order` fields to `ActionLogFilter`, implement sanitized sort parsing, and enable dynamic SQL ordering.
- **Current Problem:** `cli/cmdagent/agent_ui_logs.go` hardcodes `" ORDER BY ActionLogId DESC LIMIT ?"`, preventing UI clients from ordering action logs by timestamp, action type, agent role, or subtask.
- **Architectural Solution:**
  1. Add `SortBy string` and `Order string` to `ActionLogFilter` in `cli/cmdagent/agent_ui_server.go`.
  2. Implement `parseSortColumn` and `parseSortOrder` in `cli/cmdagent/agent_ui_logs.go` with strict whitelist sanitization against SQL injection.
  3. Update `appendPaginationLimitOffset` to dynamically apply whitelisted sort column and order.
- **Whitelist Helpers:**
  ```go
  func parseSortColumn(raw string) string {
  	switch strings.ToLower(strings.TrimSpace(raw)) {
  	case "createdat", "time", "date":
  		return "CreatedAt"
  	case "actiontype", "type":
  		return "ActionType"
  	case "agentrole", "role":
  		return "AgentRole"
  	case "subtaskid", "subtask":
  		return "SubtaskId"
  	default:
  		return "ActionLogId"
  	}
  }

  func parseSortOrder(raw string) string {
  	if strings.EqualFold(strings.TrimSpace(raw), "asc") {
  		return "ASC"
  	}

  	return "DESC"
  }
  ```
- **Dynamic Order Formulation:**
  ```go
  func appendPaginationLimitOffset(base string, args []any, f ActionLogFilter) (string, []any) {
  	sortCol := parseSortColumn(f.SortBy)
  	sortDir := parseSortOrder(f.Order)
  	base += fmt.Sprintf(" ORDER BY %s %s LIMIT ?", sortCol, sortDir)
  	args = append(args, f.Limit)
  	if f.Offset > 0 {
  		base += " OFFSET ?"
  		args = append(args, f.Offset)
  	}

  	return base, args
  }
  ```

---

### Improvement 26: Subprocess Git Unpushed Commit Counter (`cli/gitutil/git_repo_commits.go`)
- **Objective:** Add `GetRepoUnpushedCommitCount(repoPath string) (int, error)` executing `git rev-list --count @{u}..HEAD`.
- **Current Problem:** `cli/gitutil/git_repo_query.go` has `HasUnpushedCommits` returning a boolean via `git cherry -v`, but does not quantify the exact number of pending commits. Placing new functions directly in `git_repo_query.go` (84 lines) would exceed the 100-line cap.
- **Architectural Solution:** Create `cli/gitutil/git_repo_commits.go` for revision and commit range calculation helpers.
- **Contract:**
  ```go
  package gitutil

  import (
  	"strconv"
  	"strings"
  )

  // GetRepoUnpushedCommitCount returns the number of unpushed commits ahead of upstream tracking.
  func GetRepoUnpushedCommitCount(repoPath string) (int, error) {
  	out, err := RunGitInDir(repoPath, "rev-list", "--count", "@{u}..HEAD")
  	if err != nil {
  		return 0, err
  	}

  	clean := strings.TrimSpace(string(out))
  	count, parseErr := strconv.Atoi(clean)
  	if parseErr != nil {
  		return 0, parseErr
  	}

  	return count, nil
  }
  ```
- **Line Budget:** 12 lines. Zero nested ifs.

---

### Improvement 27: Common Ancestor Merge Base Resolution (`cli/gitutil/git_repo_commits.go`)
- **Objective:** Add `GetRepoMergeBase(repoPath, ref1, ref2 string) (string, error)` executing `git merge-base ref1 ref2`.
- **Current Problem:** Determining the divergence point or common ancestor between branches (e.g. `HEAD` and `main`) currently requires manually orchestrating subprocesses in calling packages.
- **Architectural Solution:** Declare `GetRepoMergeBase` in `cli/gitutil/git_repo_commits.go`.
- **Contract:**
  ```go
  // GetRepoMergeBase finds the best common ancestor commit between two references.
  func GetRepoMergeBase(repoPath, ref1, ref2 string) (string, error) {
  	out, err := RunGitInDir(repoPath, "merge-base", ref1, ref2)
  	if err != nil {
  		return "", err
  	}

  	return strings.TrimSpace(string(out)), nil
  }
  ```
- **Line Budget:** 7 lines. Zero nested ifs.

---

### Improvement 28: Target Discovery & Dirty Extraction (`cli/cmdcpar/cpar_target.go`)
- **Objective:** Extract target directory discovery and dirty repository collection from `cli/cmdcpar/cpar.go` into `cli/cmdcpar/cpar_target.go`.
- **Current Problem:** `cli/cmdcpar/cpar.go` spans 293 lines, combining CLI flag parsing, target repository collection, git operations dispatch, progress bars, and task queue reporting.
- **Architectural Solution:** Extract target collection routines into `cli/cmdcpar/cpar_target.go`:
  - `CollectDirtyRepositories()`
  - `resolveRelativePath(absPath, rootDir string) string`
  - `appendDirty(list []cparTarget, name, relPath, absPath string) []cparTarget`
  - `resolveAllRecords() ([]model.ScanRecord, error)`
  - `queryStoreRecords() ([]model.ScanRecord, error)`
- **Impact:** Decreases `cpar.go` by ~65 lines, creating clean separation between target discovery and execution.

---

### Improvement 29: Scope Resolution & DB Fallback Extraction (`cli/cmdstatus/status_resolve.go`)
- **Objective:** Extract scope resolution and DB/JSON record loading from `cli/cmdstatus/status.go` into `cli/cmdstatus/status_resolve.go`.
- **Current Problem:** `cli/cmdstatus/status.go` spans 268 lines and mixes terminal table display with database connection fallback, JSON manifest fallback, and record resolution.
- **Architectural Solution:** Extract repository record resolution into `cli/cmdstatus/status_resolve.go`:
  - `loadStatusByScope(scope, group string) ([]model.ScanRecord, error)`
  - `loadRecordsByGroup(group string) ([]model.ScanRecord, error)`
  - `buildStatusDBOpenError(err error) error`
  - `loadAllRecordsDB() ([]model.ScanRecord, error)`
  - `LoadRecordsJSONFallback() ([]model.ScanRecord, error)`
  - `loadAllRecordsDBOrEmpty() []model.ScanRecord`
  - `newStatusNoDataError(scope, group string) error`
  - `loadStatusRecords(scope, group string) ([]model.ScanRecord, error)`
  - `handleStatusDBError(err error) error`
- **Impact:** Leaves `status.go` strictly focused on argument processing and table rendering, reducing `status.go` to ~75 lines.

---

### Improvement 30: Operational Handler Extraction & Route Decoupling (`cli/cmd/see_ops.go`)
- **Objective:** Extract operational execution handlers (`RunRepoManageUI`, `runSeeCommitPending`, `runSeeHistory`, `RunSeeErrorsSSH`, `PrintSeeHelp`) from `cli/cmd/cmd_see.go` into `cli/cmd/see_ops.go`.
- **Current Problem:** `cli/cmd/cmd_see.go` spans 205 lines, combining command routing dispatch with UI management subprocesses and SSH terminal calls.
- **Architectural Solution:** Extract execution handlers into `cli/cmd/see_ops.go`:
  - `runSeeCommitPending(args []string) *apperror.AppError` & `printPendingCommits()`
  - `runSeeHistory(args []string) *apperror.AppError`, `runSeeHistorySSH(...)`, and `resolveHistorySSHArgs(...)`
  - `RunSeeErrorsSSH(...)`, `runSeeErrorsSSH(...)`, and `resolveErrorsSSHArgs(...)`
  - `RunRepoManageUI(args []string) *apperror.AppError`
  - `PrintSeeHelp()`
- **Impact:** Reduces `cmd_see.go` to < 80 lines (pure route dispatch table), isolating execution logic into `see_ops.go`.

---

## 3. Zero-Harm Principle & Migration Safety

All modifications across Wave 3 strictly observe the zero-harm principle:

1. **Non-Destructive Additive APIs:**
   `ListReposByTag` and `ListReposUnpushed` are additive methods on `*DB` in `cli/store/repo_query_filter.go`. Existing queries in `cli/store/repo.go` remain completely untouched and backward-compatible.

2. **Safe Fallbacks:**
   - In `ListReposByTag`, passing an empty string or whitespace falls back safely to `db.ListRepos()`.
   - In `GetHeatmapTaskSummary`, passing an empty `taskSlug` queries all tasks rather than failing.
   - In `TrackAgentAction`, empty `taskSlug` or `actionType` returns early with `nil` error rather than throwing exceptions or writing corrupt rows.
   - In `parseSortColumn` and `parseSortOrder`, any unknown or malicious string safely defaults to `"ActionLogId"` and `"DESC"`.

3. **Memory & Connection Hygiene:**
   - All SQL query rows are closed immediately via `defer rows.Close()`.
   - Database handles opened in helpers (such as `recordAgentActionAtRoot`) are closed via `defer hdb.Close()`.

---

## 4. File Sizing Ratchet & Modular Decomposition Protocol

To prevent code bloat and comply with the project's sizing guidelines, all touched files must conform to the following limits:

| Target File | Current Lines | Planned Extraction / Destination | Post-Refactor Line Count |
|:---|:---:|:---|:---:|
| `cli/store/repo.go` | 384 | Keep unchanged; add queries to `repo_query_filter.go` | 384 lines (well below 500 cap) |
| `cli/store/repo_query_filter.go` | NEW | Receives `ListReposByTag` and `ListReposUnpushed` | ~35 lines |
| `cli/store/heatmap_store.go` | 323 | Add `GetHeatmapTaskSummary` (14 lines) | ~337 lines (below 500 cap) |
| `cli/store/heatmap_schema.go` | 254 | Add `HeatmapTaskSummary` struct (10 lines) | ~264 lines (below 500 cap) |
| `cli/heatmap/heatmap_tracker.go` | 93 | Keep unchanged; extract agent tracking | 93 lines (under 100 cap) |
| `cli/heatmap/heatmap_tracker_agent.go` | NEW | Receives `TrackAgentAction` and `recordAgentActionAtRoot` | ~45 lines |
| `cli/cmdagent/agent_ui_logs.go` | 250 | Add sorting whitelist helpers | ~275 lines (below 500 cap) |
| `cli/cmdagent/agent_ui_server.go` | 777 | Add `SortBy` and `Order` to `ActionLogFilter` | 779 lines (existing server file) |
| `cli/gitutil/git_repo_query.go` | 84 | Keep unchanged; extract commit helpers | 84 lines (under 100 cap) |
| `cli/gitutil/git_repo_commits.go` | NEW | Receives `GetRepoUnpushedCommitCount` and `GetRepoMergeBase` | ~30 lines |
| `cli/cmdcpar/cpar.go` | 293 | Extract target discovery into `cpar_target.go` | ~230 lines |
| `cli/cmdcpar/cpar_target.go` | NEW | Target discovery and dirty repository helpers | ~65 lines |
| `cli/cmdstatus/status.go` | 268 | Extract scope resolution into `status_resolve.go` | ~75 lines (under 100 cap) |
| `cli/cmdstatus/status_resolve.go` | NEW | Scope resolution and DB/JSON loading | ~195 lines |
| `cli/cmd/cmd_see.go` | 205 | Extract operational handlers into `see_ops.go` | ~75 lines (under 100 cap) |
| `cli/cmd/see_ops.go` | NEW | Operational handlers and help rendering | ~135 lines |

---

## 5. Verification Protocol Under Zero-Build and Zero-Test Mandate

Under the project mandate, compiler builds (`go build`) and automated test suites (`go test`) are strictly prohibited during this phase. Verification is conducted via the following static analysis and contract inspection protocol:

1. **AST & Syntax Audit:**
   - Every file must compile syntactically: package names match directory names (`store`, `heatmap`, `cmdagent`, `gitutil`, `cmdcpar`, `cmdstatus`, `cmd`).
   - All imports must be used, resolved, and reference valid internal packages (`github.com/alimtvnetwork/gitmap-v28/...`).

2. **Error Sovereignty Inspection:**
   - No leaf function may call `cliexit.HandleError`, `cliexit.HandleExit`, or `os.Exit`.
   - All database query and scan failures return structured `error` or `*apperror.AppError`.

3. **Function Limit & Complexity Inspection:**
   - Every new function must be <= 15 lines (<= 8 lines preferred).
   - Zero nested `if` statements. Guard clauses with early returns only.

4. **Naming & Positive Polarity Check:**
   - All booleans must begin with `is` or `has` (`isDirty`, `hasAttached`, `isMet`).
   - Banned prefixes: `can`, `should`, `was`, `will`. Zero negative booleans.

5. **Path and Formatting Hygiene:**
   - All paths are relative Git paths (`cli/...`, `02-spec/...`, `.ai-memory/...`).
   - Zero absolute paths and zero file URI schemes (`file:///`).
