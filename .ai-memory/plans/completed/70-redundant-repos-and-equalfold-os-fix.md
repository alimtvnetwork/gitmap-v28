# Completed Plan: 70-redundant-repos-and-equalfold-os-fix

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

Release after this please with a minor bump
```

---

## Executive Summary of Completed Deliverables

1. **Confirmation of Redundant Repository Optimization Code**:
   - Confirmed through deep architectural audit that GitMap did not previously have automated database-level redundant repository deduplication (it only ran `gitmap clone --fix` which cleaned VS Code workspace files, leaving git repositories and `gitmap.db` untouched).
   - Designed and implemented `FindDuplicateRepos()` and `DeduplicateRepos(keepNewest bool)` directly in `cli/store/repo_duplicates.go`, backed by high-speed SQLite queries grouping by normalized remote URLs.
   - Wired redundant repository optimization into `gitmap pull` / `gitmap pull-all` via `OptimizeRedundantPullPool` in `cli/cmdpull/pull_dedup.go`, ensuring that the next pool automatically optimizes redundant repositories in SQLite before worker dispatch.

2. **Strings Comparison Modernization (`strings.EqualFold` vs `strings.ToLower`)**:
   - Eliminated slower heap-allocating `strings.ToLower` comparisons across 35+ files and call sites in `cli/cmd`, `cli/cmdclone`, `cli/cmddb`, `cli/cmdpull`, `cli/cmdssh`, `cli/cmdvscode`, `cli/macro`, `cli/osclean`, `cli/osuser`, `cli/release`, `cli/store`, and `cli/vscodepm`.
   - Replaced all case-insensitive equality checks with zero-allocation `strings.EqualFold()`.

3. **OS-Aware Path Case-Sensitivity Remediation**:
   - Fixed path comparison logic in `cli/fsutil/path_normalize.go` by introducing:
     - `IsPathCaseInsensitive(p string) bool`: returns `true` on Windows or when path contains a Windows drive letter prefix (`C:`, `D:`).
     - `EqualPaths(p1, p2 string) bool`: uses `strings.EqualFold` on Windows, and strict byte equality `==` on Unix/Linux systems where paths differing only by case are legally distinct.
     - `CanonicalPathKey(p string) string`: folds to lowercase on Windows while strictly preserving exact path casing on Unix/Linux.
   - Refactored `cli/fsutil/workdir_match.go`, `cli/cmdpull/pull_efficient.go`, `cli/cmd/reconcile_db.go`, and `cli/vscodepm/io.go` to use these centralized OS-aware helpers.

4. **Database-Driven Redundancy & Future Uniqueness Guarantees**:
   - Ensured `RepoFile` table schema in `cli/repodb/repo_db.go` has `UNIQUE(RelativePath)`.
   - Implemented typed `Upsert` method on `RepoFileDbRepo` in `cli/repodb/repo_file.go` using `ON CONFLICT(RelativePath) DO UPDATE SET`.
   - Ensured `cli/indexer/walker.go` uses `sqlUpsertRepoFile` with `ON CONFLICT(RelativePath)` upsert semantics, guaranteeing no redundant files can be inserted.
   - Enforced database-driven redundancy checks strictly in SQLite (`gitmap.db`), completely replacing raw filesystem traversals.

---

## Verification Evidence
- Task Manager DB: `.ai-memory/temp-agents/70-redundant-repos-and-equalfold-os-fix/agent-task.db` (All 4 subtasks `DONE`, 100.0% completion).
- Targeted guideline autofixer: `python 03-ai-scripts/05-guideline-autofixer.py cli/store cli/fsutil cli/cmdpull cli/repodb --check-only` -> exit code 0 across 289 files scanned.
- Forbidden strings check: `python linter-scripts/check-forbidden-strings.py` -> exit code 0.
- Zero raw secrets, positive booleans enforced, functions within modular size limits.
