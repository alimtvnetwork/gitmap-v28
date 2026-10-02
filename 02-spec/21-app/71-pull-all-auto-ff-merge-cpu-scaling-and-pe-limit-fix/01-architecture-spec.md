# Architecture Specification: Pull-All Auto Fast-Forward Merge, CPU Auto-Scaling, Batch Failure Remediation & PE Limit Fix

## 1. System Overview

This specification establishes the architectural design for five core capabilities in GitMap:
1. **CLI Flag Disambiguation for Pipeline Errors (`gitmap pe`)**: Proper parsing of limit flags (`-l`, `-limit`, `--limit`, `-n`, `--lines`) without false-positive repository target collisions.
2. **Local Repository Identity & No-Remote Pull Handling**: Resolving repository names from filesystem paths when remote URLs are absent (eliminating the `"unknown (repo-cache)"` bug) and skipping local-only repositories during pull.
3. **Automated Divergent Branch Fast-Forward & Merge Fallback (Strictly No Rebase)**: Multi-token divergence recognition in `safe_pull.go`, followed by clean `--no-rebase` merge execution and immediate conflict abort.
4. **Batch Failure Remediation & Interactive Console Prompt**: An interactive resolution prompt (`[1] All at once`, `[2] One by one`, `[q] Skip`) at the end of `gitmap pa`, backed by a dedicated batch command (`gitmap pull-fix` / `gitmap pf` / `gitmap pa --fix`).
5. **CPU-Aware Concurrency Auto-Scaling & Presets**: Dynamic concurrency scaling based on system logical CPU counts and load pressure, providing `--high-perf` (turbo) and `--low-cpu` (conservative) presets.

---

## 2. Architecture & Data Flow

```
                      +---------------------------------------+
                      |       User CLI Execution Flow         |
                      +---------------------------------------+
                                          |
                +-------------------------+-------------------------+
                |                                                   |
      [ gitmap pe -l 10 ]                                  [ gitmap pa / pull ]
                |                                                   |
  +---------------------------+                        +---------------------------+
  |    pipeline_flags.go      |                        |    pull_concurrency.go    |
  |  - isValueFlag check      |                        |  - CPU Auto-Scaling       |
  |  - parseLimitFlag         |                        |  - High-Perf / Low-CPU    |
  |  - Disambiguate -n vs -l  |                        +---------------------------+
  +---------------------------+                                     |
                |                                      +---------------------------+
  +---------------------------+                        |      pull_worker.go       |
  |    pipeline_logs.go       |                        |  - Skip no-remote repos   |
  |  - Apply limit to logs    |                        |  - Fallback repo name     |
  +---------------------------+                        +---------------------------+
                                                                    |
                                                       +---------------------------+
                                                       |       safe_pull.go        |
                                                       |  - ff-only pull attempt   |
                                                       |  - isDivergedOutput check |
                                                       |  - --no-rebase merge      |
                                                       |  - merge --abort on fail  |
                                                       +---------------------------+
                                                                    |
                                                       +---------------------------+
                                                       |     pull_db_sync.go       |
                                                       |  - Record pull errors     |
                                                       +---------------------------+
                                                                    |
                                                       +---------------------------+
                                                       |    Interactive Prompt     |
                                                       |  - Resolve all at once    |
                                                       |  - Step one by one        |
                                                       |  - gitmap pull-fix (pf)   |
                                                       +---------------------------+
```

---

## 3. Detailed Component Architecture

### 3.1 Pipeline Error Flags Disambiguation (`cli/cmdpipeline`)
- **Problem**: When a user executes `gitmap pe -l 10` or `gitmap pe -limit 10`, `parseRepoTarget` treats `-l` as an unknown flag, skips it, and then mistakes the trailing argument `"10"` for a repository slug. Because no repository is named `"10"`, GitMap terminates with `Repository not found: "10"`.
- **Solution**:
  - Add `Limit int` and `HasLimit bool` to `PipelineErrorFlags`.
  - Update `isValueFlag` in `pipeline_flags.go` to explicitly recognise `-l`, `-limit`, `--limit`, `-n`, `--lines`, and `-lines`.
  - Disambiguate `-n`: if followed by a valid positive integer, treat as `Limit = N`; otherwise, treat as `--no-output-log`.
  - Ensure `isConsumedFlagValue` recognizes `flags.Limit` string representation so numeric limits are never captured as repository targets.
  - Apply `flags.Limit` in `pipeline_logs.go` to cap rendered errors, warnings, and failure history runs.

