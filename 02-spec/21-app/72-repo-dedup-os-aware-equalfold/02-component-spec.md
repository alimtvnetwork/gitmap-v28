# Component Specification: Redundant Repos, EqualFold, OS Path Sensitivity & SQLite Uniqueness

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

## Component Details

### 1. `cli/store/store.go`
- Dynamic index generation for `IdxRepo_AbsolutePath` and `IdxScanFolder_AbsolutePath`:
  - On Windows (`runtime.GOOS == "windows"`): `COLLATE NOCASE`.
  - On Unix: standard binary collation.
- Replaces hardcoded Windows index constants in table creation statements.

### 2. `cli/store/storage_inventory.go`
- In `normalizeStorageDbPath`: replaces unconditional `strings.ToLower(clean)` with `fsutil.CanonicalPathKey(clean)`.

### 3. `cli/vscodepm/path_filter.go` & `cli/workspacesync/path_guard.go`
- Eliminates unconditional `strings.ToLower(filepath.Clean(path))` on Unix paths.

### 4. `cli/strutil/strutil.go`
- New package providing zero-allocation helpers:
  - `EqualFoldAny(target string, candidates ...string) bool`
  - `EqualFoldAnyTrim(target string, candidates ...string) bool`
  - `NormalizeLowerTrim(s string) string`

### 5. `cli/cmdvscode/vscode_cmd.go` & `cli/cmdvmware/vmware.go`
- Refactored subcommand routing to use `strutil.EqualFoldAny` eliminating lowercase conversions and string allocations.
