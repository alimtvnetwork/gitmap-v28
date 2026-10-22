# Subtask Plan 04: Convert Void/Swallowed Helpers in `pe_gh.go` to Return `*appfault.AppError`

- **Subtask ID:** `subtask-04`
- **Subtask Name:** `04-cmdpe-gh-error-appfault-wrapper`
- **Parent Spec:** `02-spec/21-app/289-improvements-planning-and-execution/01-architecture-spec.md`
- **Assigned Agent Role:** Spec Author 01 / Implementation Worker 01
- **Target Files:**
  - `cli/cmdpe/pe_gh.go`
  - `cli/cmdpe/pe_triage.go`
- **Status:** Completed

---

## 1. Targeted Problem & Architectural Rationale

In `cli/cmdpe/pe_gh.go`, multiple helper routines enumerate repositories and resolve filesystem paths for pipeline error triage:

1. `listGitRepos(workdir string) []triageRepo`:
   Swallows `os.ReadDir(workdir)` failures by returning `nil` directly, obscuring missing workspace directories, permission denials, or invalid paths.
2. `appendTriageRepo(repos []triageRepo, workdir string, e os.DirEntry) []triageRepo`:
   Discards the boolean result of `githubSlugForDir(repoPath)` (`slug, _ := githubSlugForDir(repoPath)`).
3. `gitLastCommitAge(dir string) string`:
   Discards execution errors from `gitutil.RunGitInDir` and unconditionally returns `"unknown"`, preventing diagnostics when Git fails or repositories are corrupt.

Under **Rule 1 (Error Return Sovereignty)** and **Rule 2 (Function Line Limits)**:
- `listGitRepos` must return `([]triageRepo, *appfault.AppError)`.
- `pe_triage.go` must handle the returned `*appfault.AppError` and propagate it to the command runner.
- Discarded booleans and errors in `appendTriageRepo` and `gitLastCommitAge` must be replaced with explicit positive checks and structured Result wrappers.

---

## 2. Component Boundaries & Code Signatures

### 2.1 Refactor `listGitRepos` in `cli/cmdpe/pe_gh.go`
- **Signature:** `func listGitRepos(workdir string) ([]triageRepo, *appfault.AppError)`
- **Requirements:**
  - Check error returned by `os.ReadDir(workdir)`.
  - If `readErr != nil`, wrap using `apperror.WrapSimple(readErr, "pe.gh.listGitRepos")` and return immediately.
  - Sort repos deterministically.
  - Return `(repos, nil)`.
  - Max 15 lines.

```go
func listGitRepos(workdir string) ([]triageRepo, *appfault.AppError) {
	entries, readErr := os.ReadDir(workdir)
	if readErr != nil {
		return nil, apperror.WrapSimple(readErr, "pe.gh.listGitRepos: "+workdir)
	}

	var repos []triageRepo
	for _, e := range entries {
		repos = appendTriageRepo(repos, workdir, e)
	}

	sort.Slice(repos, func(i, j int) bool { return repos[i].name < repos[j].name })

	return repos, nil
}
```

### 2.2 Refactor Positive Slug Check in `appendTriageRepo`
- **File:** `cli/cmdpe/pe_gh.go`
- **Changes:**
  - Replace `slug, _ := githubSlugForDir(repoPath)` with:
    ```go
    slug, hasSlug := githubSlugForDir(repoPath)
    if !hasSlug {
    	return repos
    }
    ```
  - Enforces positive boolean check `hasSlug` and only registers valid GitHub repository checkouts.

### 2.3 Propagate Error in `cli/cmdpe/pe_triage.go`
- **File:** `cli/cmdpe/pe_triage.go`
- **Changes:**
  - In `runPeTriage`:
    ```go
    workdir := resolveWorkdir(*workdirFlag)
    repos, errRepos := listGitRepos(workdir)
    if errRepos != nil {
    	return errRepos
    }
    if len(repos) == 0 {
    	fmt.Printf("pe triage: no git repos found in %s\n", workdir)
    	return nil
    }
    ```

---

## 3. Strict Coding Guidelines Compliance Checklist

- [ ] **Error Sovereignty:** `listGitRepos` returns structured `*appfault.AppError`; callers propagate directly.
- [ ] **15-Line Limit:** Refactored functions remain strictly <= 15 lines.
- [ ] **Zero Nested Ifs:** Flattened with guard clauses.
- [ ] **Positive Booleans:** Replaces blank identifier with positive `hasSlug`.
- [ ] **Return Styling:** Standard vertical whitespace before returns.
- [ ] **Zero Git / Zero Build / Zero Test:** Strictly spec and plan authoring.