### 3.2 Repository Identity & Remote-Less Handling (`cli/mapper`, `cli/store`, `cli/cmdpull`)
- **Problem**: Local repositories created without remote origins (such as `d:\work\repo-cache` created by special repository helpers) have empty `git config --get remote.origin.url`. `extractRepoName("")` returns `"unknown"`, causing the repository to be saved in SQLite with `Slug = "unknown"` and `RepoName = "unknown"`. During `gitmap pa`, `SafePull` fails 4 times trying to pull a repository without remotes.
- **Solution**:
  - In `cli/mapper/mapper.go`, when `remoteURL` is empty, derive `repoName` from `filepath.Base(repo.AbsolutePath)`.
  - In `cli/store/repo_sanitize.go` / `repo_duplicates.go`, add self-healing: on startup or maintenance, any repository where `Slug = "unknown"` or `RepoName = "unknown"` is updated to `filepath.Base(AbsolutePath)`.
  - In `cli/cmdpull/pull_worker.go`, inspect repository remotes before launching git pull. If the repository has neither HTTPS nor SSH remotes, mark it as `PullStepTypeSkipped` with reason `"local-only (no remote)"` and do not execute retry loops.

### 3.3 Fast-Forward & Merge Fallback (Strictly No Rebase) (`cli/cloner`)
- **Problem**: Git pull divergence detection in `safe_pull.go` only checked `strings.Contains(output, "diverged")`, failing to match standard git messages like `hint: You have divergent branches` or `fatal: Not possible to fast-forward, aborting`.
- **Solution**:
  - Broaden `isDivergedOutput(output string) bool` to check multiple tokens case-insensitively using `strings.Contains`: `"diverg"`, `"not possible to fast-forward"`, `"cannot fast-forward"`, `"divergent branches"`, `"need to specify how to reconcile"`.
  - In `attemptAutoMergePull`, run:
    `git -C <dir> pull --progress --no-rebase --no-edit --autostash`
  - If a merge conflict occurs, immediately run `git -C <dir> merge --abort` and return structured `ErrMergeConflict`, avoiding unnecessary 4x retry loops.

### 3.4 Batch Failure Remediation & Interactive Prompt (`cli/cmdpull`)
- **Problem**: When multiple repositories fail during `gitmap pa`, the CLI immediately exits, requiring the user to run remediation manually one repo at a time.
- **Solution**:
  - In `pull_db_sync.go`, ensure all pull failures are recorded in the SQLite database (`PullError` / `PullBatchError`).
  - At the end of `gitmap pa`, if failures occurred and standard input is an interactive terminal (`isatty` / `term.IsTerminal`):
    - Display interactive prompt:
      ```
      [?] Remediate failed repositories?
        [1/a] Resolve all failed repositories at once
        [2/s] Step through one by one
        [q/n] Skip / Exit
      ```
    - If user chooses `1`, invoke batch remediation across all failed repositories.
    - If user chooses `2`, step through failed repositories sequentially.
  - Expose direct non-interactive CLI commands:
    `gitmap pull-fix` (alias `gitmap pf`, `gitmap pa --fix`).

### 3.5 CPU-Aware Concurrency Auto-Scaling (`cli/cloneconcurrency`, `cli/cmdpull`)
- **Problem**: Unconditional concurrency (e.g. `runtime.NumCPU()`) causes severe thread thrashing and memory pressure on systems with 16+ logical cores when running intensive git operations.
- **Solution**:
  - Adaptive Concurrency Table based on logical core count:
    - 1–2 cores: Concurrency = 1
    - 3–4 cores: Concurrency = 2
    - 5–8 cores: Concurrency = 3
    - 9–16 cores: Concurrency = 4
    - 17+ cores: Concurrency = 6 (capped to avoid I/O starvation)
  - Configuration Presets:
    - `--high-perf` / `--turbo`: concurrency = `min(NumCPU, 12)`, no inter-repo delay.
    - `--low-cpu` / `--conservative`: concurrency = `max(1, NumCPU / 4)` (capped at 2), throttling delay enabled.
  - Expose `--auto-scale`, `--high-perf`, `--low-cpu`, and `--concurrency <N>` flags in `gitmap pa` and update `cli/helptext/pull-all.md`.
