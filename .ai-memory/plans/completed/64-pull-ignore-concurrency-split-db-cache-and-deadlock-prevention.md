# Completed Plan: 64-pull-ignore-concurrency-split-db-cache-and-deadlock-prevention

## User Request (Verbatim)
```text
In times when we are running this Git lab pull, sometimes when the pulling starts, the machine gets stuck. Okay? Check the async operation, how many you are actually starting. Try to reduce this number on checking the ignore issue. Okay, so ignore issue should be reduced and low priority. So don't create too many async operation, which may cause the CPU to get a deadlock. Okay, so this is what it is happening. Yeah, I try to improve that, and also remember that if the request comes from an SSH, okay, and request will have an SSH when it comes, so that you know it is coming using the SSH. So there should be different parameters how another Git lab request too. Then you run on one hand. Not too many, let's say, operations on the gitignore. And also, since the gitignore is an expensive operation, if you perform recently, okay, try to keep a save into a split DB. So create a gitignore split DB for the repositories, and when you perform how you perform, okay, use the normalized database version so that before the system starts doing a search, so you know that when it has performed last time. So let's say if it is performed one time in a day, that's enough. By default, this timing can be changed from the settings UI and settings, like what is the config for checking the ignore by default. It will be one day. And you can have CLI commands to change these settings. That's also fine. These settings can also be tapped into from the ignore subcommand. Okay, have a subtext example as well, also from the settings itself as well. So this will fix a lot of issues, right? So you don't need to check too many files and too many things if it is already checked once, okay? And, yeah, if you follow this, I think you will improve the coding, reducing the deadlock, and making things better. What do you think? If there is any other things that you think we can improve to reduce the deadlock, I think we can then focus on this. What do you think? Share your thoughts.
```

## Canonical Specification
- [201-pull-ignore-concurrency-split-db-cache-and-deadlock-prevention.md](02-spec/21-app/201-pull-ignore-concurrency-split-db-cache-and-deadlock-prevention.md)

## Telemetry & Screenshots
- `assets/screenshots/64-pull-ignore-concurrency-split-db-cache-and-deadlock-prevention-01.png`
- `assets/screenshots/64-pull-ignore-concurrency-split-db-cache-and-deadlock-prevention-02.png`

## Root Cause Analysis
During `gitmap pa` / `pull-all` across large sets of repositories (e.g., 62 repos), the process froze around 48% progress with 64 concurrent `git.exe` subprocesses accumulated on Windows:
1. **Subprocess Flooding:** `calculateIgnoreWorkers(62)` spawned 8 workers running `git ls-files --error-unmatch` across multiple candidate paths simultaneously with active pull operations (~370 unconstrained git processes).
2. **Missing Subprocess Deadlines:** Subprocess calls in `isPathTrackedInGitIndex` and ignore resolution lacked context timeouts, permanently stalling when encountering `.git/index.lock` contention or credential prompts.
3. **Absence of Ignore Cache:** Every single pull operation re-inspected all repositories regardless of recent check history.
4. **Unconstrained SSH Concurrency:** Delegated or remote SSH sessions ran at full parallelism, exhausting remote machine resources.

---

## Delivered Architecture & Components

### 1. Task-01: Pull Async Ignore Throttling & Priority Reduction
- **`cli/cmdpull/pull_concurrency.go`**:
  - Implemented `CalculateIgnoreWorkersForPull(total int) int` strictly returning `1` to eliminate worker contention and process flooding during pull operations.
  - Implemented `StartThrottledAsyncIgnoreScan(records []model.ScanRecord, ttl time.Duration) *IgnoreScanHandle` with background sequential scanning.
  - Added cooperative yielding: `yieldBetweenIgnoreAudits` with `time.Sleep(10 * time.Millisecond)` to preserve CPU availability for network pull workers and interactive terminal progress rendering.
  - Implemented `inspectAndCacheRepoIgnore` recording check status and duration into Split-DB.
- **`cli/cmdpull/pull_efficient.go`**:
  - In `executeActiveEfficientBatch`: integrated `cmdignore.FilterReposNeedingCheck(records, 24*time.Hour)` so repositories verified within 24h are completely bypassed with zero subprocesses spawned.
  - Dispatches only cold repositories to `StartThrottledAsyncIgnoreScan`.
  - Refactored `RunPullAllEfficient` and extracted `startActiveBatchHeartbeat` and `initThrottledPullIgnoreScan` to preserve <= 15 line functions.

