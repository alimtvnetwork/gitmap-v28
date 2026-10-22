# Subtask Plan 28: Target Directory Resolution Extraction (`cli/cmdpe/pe_target.go`)

> **Subtask ID:** `subtask-28`  
> **Parent Plan:** `.ai-memory/plans/pending/291-improvements-planning-and-execution.md`  
> **Parent Spec:** `02-spec/21-app/291-improvements-planning-and-execution/03-migration-and-verification-spec.md`  
> **Target Path:** `cli/cmdpe/pe_target.go`  
> **Status:** Completed  

---

## 1. Objectives & Scope
1. Create a dedicated target resolution file `cli/cmdpe/pe_target.go`.
2. Extract common repository directory resolution and filtering helpers from `pe_gh.go` and `pe_triage.go`.
3. Provide cohesive, reusable helpers: `ResolveWorkdir(flagVal string) string`, `IsGitRepoDir(dirPath string) bool`, and `MatchRepoFilter(repoName, filter string) bool`.
4. Reduce file size and coupling in `pe_gh.go` and `pe_triage.go` while keeping all helpers <= 15 lines.

---

## 2. Targeted Problem & Architectural Rationale
Both `pe_gh.go` and `pe_triage.go` duplicate or tangle logic for:
1. Determining the effective working directory from CLI `--workdir` flags or `os.Getwd()`.
2. Inspecting whether a filesystem directory contains a valid `.git` repository folder.
3. Matching repository names against filter strings case-insensitively.

By extracting these utilities into `cli/cmdpe/pe_target.go`:
- Duplicate filesystem probing logic is eliminated.
- Files adhere to single-responsibility and keep file lengths well below limits.
- Diagnostic workflows (`pe triage`, `pe gh`, `pe rerun`) share identical target resolution behavior.

---

## 3. Concrete Implementation Details

### 3.1 New File `cli/cmdpe/pe_target.go`

```go
package cmdpe

import (
	"os"
	"path/filepath"
	"strings"
)

// ResolveWorkdir determines the target working directory from flags or current directory.
func ResolveWorkdir(flagVal string) string {
	clean := strings.TrimSpace(flagVal)
	if clean != "" {
		return clean
	}

	cwd, err := os.Getwd()
	if err != nil {
		return "."
	}

	return cwd
}

// IsGitRepoDir checks whether a directory contains a .git metadata entry.
func IsGitRepoDir(dirPath string) bool {
	gitPath := filepath.Join(dirPath, ".git")
	stat, err := os.Stat(gitPath)
	if err != nil {
		return false
	}

	return stat.IsDir()
}

// MatchRepoFilter evaluates if a repository name matches a user-specified filter.
func MatchRepoFilter(repoName, filter string) bool {
	cleanFilter := strings.TrimSpace(filter)
	if cleanFilter == "" {
		return true
	}

	return strings.Contains(strings.ToLower(repoName), strings.ToLower(cleanFilter))
}
```

### 3.2 Refactoring in `pe_gh.go` and `pe_triage.go`
- In `cli/cmdpe/pe_gh.go`, redirect internal `resolveWorkdir` and `isGitRepoDir` calls to `ResolveWorkdir` and `IsGitRepoDir`.
- In `cli/cmdpe/pe_triage.go`, redirect internal `matchRepoFilter` to `MatchRepoFilter`.

---

## 4. Strict Coding Guidelines Compliance
- **Rule 1 (Error Return Sovereignty):** Filesystem lookup errors return clean boolean or fallback strings.
- **Rule 2 (Function Limit):** `ResolveWorkdir` (12 lines), `IsGitRepoDir` (9 lines), `MatchRepoFilter` (8 lines). All <= 15 lines.
- **Rule 3 (Zero Nested Ifs):** Guard clauses on empty filter and empty string.
- **Rule 4 (Positive Booleans):** Positive naming conventions only (`IsGitRepoDir`).
- **Rule 6 (Vertical Whitespace):** Blank line before `return` and after `}`.

---

## 5. Verification Criteria
- [x] `cli/cmdpe/pe_target.go` is created with package `cmdpe`.
- [x] `ResolveWorkdir`, `IsGitRepoDir`, and `MatchRepoFilter` exported and functioning.
- [x] References in `pe_gh.go` and `pe_triage.go` updated without symbol collisions.
- [x] All functions are <= 15 lines with zero nested conditionals.
- [x] All file paths are strictly relative lowercase Git paths.
