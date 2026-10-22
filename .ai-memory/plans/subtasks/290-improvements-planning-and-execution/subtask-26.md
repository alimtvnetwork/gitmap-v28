# Subtask Plan 26: Tracking Remote Branch Inspection Helper (`cli/gitutil/git_cmd_exec.go`)

> **Subtask ID:** `subtask-26`  
> **Parent Plan:** `.ai-memory/plans/pending/290-improvements-planning-and-execution.md`  
> **Parent Spec:** `02-spec/21-app/290-improvements-planning-and-execution/03-migration-and-verification-spec.md`  
> **Target Path:** `cli/gitutil/git_cmd_exec.go`  
> **Status:** Completed  

---

## 1. Objectives & Scope
1. Implement `GetRepoTrackingRemote(repoPath string) (string, error)` in `cli/gitutil/git_cmd_exec.go`.
2. Inspect the upstream tracking branch configured for the current branch using `git rev-parse --abbrev-ref @{u}`.
3. Propagate errors cleanly when no tracking branch is configured or when the directory is not a git repository.
4. Adhere to Error Return Sovereignty (Rule 1) and Function Length Limits (<= 8 lines target, Rule 2).

---

## 2. Targeted Problem & Architectural Rationale
GitMap commands such as `cmdpull`, `cmdfixauth`, and CI/CD triage tools frequently need to know the specific upstream remote branch that a repository tracks (e.g. `origin/main` vs `upstream/master`).

Currently, callers must either:
1. Hardcode assumptions that the upstream branch is always `origin/<branch>`.
2. Execute custom shell commands or manual ref parsing.

Providing a centralized, safe helper in `gitutil` standardizes upstream tracking branch resolution across the entire CLI engine.

---

## 3. Concrete Implementation Details

### 3.1 Implementation in `cli/gitutil/git_cmd_exec.go`
Add `GetRepoTrackingRemote`:

```go
// GetRepoTrackingRemote returns the upstream tracking branch name for the current branch (e.g., "origin/main").
func GetRepoTrackingRemote(repoPath string) (string, error) {
	out, err := RunGitInDir(repoPath, constants.GitRevParse, constants.GitAbbrevRef, "@{u}")
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(out)), nil
}
```

### 3.2 Constants Verification
- `constants.GitRevParse = "rev-parse"`
- `constants.GitAbbrevRef = "--abbrev-ref"`
- `@{u}` passed directly as upstream tracking specifier.

---

## 4. Strict Coding Guidelines Compliance
- **Rule 1 (Error Return Sovereignty):** Subprocess error is returned directly via `return "", err`. No exit handlers or log dumps.
- **Rule 2 (Function Limit):** 8 lines total (matches the <= 8 line preferred target).
- **Rule 3 (Zero Nested Ifs):** Flat guard clause on error.
- **Rule 6 (Vertical Whitespace):** Single blank line after `}` and before final `return`.

---

## 5. Verification Criteria
- [x] Exported function `GetRepoTrackingRemote(repoPath string) (string, error)` exists in `cli/gitutil/git_cmd_exec.go`.
- [x] Executes `RunGitInDir` with `constants.GitRevParse`, `constants.GitAbbrevRef`, and `"@{u}"`.
- [x] Returns trimmed tracking branch string or error.
- [x] Exactly 8 lines of code.
- [x] All file paths are strictly relative lowercase Git paths.
