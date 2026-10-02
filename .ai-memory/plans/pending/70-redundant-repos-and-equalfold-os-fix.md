# Execution Plan: 70-redundant-repos-and-equalfold-os-fix

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

## 1. Architectural Synthesis & Findings

### Findings on Item 1 & 4 (Redundant Repo Optimization in SQLite)
- **Status Confirmed**: GitMap does NOT currently possess automatic database-level optimization for redundant repositories.
- `gitmap find-duplicates git` only prints duplicate remote URLs in memory and suggests `gitmap clone --fix`, which only touches VS Code configuration files without optimizing redundant repositories or `gitmap.db`.
- SQLite schema `Repo` lacks database-level duplicate detection queries.
- **Solution**:
  1. Add `db.FindDuplicateRepos()` in `cli/store/repo.go` that queries duplicate repositories using SQL `GROUP BY cleanUrl HAVING COUNT(*) > 1`.
  2. Add `db.DeduplicateRepos(keepNewest bool)` to delete duplicate rows from SQLite and clean cascading associations.
  3. Wire redundant repository optimization into `gitmap pull` / `gitmap pull-all` (the "next pool/pull") so redundant records are automatically deduplicated in SQLite before pull execution starts.
  4. Ensure redundancy check runs against SQLite (`gitmap.db`), NEVER traversing the filesystem.

### Findings on Item 2 (Zero-Alloc `strings.EqualFold` vs `strings.ToLower`)
- `strings.ToLower(a) == strings.ToLower(b)` causes heap allocations and full lowercase copying.
- Identified 39+ comparison sites across `cli/cmd`, `cli/cmddb`, `cli/cmdclone`, `cli/cmdpull`, `cli/cmdssh`, `cli/vscodepm`, `cli/release`, and `cli/osclean`.
- **Solution**: Refactor case-insensitive string equality comparisons to use zero-alloc `strings.EqualFold(a, b)`.

### Findings on Item 3 (OS-Aware Path Case Sensitivity Bug)
- **Status Confirmed**:
  - `fsutil.EqualPaths(p1, p2)` in `cli/fsutil/path_normalize.go` previously did `==` without case-folding, which broke on Windows where drive letters/paths have mixed casing.
  - `CanonicalRepoPathKey` in `cli/cmdpull/pull.go` and `cli/cmdpull/pull_efficient.go` unconditionally ran `strings.ToLower(path)`, which broke on Unix (Linux/macOS) where paths differing by case are distinct directories.
- **Solution**:
  1. Upgrade `cli/fsutil/path_normalize.go` with OS-aware helpers:
     - `IsPathCaseInsensitive(p string) bool`: returns `true` if `runtime.GOOS == "windows"` or path has Windows drive prefix.
     - `EqualPaths(p1, p2 string) bool`: uses `strings.EqualFold` on Windows, exact `==` on Unix/Linux.
     - `CanonicalPathKey(p string) string`: returns lowercase on Windows, preserves casing on Unix.
  2. Update callers (`cli/cmdpull/pull.go`, `cli/cmdpull/pull_efficient.go`, `cli/cmd/reconcile_db.go`, `cli/fsutil/workdir_match.go`) to use `fsutil.CanonicalPathKey` and `fsutil.EqualPaths`.

### Findings on Item 5 (Future Uniqueness Guarantee in SQLite)
- In `RepoFile` (repodb) and `Repo` (`gitmap.db`), ensure unique constraints and upsert queries prevent duplicate records.
- In `cli/indexer/walker.go` and `cli/repodb/repo_file.go`, enforce `ON CONFLICT` and deduplication guards before insert.

---

## 2. Work Breakdown & Subtask Ownership

| Task ID | Title | Owner | Target Files |
|---|---|---|---|
| `Task-01` | SQLite Redundant Repo Deduplication & Pull Integration | Worker 01 | `cli/store/repo_duplicates.go`, `cli/store/repo.go`, `cli/cmdpull/pull_dedup.go`, `cli/cmd/find_duplicates_git.go` |
| `Task-02` | Refactor String Comparisons to `strings.EqualFold` | Worker 02 | `cli/cmd/*.go`, `cli/cmdclone/*.go`, `cli/cmdssh/*.go`, `cli/release/*.go`, `cli/osclean/*.go` |
| `Task-03` | OS-Aware Path Sensitivity & `cli/fsutil` Normalization | Worker 01 | `cli/fsutil/path_normalize.go`, `cli/fsutil/workdir_match.go`, `cli/cmdpull/pull.go`, `cli/cmdpull/pull_efficient.go` |
| `Task-04` | SQLite Ingestion Uniqueness Guarantee for Files & Repos | Worker 02 | `cli/repodb/repo_file.go`, `cli/indexer/walker.go`, `cli/store/repo_upsert.go` |
| `Task-05` | Verification, Documentation, and Release | Lead Orchestrator | `02-spec/21-app/...`, `.ai-memory/...`, `version.json` |
