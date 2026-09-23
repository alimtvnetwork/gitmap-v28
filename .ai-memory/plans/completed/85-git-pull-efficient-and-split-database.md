# Plan 85: Git Pull Efficient Engine, Non-Git Directory Fallback, and Dedicated SQLite Split-Database (`gitmap-pull.db`)

- **Status:** COMPLETED
- **Completed At:** 2026-09-23
- **Spec Reference:** [02-spec/21-app/136-git-pull-efficient-and-split-database.md](../../../02-spec/21-app/136-git-pull-efficient-and-split-database.md)
- **Execution Lifecycle:** 1 continuous self-loop across 8 discrete tasks, fully adhering to coding guidelines, zero-nesting, and split-DB standards.

---

## User Request (Verbatim)

```text
Let's introduce a new command for Git Map to pull. Okay? It would be a bit different than the traditional pull that we have. Let me share my ideas. So, okay. The first issue that when we do the Git pull, it does the pull. Okay? So let's say if we do the Git pull, and if it is not a Git repository, then you will by default go into the Git pull all. Okay? So that is like the traditional one that I want. Okay? Now, on top of this, we want to make a new command introduced that would be git map pull all hyphen efficient or AE. So that would also have a short form, git map PAE, if that is not there, or it will be pull hyphen AE. Okay, so all of these will be aliasing of the same efficient. Now, I will describe how the efficient will work. From now on, any pull that we do using the git map, that will be recorded to its SQLite database and it not the root database, but it will create its own split database called git map hyphen pull.db. It's the same folder as the data folder. The root db is close to that. It have a pull git map hyphen pull db. That db will only contain the information for all the repositories that we are pulling, the pull command when it runs, where it runs, and which repositories. So that means which repositories when I'm saying that means it needs to have a subtable, right? One-to-many relationship. So it needs to have this information like how many repositories has been, let's say, pulled. Okay, that's one table information. In another table, we need to know how many data actually pulled. We don't update the file in the db. Remember that. That will actually slow us down. But we update how many files it has changed and what was the last commit. Okay? Now, let's say user is running the git map pull or pull all, whatever that is running, it always saves into that database. Now, the efficients, when we say git map pull all efficient, that will have a hidden trigger. So it will observe the last 20 to 30 repositories, or 30 pulls, pull all. So basically, the efficient pull will behave like pull all, but a subtle difference. The difference is when we have 20 pulls, 20 pull information, and from the 20 pull information, if we find zero, let's say, commits are coming or data is changed on certain repositories, those will be marked as inactive. So in future, when we do the, let's say, pull effective in, let's say, in a one-day window. One day, remember that, one-day window. It will only pull the active ones. So it will skip the inactive ones. And it should mention like, these are the inactive ones we did not pull. If you want to pull the inactive ones, run pull all. Okay? So it will give a summary like these are the pulls that it has done. Okay? And also we can pull effective should not pull the table. Okay? So that it's not necessary. But also we can have a hyphen, hyphen flag like a hyphen, hyphen status or hyphen, hyphen status table. Both would actually show the table just like it does on the pull all or probably, we can have another command. So these flags will work, but also we can have another command like git map space pull hyphen all hyphen efficient hyphen table. So in short, it will have another short form will be git map space PAET. PAET. PAET. Okay? So when we run the short forms, it should always tell us the full form of the command, where it is running, so that user is aware of it. Also during the runtime, always include the version number of the git map so that we are aware of which version is running it. Are you clear with the requirements? Can you please implement these requirements? What do you think about your confidence level?
```

---

## 1. Executive Summary of Implementation

This release delivers the complete Git Pull Efficient suite:
1. **Non-Git Directory Fallback:** When executing `gitmap pull` outside a Git repository (without explicit slug/group), the command automatically defaults to `gitmap pull all` and notifies the user.
2. **Dedicated Split Database (`gitmap-pull.db`):** Co-located in the binary data directory (`BinaryDataDir()`) alongside the root database. Managed via `store.PullSplitDB` with normalized `PullRun` and `PullRepoRun` tables. Zero file blob writes to ensure sub-millisecond execution.
3. **Automated Telemetry Sync:** Both batch pulls and single-repo CWD pulls automatically sync execution metrics into `gitmap-pull.db`, recording commit hashes, commit message, author, duration, files changed count, and status flags.
4. **Inactivity Heuristic Engine:** Analyzes the last 20 to 30 runs from `gitmap-pull.db`. If a repository has 20 consecutive runs with zero commits and zero changed files within the last 24 hours, it is classified as inactive and skipped in efficient mode.
5. **Commands & Aliasing:**
   - Full command: `gitmap pull all-efficient` (Aliases: `pull-all-efficient`, `pull-ae`, `pae`)
   - Table command: `gitmap pull-all-efficient-table` (Alias: `paet`, flags `--status`, `--status-table`)
   - Direct subcommand dispatch: `gitmap pull all-efficient`, `gitmap pull ae`, `gitmap pull pae`, `gitmap pull all-efficient-table`, `gitmap pull paet`
