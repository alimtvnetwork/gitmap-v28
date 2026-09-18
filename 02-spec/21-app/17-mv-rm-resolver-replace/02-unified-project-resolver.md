# Specification 17 — Chapter 2: Unified Project Resolver Engine

## 1. The Core Problem

In prior versions, commands like `gitmap rm .\prompt-architect` or `gitmap cd ./my-repo/` failed because target lookup performed direct string comparisons (`r.Slug == target` or byte-exact `r.AbsolutePath == abs`). On Windows, path separators (`\` vs `/`), trailing slashes, drive letter casing (`D:\` vs `d:\`), and relative path notations (`.\`, `..\`, `.`) created mismatches against the normalized paths in SQLite.

## 2. Resolver Resolution Hierarchy

When a user specifies a target string `T`, the resolver executes the following ordered evaluation:

1. **Working Directory Fallback (Empty or `.`)**:
   - If `T == ""` or `T == "."`, resolve the current working directory (`os.Getwd()`).
   - Look up any `Repo` whose canonical path matches PWD.
2. **Canonical Absolute Path Match**:
   - Compute `filepath.Abs(T)` and clean via Windows long-path / case-folding canonicalization.
   - Match against `Repo.AbsolutePath` in SQLite using normalized, case-insensitive comparison.
3. **Relative Path Resolution**:
   - If `T` contains path separators (`/` or `\`) or starts with `.` or `..`:
     - Clean trailing slashes: `strings.TrimRight(T, "/\\")`.
     - Resolve relative to PWD and check if directory exists on disk and in database.
4. **Alias Table Match**:
   - Query `SELECT RepoPath FROM Alias WHERE ShortName = ? COLLATE NOCASE`.
   - If matched, load `Repo` at `RepoPath`.
5. **Exact Repository Slug Match**:
   - Query `SELECT * FROM Repo WHERE Slug = ? COLLATE NOCASE`.
6. **Path Basename Match**:
   - If `T` matches `filepath.Base(r.AbsolutePath)` (e.g. `prompt-architect`), resolve the unique repo.
7. **Glob Pattern Matching**:
   - If `T` contains `*`, `?`, or `[`, expand via `filepath.Match` against both `Slug` and directory basename.

## 3. Path Normalization Contract

```go
type CanonicalPath struct {
    Original   string
    CleanPath  string
    IsAbs      bool
    ExistsDisk bool
}
```
All path comparisons within SQLite and file system operations must pass through `fsutil.NormalizePath()` or `fsutil.CanonicalRepoPath()`.
