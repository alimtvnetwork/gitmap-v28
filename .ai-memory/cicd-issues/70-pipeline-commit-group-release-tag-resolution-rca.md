# RCA: Recent Commits Summary Table Release Tag Resolution Failure Outside Repository Root

- **Date:** 2026-09-21
- **Affected Subcommand:** `gitmap pe` (Recent Commits Pipeline Summary table)
- **Target Release:** `v6.292.0`
- **Scope:** Pipeline CLI rendering, multi-run commit group release extraction, and external working directory resolution

---

## 1. Symptom

When executing `gitmap pe` from an arbitrary directory outside a local git repository (e.g., `PS C:\Users\Administrator> gitmap pe`), the Recent Commits Pipeline Summary table rendered the latest release commit with a hyphen (`-`) in the `Release` column instead of the release version tag `v6.292.0`, despite the commit having been tagged and merged:
```text
  ● Recent Commits Pipeline Summary (Last 5 Commits):
    Offset   Commit    Branch               Release     Status     Workflows                        Failures
    ------   ------    ------               -------     ------     ---------                        --------
    latest   ba9ccc9   main                 -           RUNNING    CI [], CI Beacon [PASS] (+4)     0
    -1       3a629d9   release/v6.291.0     v6.291.0    RUNNING    CI Beacon [PASS], CI [RUN] (+4)  0
    -2       3ac9771   release/v6.290.0     v6.290.0    PASS       CI Beacon [PASS], CI [PASS] (+4) 0
```

---

## 2. Root Cause

The issue was caused by three interacting bugs in the commit grouping and release resolution architecture:

1. **Commit Group Multi-Run Branch Discard**: In `cli/cmdpipeline/pipeline_commit_groups.go`, `GroupRunsByCommit` processes runs sequentially. For a single release commit (e.g. `ba9ccc9`), GitHub Actions triggers workflow runs across multiple refs: `main`, `release/v6.292.0`, and `v6.292.0`. The first run returned by the GitHub API happened to be on branch `main`, so `initCommitGroup` recorded `HeadBranch = "main"`. When subsequent runs arrived for the same commit with `HeadBranch = "release/v6.292.0"` or `v6.292.0`, `groupBuilder.addRun` appended the workflow but completely ignored `r.HeadBranch`, discarding the release branch information.
2. **Missing `DisplayTitle` in GitHub Run Query**: Commit titles for release commits contain explicit release information (e.g., `release: v6.292.0 fix(...)`). However, `queryWorkflowRuns` in `pipeline_query.go` did not request `displayTitle` from `gh run list --json`, and `ghRunItem` lacked the field.
3. **Unanchored `git tag --points-at` Execution**: In `cli/cmdpipeline/pipeline_history.go`, `extractReleaseFromSha` invoked `exec.Command("git", "tag", "--points-at", sha)` in the current working directory. When `gitmap` was executed from `C:\Users\Administrator` (which is not a git repository), `git` failed with exit code 128 (`fatal: not a git repository`), causing `extractReleaseFromSha` to return an empty string and fall back to `"-"`.

---

## 3. Resolution

1. **Multi-Run Release Aggregation on Commit Groups**:
   - Added `Release string` to `CommitPipelineGroup` in `cli/cmdpipeline/pipeline_commit_groups.go`.
   - Updated `initCommitGroup` and `addRun` via `updateGroupRelease` to evaluate `r.HeadBranch` and `r.DisplayTitle` for every run in the commit group. If any run is on a release branch/tag or contains a release title, the group stores the resolved release version.
2. **Include `displayTitle` in GitHub Actions Query**:
   - Added `DisplayTitle` and `Event` to `ghRunItem` in `cli/cmdpipeline/pipeline.go`.
   - Updated `queryWorkflowRuns` in `cli/cmdpipeline/pipeline_query.go` to include `displayTitle,event` in the `--json` payload.
3. **Add `extractReleaseFromTitle`**:
   - Implemented title tokenization and semver parsing in `cli/cmdpipeline/pipeline_history.go` to extract release version tags directly from commit titles.
4. **Repo-Aware Git Tag Fallback**:
   - Updated `extractReleaseFromSha` in `cli/cmdpipeline/pipeline_history.go` to safely query the local repository path via read-only SQLite lookup from `gitmap.db` (`Repo` table) when the current directory is not a git repository, executing `git -C <repoPath> tag --points-at <sha>`.
5. **Unit Tests**:
   - Added unit tests in `cli/cmdpipeline/pipeline_history_test.go`:
     - `TestExtractReleaseFromTitle`
     - `TestResolveReleaseFromRun`
     - `TestGroupRunsByCommit_PreservesReleaseFromSubsequentRuns`
   - Verified 100% green tests and linters.

---

## 4. Prevention & Learnings

- **Multi-Source Ref Aggregation**: Workflows for a single commit SHA can be triggered by different branches, tags, and PR refs. Commit aggregation must harvest metadata (like release versions and branches) across all associated runs rather than relying solely on the first run processed.
- **CWD Independence for Git Commands**: Any internal `git` command execution that operates on a known repository must provide an explicit repository path (`git -C <repoDir>`) when running outside a git repository.