6. **Command Expansion Banner:** Short-form commands (`pae`, `paet`, `pull-ae`) display the expanded command name, working directory (`cwd`), and GitMap version (`vX.Y.Z`).
7. **Concise Default Summary:** Efficient pull suppresses the heavy table by default, listing active pulls and clearly detailing skipped inactive repositories with guidance on how to force a full pull (`gitmap pull all`). Full table is available on demand via `paet`, `--status`, or `--status-table`.

---

## 2. Consolidated Subtask Register

### Subtask 01: Non-Git Directory Fallback to Pull-All
- **Traceability ID:** Task-02
- **Files Modified:** `cli/cmdpull/pull.go`, `cli/cmdpull/pull_fallback.go`, `cli/cmdpull/pull_fallback_test.go`
- **Delivered:** `ShouldFallbackToPullAll` and `AnnounceNonGitPullFallback` in `pull_fallback.go`. When outside a git worktree, `runPull` and `dispatchPullExecution` seamlessly set `opts.all = true`.

### Subtask 02: Dedicated Split Database Architecture (`gitmap-pull.db`)
- **Traceability ID:** Task-03
- **Files Created:** `cli/store/split_db_pull.go`, `cli/store/split_db_pull_types.go`, `cli/store/split_db_pull_ops.go`, `cli/store/split_db_pull_test.go`
- **Delivered:** Initialized `PullRun` (master) and `PullRepoRun` (detail) schemas in SQLite, connection pooling via `ConfigureSQLiteConn`, transaction-based batch insertions, and `EvaluateRepoInactivity`.

### Subtask 03: Pull Execution Recording & Telemetry Ingestion
- **Traceability ID:** Task-04
- **Files Created/Modified:** `cli/cmdpull/pull_db_sync.go`, `cli/cmdpull/pull_db_sync_types.go`, `cli/cmdpull/pull.go`
- **Delivered:** Intercepted `executePullBatchLifecycle` and `runPullCWDTracked` to record `PullRunRecord` and `PullRepoRunRecord` entries, capturing git commit metadata, files changed count, and durations without file content overhead.

### Subtask 04: Inactivity Analysis Engine & 24-Hour Active Repository Filtering
- **Traceability ID:** Task-05
- **Files Created:** `cli/cmdpull/pull_efficient.go`, `cli/cmdpull/pull_efficient_types.go`, `cli/cmdpull/pull_efficient_test.go`
- **Delivered:** Partitioned candidate repositories into active and inactive slices. Evaluated 20 consecutive runs within 24 hours. Formatted clear skip summaries with `gitmap pull all` guidance.

### Subtask 05: Efficient Pull Commands, Aliases & Table Display Modes
- **Traceability ID:** Task-06
- **Files Modified:** `cli/constants/constants_cli.go`, `cli/constants/cmd_constants_test.go`, `cli/cmd/rootcore.go`, `cli/cmd/rootgit.go`, `cli/cmd/clihelpers.go`, `cli/cmd/rootsuggest.go`
- **Delivered:** Registered `CmdPullAllEfficient`, `CmdPullAllEfficientAlias`, `CmdPullAE`, `CmdPullAllEfficientTable`, and `CmdPullAllEfficientTableAlias`. Connected root dispatcher and `git` subcommand passthrough.

### Subtask 06: Version Banner, CWD & Short-Form Command Expansion Output
- **Traceability ID:** Task-07
- **Files Created:** `cli/cmdpull/pull_banner.go`, `cli/cmdpull/pull_banner_test.go`
- **Delivered:** Formatted Lipgloss/ANSI banner announcing full command name, working directory, and normalized GitMap version string on short-form invocations.

### Subtask 07: Help Text, Documentation, Targeted Linting & Atomic Commit
- **Traceability ID:** Task-08
- **Files Created/Modified:** `cli/helptext/pull.md`, `cli/helptext/pull-all-efficient.md`, `cli/helptext/pull-all-efficient-table.md`, `cli/helptext/pull-ae.md`, `cli/helptext/pae.md`, `cli/helptext/paet.md`
- **Delivered:** Comprehensive markdown documentation with examples, flags, and cross-references. Verified 0 linting violations across all touched files.
