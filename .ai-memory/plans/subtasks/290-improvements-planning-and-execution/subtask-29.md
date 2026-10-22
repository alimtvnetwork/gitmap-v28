# Subtask Plan 29: Vertical Whitespace & Guard Hygiene in Missing Repo Remediation (`cli/cmdpull/pull_missing_remediate.go`)

> **Subtask ID:** `subtask-29`  
> **Parent Plan:** `.ai-memory/plans/pending/290-improvements-planning-and-execution.md`  
> **Parent Spec:** `02-spec/21-app/290-improvements-planning-and-execution/03-migration-and-verification-spec.md`  
> **Target Path:** `cli/cmdpull/pull_missing_remediate.go`  
> **Status:** Completed  

---

## 1. Objectives & Scope
1. Enforce strict vertical whitespace styling in `cli/cmdpull/pull_missing_remediate.go`: exactly one blank line after `}` and exactly one blank line before `return`.
2. Flatten conditional branching using early guard clauses across missing repository remediation and interactive prompt dispatch.
3. Verify that all functions remain strictly within the <= 15 lines budget (Rule 2).
4. Preserve 100% existing functionality, terminal messaging, and interactive prompt choices.

---

## 2. Targeted Problem & Architectural Rationale
In `cli/cmdpull/pull_missing_remediate.go`, vertical whitespace inconsistencies violate the repository's coding guidelines:
1. `handlePromptDeleteAction`:
   ```go
   if err := DeleteRepoRecordByPath(f.RepoPath); err != nil {
       fmt.Printf(...)
       return true, false
   }
   fmt.Printf(...) // VIOLATION: Missing blank line after closing brace '}'!
   return true, true
   ```
2. Inconsistent blank lines before return statements across remediation helpers.
3. Chained conditional flows that can be structured with clear early guard returns.

Enforcing consistent vertical whitespace prevents merge conflicts, improves automated diff analysis, and ensures automated linters pass cleanly.

---

## 3. Concrete Implementation Details

### 3.1 Refactoring `handlePromptDeleteAction`
```go
func handlePromptDeleteAction(f PullFailureSummary) (bool, bool) {
	if err := DeleteRepoRecordByPath(f.RepoPath); err != nil {
		fmt.Printf("  %s✗%s [%s] Failed to remove record from database: %v\n",
			constants.ColorRed, constants.ColorReset, f.RepoName, err)

		return true, false
	}

	fmt.Printf("  %s✓%s [%s] Stale record removed from database.\n",
		constants.ColorGreen, constants.ColorReset, f.RepoName)

	return true, true
}
```

### 3.2 Whitespace & Guard Verification across Helpers
Review and format:
```go
func execCloneAndRecord(repoPath, remoteURL string) bool {
	if err := ExecCloneRepo(remoteURL, repoPath); err != nil {
		printCloneFailure(err)

		return false
	}

	if errTrans := remediateMissingRepoTransport(repoPath, remoteURL); errTrans != nil {
		fmt.Printf("    %s %s⚠%s Notice: cloned successfully, but failed to record transport in DB: %v\n",
			constants.TreeBranch, constants.ColorYellow, constants.ColorReset, errTrans)
	}

	fmt.Printf("    %s %s✓%s Successfully restored repository via clone fallback.\n",
		constants.TreeCorner, constants.ColorGreen, constants.ColorReset)

	return true
}

func RemediateMissingRepo(f PullFailureSummary) bool {
	printMissingRepoNotice(f)

	remoteURL := resolveMissingRemoteURL(f)
	if remoteURL == "" {
		printMissingRemoteURLError()

		return false
	}

	fmt.Printf("    %s Attempting clone fallback from %s...\n", constants.TreeBranch, remoteURL)

	return execCloneAndRecord(f.RepoPath, remoteURL)
}

func handleCloneOrFallback(choice, remoteURL string, f PullFailureSummary) (bool, bool) {
	if isCloneChoice(choice) {
		f.RemoteURL = remoteURL

		return true, RemediateMissingRepo(f)
	}

	return true, false
}
```

---

## 4. Strict Coding Guidelines Compliance
- **Rule 2 (Function Limit):** All functions <= 15 lines.
- **Rule 3 (Zero Nested Ifs):** Every conditional uses early guard returns.
- **Rule 6 (Vertical Whitespace):**
  - Exactly one blank line after every `}` block closure.
  - Exactly one blank line before every `return` statement (except when function body is single return).

---

## 5. Verification Criteria
- [x] Blank line exists after `}` in `handlePromptDeleteAction`.
- [x] Exactly one blank line precedes all `return` statements.
- [x] No nested `if` statements anywhere in `pull_missing_remediate.go`.
- [x] Every function <= 15 lines.
- [x] All file paths are strictly relative lowercase Git paths.
