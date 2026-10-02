# Completed Plan: 71-repo-dedup-equalfold-os-guarantee

## User Request (Verbatim)
```text
Do you have the code to improve the optimize the redundant repos that has been created? That's it. Do you have the code to do that? Confirm it. And also, I've seen in the code you are trying to compare the code in lowercase. Try not to do that. the specific method strings unfold, which is a lot more faster. Try to use that when you are comparing two strings without the case sensitivity. So try to apply that in terms of check-in, because in Unix, the different paths, different cases actually mean same thing. So you can only check ignoring the path in Windows. So you need to understand which OS you are in. So I think that is a bug we need to fix, and again, make a release and make sure you also optimize in the next pool if there is a redundancy. Okay? You can find the redundancy from your SQLite database, not from file system. Okay, so in future, when you add new files, you make sure that it is unique. Do you understand?

# Actionable Items Must Follow Non-Negotiable

1. Confirm if the code to optimize redundant repositories exists.
2. Avoid comparing code in lowercase; use the `strings unfold` method for case-insensitive comparisons.
3. Identify and fix the OS-specific bug related to path case sensitivity.
4. Ensure redundancy checks are performed using the SQLite database.
5. Guarantee uniqueness of new files added in the future.

Must follow and spawn agent using 

@[.agents/skills/execute-parent-task-with-n-steps-v6]

## Additional Instructions

learn /learn if you have to learn something and /plan stuff before working please./plan
```

---

## Executive Summary of Completed Deliverables

1. **Confirmation & Enhancement of Redundant Repo Optimization**:
   - Confirmed existing automated redundant repository optimization in `cli/store/repo_duplicates.go`, `cli/cmdpull/pull_dedup.go`, and `cli/cmd/find_duplicates_git.go`.
   - Enhanced `DeduplicateRepos()` in `cli/store/repo_duplicates.go` to prune BOTH remote URL duplicate groups and duplicate records sharing identical `AbsolutePath`s in a single transactional SQLite operation using OS-aware queries (`constants.SQLDeduplicateReposWindows` / `constants.SQLDeduplicateReposUnix`).
   - Ensured `cli/cmdpull/pull_dedup.go` runs this enhanced deduplication automatically before pulling.

2. **Strings Comparison Modernization (`strings.EqualFold`)**:
   - Refactored `cli/cmdignore/ignore_cmd.go` to replace `strings.ToLower(args[0]) == "interval"` with `strings.EqualFold(args[0], "interval")`.
   - Refactored `cli/cmdinstall/installantigravity_deploy_windows.go` in `hasDirInPathString` to eliminate per-iteration `strings.ToLower` allocations and compare path segments using `strings.EqualFold`.

3. **OS-Aware Path Case-Sensitivity Remediation**:
   - Refactored `cli/cmdvscode/find_duplicates_vscode.go` to use `fsutil.CanonicalPathKey(e.RootPath)` instead of unconditionally lowercasing `e.RootPath`.
   - Refactored `cli/cmdchromeprofile/find_duplicates_chrome.go` to use `fsutil.CanonicalPathKey(p.SourcePath)` instead of lowercasing `p.SourcePath`.
   - Windows paths fold case, while Unix paths strictly preserve case.

4. **SQLite Ingestion Uniqueness Guarantee for Future Files**:
   - Updated `cli/repodb/repo_db.go` so `IdxRepoFile_RelativePath` is created as `CREATE UNIQUE INDEX IF NOT EXISTS IdxRepoFile_RelativePath ON RepoFile(RelativePath COLLATE NOCASE);` on Windows, and standard binary collation on Unix.
   - Added `EnsureFileUniqueInDB(ctx context.Context, db *sql.DB, relPath string) (bool, error)` in `cli/repodb/repo_file.go`.

5. **Release & Version Bump**:
   - Bumped version to `6.462.0` in `version.json`.

---

## Verification Evidence
- Task DB: `.ai-memory/temp-agents/75-71-repo-dedup-equalfold-os-guarantee/agent-task.db` (4/4 subtasks DONE, 100%).
- Guideline autofixer: `python 03-ai-scripts/05-guideline-autofixer.py <folder> --check-only` passed exit 0 across `cli/store`, `cli/repodb`, `cli/cmdpull`, `cli/cmdvscode`, `cli/cmdchromeprofile`.
- Forbidden strings: `python linter-scripts/check-forbidden-strings.py` passed exit 0.
- Relative paths: `python linter-scripts/check-relative-paths.py` passed exit 0 across 8850 files.
