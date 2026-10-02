# Completed Plan: 72-repo-dedup-os-aware-equalfold

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

learn /learn if you have to learn something and /plan stuff before working please.
```

---

## Executive Summary of Completed Deliverables

1. **Confirmation of Redundant Repository Optimization Code**:
   - Confirmed that automated database-level redundant repository discovery and deduplication is implemented in:
     - `cli/store/repo_duplicates.go`: `FindDuplicateRepos()` and `DeduplicateRepos(keepNewest bool)` directly query SQLite `gitmap.db` by remote URL and paths.
     - `cli/cmdpull/pull_dedup.go`: `OptimizeRedundantRepos()` executes before batch pull / pool worker dispatch.
     - `cli/cmd/find_duplicates_git.go`: CLI command `gitmap find-duplicates git --fix` enables immediate on-demand deduplication.

2. **Strings Comparison Modernization (`strings.EqualFold` & `strutil.EqualFoldAny`)**:
   - Created `cli/strutil/strutil.go` providing zero-allocation Unicode case-folding helpers: `EqualFoldAny`, `EqualFoldAnyTrim`, and `NormalizeLowerTrim`.
   - Modernized `cli/cmdvscode/vscode_cmd.go` and `cli/cmdvmware/vmware.go` to eliminate naive lowercase conversions and evaluate arguments using `strutil.EqualFoldAny`.

3. **OS-Aware Path Case-Sensitivity Remediation**:
   - Fixed critical index creation bug in `cli/store/store.go`:
     - Introduced `getRepoAbsPathIndexQuery()` and `getScanFolderPathIndexQuery()` checking `runtime.GOOS == "windows"`.
     - On Windows: enforces `COLLATE NOCASE`.
     - On Unix: enforces standard binary collation, preserving distinct casing for legal Unix paths.
   - Refactored `cli/store/storage_inventory.go`, `cli/vscodepm/path_filter.go`, and `cli/workspacesync/path_guard.go` to use `fsutil.CanonicalPathKey` and `fsutil.EqualPaths`, ensuring path casing is preserved on Unix and only folded on Windows.

4. **Database-Driven Redundancy & Future Uniqueness Guarantees**:
   - Ensured `RepoFile` table schema in `cli/repodb/repo_db.go` has `UNIQUE(RelativePath)` with OS-aware collation.
   - Verified `RepoFileDbRepo.Upsert` in `cli/repodb/repo_file.go` uses `ON CONFLICT(RelativePath) DO UPDATE SET` upsert semantics.
   - Enforced database-driven redundancy checks strictly in SQLite (`gitmap.db`), completely avoiding raw filesystem traversals.

5. **Release & Version Bump**:
   - Bumped version to `6.464.0` in `version.json`, `package.json`, `cli/constants/constants.go`, `changelog.md`, and synchronized release manifests.
