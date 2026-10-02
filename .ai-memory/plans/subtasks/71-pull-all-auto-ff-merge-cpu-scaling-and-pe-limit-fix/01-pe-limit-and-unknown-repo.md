# Subtask Plan 01: Pipeline Error Limit Flag Parsing & Unknown Repo Resolution

## 1. Overview & Assigned Subtasks
- **Subtask 1 (Task-01)**: Fix `gitmap pe -l 10`, `-limit 10`, `--limit 10`, `-n 10` flag parsing and output limit rendering.
- **Subtask 2 (Task-02)**: Resolve `"unknown"` repo-cache name and slug resolution, DB sanitization, and skip no-remote repos in `gitmap pa`.

---

## 2. Implementation Steps

### Subtask 1: PE Limit Flag Parsing
1. **Edit `cli/cmdpipeline/pipeline_flags.go`**:
   - Add fields `Limit int` and `HasLimit bool` to `PipelineErrorFlags`.
   - Implement `parseLimitFlag(args []string) (int, bool)` supporting `-l`, `-limit`, `--limit`, `--lines`, `-lines`, and `-n` (when followed by a number), plus inline syntax `-l=N`, `--limit=N`, etc.
   - Update `hasSuppressOutputArg(args []string)` so `-n` is NOT treated as output suppression if followed by a number.
   - Update `isValueFlag(flag string)` to include `-l`, `-limit`, `--limit`, `--lines`, `-lines`, and `-n`.
   - Update `isConsumedFlagValue(val string, flags *PipelineErrorFlags)` to return true if `val == strconv.Itoa(flags.Limit)`.
   - Wire `parseLimitFlag` into `parseCommonErrorFlags`.
2. **Edit `cli/cmdpipeline/pipeline_logs.go`**:
   - In `processAndRenderErrorLogs`, pass `flags.Limit` to control rendered line count.
   - In `capFailedRuns`, respect `flags.Limit` if set.
   - In `printPipelineErrorLogsHelp`, document `-l, --limit, -n, --lines <N>`.

### Subtask 2: Unknown Repo Resolution & No-Remote Pull Handling
1. **Edit `cli/mapper/mapper.go`**:
   - In `buildOneRecord`, if `remoteURL` is empty or `extractRepoName(remoteURL)` is `"unknown"`, derive repository name from `filepath.Base(repo.AbsolutePath)`.
   - Ensure `buildSlug` produces a valid slug based on the resolved repository name.
2. **Edit `cli/store/repo_sanitize.go` (or create if needed)**:
   - Implement `SanitizeUnknownRepos(db *sql.DB)`: check for any repo with `Slug = 'unknown'` or `RepoName = 'unknown'`. If its `AbsolutePath` exists on disk, update `RepoName = filepath.Base(AbsolutePath)` and `Slug = strings.ToLower(filepath.Base(AbsolutePath))`.
   - Call this during DB initialization or startup maintenance.
3. **Edit `cli/cmdpull/pull_worker.go`**:
   - In `runTrackedPullLifecycle`, check if the repo has an upstream remote (`rec.HTTPSUrl == "" && rec.SSHUrl == "" && !gitutil.HasRemote(rec.AbsolutePath)`).
   - If no remote is configured, mark the step as `PullStepTypeSkipped` with reason `"local-only (no remote)"` instead of executing failing pull retries.

---

## 3. Verification Commands
- `python 03-ai-scripts/05-guideline-autofixer.py cli/cmdpipeline cli/mapper cli/store cli/cmdpull --check-only`
- `python linter-scripts/check-relative-paths.py`