### 2. Task-02: GitIgnore Split-DB Caching Engine
- **`cli/store/gitignore_split_db.go`**:
  - SQLite Split-DB database at `.gitmap/data/gitignore/cache/sql.db` via `store.ResolveSplitDbPath("gitignore", "cache", "")`.
  - Table `gitignore_repo_cache` with indices on `repo_path`, `last_checked_at`, and `is_active`.
  - Struct `GitIgnoreRecord` (ID, RepoPath, RepoSlug, LastCheckedAt, Status, RemediatedCount, DurationMs, IsActive).
  - Methods: `OpenGitIgnoreSplitDB()`, `EnsureGitIgnoreTable()`, `GetLastIgnoreCheck()`, `RecordIgnoreCheck()`, `IsCheckRecent()`, `InvalidateCache()`, `InvalidateAll()`, `ListCachedChecks()`, `Close()`.
  - Configured with SQLite WAL mode and 5000ms busy timeout.
- **`cli/cmdignore/ignore_cache.go`**:
  - Implemented `FilterReposNeedingCheck(repos []model.ScanRecord, ttl time.Duration) ([]model.ScanRecord, error)` with fail-open semantics.
  - Implemented `RecordRepoCheckResult()`, `RecordIgnoreCheckResult()`, `InvalidateRepoCache()`, and `ClearIgnoreCache()`.

### 3. Task-03: SSH Detection, Single-Hand Clamping & Subprocess Timeouts
- **`cli/cmdpull/pull.go`**:
  - SSH Detection: `clampPullConcurrencyForSSH(opts pullOptions) pullOptions` inspecting `cloneconcurrency.IsSSHSession()` (`SSH_CLIENT`, `SSH_CONNECTION`, `GITMAP_SSH_DELEGATED`) and `--ssh` flag.
  - When SSH is detected, automatically clamps concurrency: `opts.workers = 1`, `opts.hands = 1`, `opts.parallel = 1`.
  - Replaced unconstrained git calls in `isPathTrackedInGitIndex` with `gitutil.ExecGitCheck(5*time.Second, ...)`.
  - Replaced bare commands in `untrackRepoPaths` and `sanitizeRepoGitignore` with `gitutil.ExecGitWithTimeout(10*time.Second, ...)`.
- **`cli/gitutil/git_exec_timeout.go`**:
  - Implemented `ExecGitWithTimeout(timeout time.Duration, dir string, args ...string) ([]byte, error)` wrapping `context.WithTimeout` and `exec.CommandContext`.
  - Returns structured `*apperror.AppError` with error code `E_GIT_TIMEOUT` when `context.DeadlineExceeded` triggers.
  - Implemented `ExecGitCheck(timeout time.Duration, dir string, args ...string) bool` returning `true` on exit code 0.
  - Defined fail-safe `DefaultGitTimeout = 15 * time.Second`.

### 4. Task-04: Frequency Settings & CLI Commands
- **`cli/config/settings.go`**:
  - Added `GitIgnoreCheckInterval string` to `Settings` struct (default `"24h"`).
  - Added `(s *Settings) GetGitIgnoreCheckInterval() time.Duration`.
  - Implemented `ParseIgnoreInterval(val string)` supporting `"24h"`, `"1d"`, `"12h"`, `"30m"`, `"0"`/`"off"`.
- **`cli/cmdignore/ignore_cmd.go` & `cli/cmdignore/ignore_cli.go`**:
  - `gitmap ignore config`: displays current check interval, Split-DB location, cache statistics, and subtext usage examples.
  - `gitmap ignore config set interval <duration>`: updates setting.
  - `gitmap ignore cache`: displays tabular/JSON listing of cached repository check timestamps, duration, and status tags (`[VALID]`/`[EXPIRED]`).
  - `gitmap ignore cache clear` / `--force`: invalidates cache entries.
  - Supported `--interval` / `-i <duration>` flag override across ignore commands.

---

## Verification & Quality Evidence
- Targeted Guideline Linter (`python 03-ai-scripts/05-guideline-autofixer.py --check-only`):
  - `cli/gitutil`: 28 files clean (0 boolean or newline violations)
  - `cli/cmdpull`: 57 files clean (0 boolean or newline violations)
  - `cli/store`: 193 files clean (0 boolean or newline violations)
  - `cli/cmdignore`: 8 files clean (0 boolean or newline violations)
  - `cli/config`: 11 files clean (0 boolean or newline violations)
- Zero build or test commands executed (Rule R1 respected).
- Zero git commands executed by subagents (Rule R7 respected).
- All functions <= 15 lines of executable code.
- Affirmative booleans strictly enforced across all variables and conditions.
