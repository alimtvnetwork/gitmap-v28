# Subtask Plan 02: Pull-All Fast-Forward Merge, Batch Remediation & CPU Auto-Scaling

## 1. Overview & Assigned Subtasks
- **Subtask 3 (Task-03)**: Divergent pull detection and automatic Fast-Forward merge fallback (Strictly `--no-rebase`).
- **Subtask 4 (Task-04)**: Batch failure remediation command (`gitmap pull-fix`, `gitmap pf`, `gitmap pa --fix`) and interactive resolution prompt.
- **Subtask 5 (Task-05)**: CPU-aware concurrency auto-scaling, configuration presets (`--high-perf`, `--low-cpu`), and help documentation.

---

## 2. Implementation Steps

### Subtask 3: Fast-Forward & Merge Fallback (Strictly No Rebase)
1. **Edit `cli/cloner/safe_pull.go`**:
   - Broaden `isDivergedOutput(output string)` to check multiple tokens case-insensitively: `"diverg"`, `"not possible to fast-forward"`, `"cannot fast-forward"`, `"reconcile divergent"`.
   - Update `attemptAutoMergePull` to execute:
     `git -C <dir> pull --progress --no-rebase --no-edit --autostash`
   - If merge conflicts occur (detected via `"CONFLICT"` in output or non-zero exit status with unmerged paths), immediately run `git -C <dir> merge --abort` and return structured `ErrMergeConflict`, avoiding multiple 4x retry loops.

### Subtask 4: Batch Remediation & Interactive Prompt
1. **Edit `cli/cmdpull/pull_db_sync.go`**:
   - Ensure all pull failures from the current session are recorded in SQLite (`pull_errors` table or store).
2. **Create/Update `cli/cmdpull/pull_remediation.go`**:
   - Create `PromptInteractiveBatchFix(failedRepos []PullFailureItem) error`:
     - Checks if standard input is an interactive terminal (`isatty.IsTerminal`). If not, skips prompt.
     - Displays options: `[1/a] Resolve all failed repositories at once`, `[2/s] Step through one by one`, `[q/n] Skip / Exit`.
     - Dispatches appropriate remediation actions based on user selection.
   - Implement `RunBatchPullFix(args []string) error` to remediate failed repositories in batch.
3. **Edit `cli/cmdpull/pull.go` & `cli/cmd/rootutility.go`**:
   - Call `PromptInteractiveBatchFix` at the end of `finalizePullBatchTask` before exit.
   - Register CLI command `gitmap pull-fix` (alias `gitmap pf`, `gitmap pa --fix`).

### Subtask 5: CPU Auto-Scaling & Concurrency Presets
1. **Create/Update `cli/cloneconcurrency/pull_concurrency.go`**:
   - Implement `ResolveAdaptiveConcurrency(preset ConcurrencyPreset, userLimit int) int`:
     - Default Auto-Scale: 1-2 cores -> 1; 3-4 cores -> 2; 5-8 cores -> 3; 9-16 cores -> 4; 17+ cores -> 6.
     - `--high-perf` / `--turbo`: `min(NumCPU, 12)` workers, 0ms delay.
     - `--low-cpu` / `--conservative`: `max(1, NumCPU / 4)` workers (capped at 2), throttling pause.
2. **Edit `cli/cmdpull/pull.go` & `cli/cmdpull/pull_efficient.go`**:
   - Parse `--auto-scale`, `--high-perf`, `--turbo`, `--low-cpu`, `--conservative`, `--concurrency`.
   - Apply calculated concurrency to worker pool.
3. **Update Help Documentation**:
   - Edit `cli/helptext/pull-all.md` and `cli/cmdpull/pull_help.go` to document CPU presets and flags.

---

## 3. Verification Commands
- `python 03-ai-scripts/05-guideline-autofixer.py cli/cloner cli/cmdpull cli/cloneconcurrency --check-only`
- `python linter-scripts/check-relative-paths.py`
