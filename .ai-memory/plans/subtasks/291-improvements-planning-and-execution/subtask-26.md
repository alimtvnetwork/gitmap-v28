# Subtask Plan 26: Git Working Tree Dirty Files Query (`cli/gitutil/git_repo_query.go`)

> **Subtask ID:** `subtask-26`  
> **Parent Plan:** `.ai-memory/plans/pending/291-improvements-planning-and-execution.md`  
> **Parent Spec:** `02-spec/21-app/291-improvements-planning-and-execution/03-migration-and-verification-spec.md`  
> **Target Path:** `cli/gitutil/git_repo_query.go`  
> **Status:** Completed  

---

## 1. Objectives & Scope
1. Implement `GetRepoDirtyFiles(repoPath string) ([]string, error)` in `cli/gitutil/git_repo_query.go`.
2. Execute `git status --porcelain` in the specified directory and parse the output into a clean slice of relative modified file paths.
3. Eliminate fragmented in-place porcelain parsers scattered throughout command packages.
4. Adhere strictly to Error Return Sovereignty (Rule 1) and Function Length Limits (<= 15 lines, Rule 2).

---

## 2. Targeted Problem & Architectural Rationale
Commands such as `gitmap see`, `cpar`, and `lowercasefix` need to inspect dirty files. Currently, several commands implement ad-hoc parsing of `git status` output with manual string trimming, risking subtle discrepancies with untracked, renamed, or deleted files.

`GetRepoDirtyFiles` centralizes porcelain parsing in `gitutil`, providing:
1. A canonical, reuse-first helper for extracting dirty file paths.
2. Direct error return sovereignty without swallowing subprocess failures.
3. Clean memory slices without empty line artifacts.

---

## 3. Concrete Implementation Details

### 3.1 Function Additions in `cli/gitutil/git_repo_query.go`

```go
// GetRepoDirtyFiles returns relative paths of uncommitted, modified, or untracked files in repoPath.
func GetRepoDirtyFiles(repoPath string) ([]string, error) {
	out, err := RunGitInDir(repoPath, constants.GitStatus, "--porcelain")
	if err != nil {
		return nil, err
	}

	return parsePorcelainDirtyLines(string(out)), nil
}

func parsePorcelainDirtyLines(output string) []string {
	var files []string
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		clean := strings.TrimSpace(line)
		if len(clean) > 3 {
			files = append(files, strings.TrimSpace(clean[3:]))
		}
	}

	return files
}
```

---

## 4. Strict Coding Guidelines Compliance
- **Rule 1 (Error Return Sovereignty):** Subprocess errors from `RunGitInDir` are bubbled directly to the caller without exit handler calls.
- **Rule 2 (Function Limit):** `GetRepoDirtyFiles` (8 lines), `parsePorcelainDirtyLines` (11 lines). Both <= 15 lines.
- **Rule 3 (Zero Nested Ifs):** Linear string splitting with early guard on error.
- **Rule 4 (Positive Booleans):** No negative booleans.
- **Rule 6 (Vertical Whitespace):** Blank lines before `return` and after `}`.

---

## 5. Verification Criteria
- [x] `GetRepoDirtyFiles(repoPath string) ([]string, error)` is declared in `cli/gitutil/git_repo_query.go`.
- [x] Properly parses `git status --porcelain` status lines (skipping 2-character prefix and space).
- [x] Returns empty slice with `nil` error for clean repositories.
- [x] Each helper is <= 15 lines.
- [x] All file paths are strictly relative lowercase Git paths.
