# Subtask 223.3: Git Command Split-DB Heatmap Tracing

> **Subtask ID:** 223.3  
> **Target File:** `.ai-memory/plans/subtasks/223-ubuntu-cursor-gitmap-python-git-tracing-and-aum-search/03-git-split-db-heatmap-command-tracing.md`  
> **Parent Plan:** [.ai-memory/plans/223-ubuntu-cursor-gitmap-python-git-tracing-and-aum-search.md](../../223-ubuntu-cursor-gitmap-python-git-tracing-and-aum-search.md)  
> **Spec Reference:** [02-spec/21-app/223-ubuntu-cursor-gitmap-python-git-tracing-and-aum-search/02-component-and-search-spec.md](../../../../02-spec/21-app/223-ubuntu-cursor-gitmap-python-git-tracing-and-aum-search/02-component-and-search-spec.md)  
> **Status:** Pending  
> **Target Area:** `cli/cmd/rootgit.go`, `cli/cmd/commit_push.go`, `cli/store/command_history_split_db.go`  

---

## 1. Objective

Instrument all raw Git passthrough commands (`gitmap git <cmd>`) and internal Git execution workflows in `commit_push.go` to capture execution duration, timestamp, command line, and exit status into the Split-DB `commands.db` (`CommandHistory` table). This eliminates telemetry blind spots in developer contribution heatmaps and activity dashboards.

---

## 2. Problem Statement

1. **Unmonitored Passthrough Commands:** When users or AI agents execute commands such as `gitmap git status`, `gitmap git diff`, or `gitmap git log`, `runGitPassthrough` in `cli/cmd/rootgit.go` delegates directly to `execGitInheritCP(args...)` without measuring elapsed execution time or persisting telemetry.
2. **Missing Internal Git Telemetry:** In `cli/cmd/commit_push.go`, operations like `performCommitPush` (invoking `git commit` and `git push`) and `executePullCommitPush` (invoking `git pull --rebase`) execute raw child processes without recording entries into `CommandHistorySplitDB`.
3. **Incomplete Heatmap Representation:** The developer activity visualization rendered by `cli/dashboard/` and analytics components queries `CommandHistory` for daily and hourly event frequencies. Because raw Git operations were unlogged, Git-heavy development sessions appeared deceptively inactive in contribution heatmaps.

---

## 3. Implementation Details

### Step 3.1: Duration Measurement & Interception in `cli/cmd/rootgit.go`
- Update `runGitPassthrough(args []string) error` to:
  1. Record `start := time.Now()`.
  2. Execute `err := execGitInheritCP(args...)`.
  3. Calculate `durationMs := time.Since(start).Milliseconds()`.
  4. Extract `exitCode` using helper `extractExitCode(err)`.
  5. Asynchronously/safely call `recordGitCommandHistory(args, exitCode, durationMs)`.
  6. Return wrapped error if execution failed.

### Step 3.2: Safe Exit Code Extraction Helper
- Implement `extractExitCode(err error) int`:
  - If `err == nil`, return `0`.
  - Type-assert `err` to `*exec.ExitError` and return `exitErr.ExitCode()`.
  - Fall back to `1` for general execution errors.

### Step 3.3: Split-DB Persistence Routine
- Implement `recordGitCommandHistory(gitArgs []string, exitCode int, durationMs int64)`:
  - Guard against empty arguments.
  - Construct `cmdName := "git " + gitArgs[0]` and `cmdLine := "git " + strings.Join(gitArgs, " ")`.
  - Open `store.OpenCommandHistorySplitDB("")`.
  - Execute `histDB.InsertCommandRecord(cmdLine, cmdName, exitCode, durationMs)`.
  - Ensure any database error is silently discarded so user Git commands never fail due to database access issues.

### Step 3.4: Internal Git Operation Hooking in `cli/cmd/commit_push.go`
- In `cli/cmd/commit_push.go`:
  - Instrument `performCommitPush`:
    - After `git commit -m <msg>`, record entry `CommandName: "git commit"`.
    - After `git push`, record entry `CommandName: "git push"`.
  - Instrument `executePullCommitPush`:
    - After `git pull --rebase`, record entry `CommandName: "git pull"`.
  - Ensure all timing measurements use high-precision monotonic clocks (`time.Since`).

---

## 4. Verification & Testing Plan

1. **Unit Tests:**
   - Author tests in `cli/cmd/rootgit_test.go` verifying that `recordGitCommandHistory` inserts valid records with proper duration and exit codes into a temporary SQLite `CommandHistorySplitDB`.
2. **Integration Verification:**
   - Run `gitmap git status` locally.
   - Query `commands.db` via `sqlite3` or Go test to confirm that a record with `CommandName = "git status"` exists and `DurationMs >= 0`.
3. **Heatmap Query Check:**
   - Verify that `ListRecentCommands` in `cli/store/command_history_split_db.go` returns the newly traced Git commands.

---

## 5. Execution Checklist

- [ ] Inspect existing `cli/cmd/rootgit.go` and `cli/cmd/commit_push.go`.
- [ ] Add `extractExitCode` helper in `cli/cmd/rootgit.go`.
- [ ] Add `recordGitCommandHistory` helper in `cli/cmd/rootgit.go`.
- [ ] Wrap `runGitPassthrough` with execution timer and history recorder.
- [ ] Hook `performCommitPush` in `cli/cmd/commit_push.go` to record commit and push actions.
- [ ] Hook `executePullCommitPush` in `cli/cmd/commit_push.go` to record pull actions.
- [ ] Author unit test covering `recordGitCommandHistory` in `cli/cmd/rootgit_test.go`.
- [ ] Verify clean build with no compiler warnings or linter violations.

---

## 6. Acceptance Criteria

- [ ] **AC-223.3-1:** Running `gitmap git status` records a new entry in `commands.db` with `CommandName = "git status"` and accurate `DurationMs`.
- [ ] **AC-223.3-2:** Non-zero Git exit codes (e.g. `gitmap git checkout invalid-branch-name`) are accurately logged with `ExitCode != 0`.
- [ ] **AC-223.3-3:** If `commands.db` is temporarily locked or inaccessible, `gitmap git <cmd>` continues execution without error.
- [ ] **AC-223.3-4:** Internal Git commands from `gitmap commit-push` are logged into `CommandHistory` with granular sub-operation names (`git commit`, `git push`).
- [ ] **AC-223.3-5:** All new functions follow Go coding guidelines: zero nested ifs, guard clauses, and positive boolean names.
