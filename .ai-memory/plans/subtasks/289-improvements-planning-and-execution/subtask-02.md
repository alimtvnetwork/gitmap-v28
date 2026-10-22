# Subtask Plan 02: Error Return Sovereignty in `clearResolvedPullError`

- **Subtask ID:** `subtask-02`
- **Subtask Name:** `02-pull-remediation-clear-error-sovereignty`
- **Parent Spec:** `02-spec/21-app/289-improvements-planning-and-execution/01-architecture-spec.md`
- **Assigned Agent Role:** Spec Author 01 / Implementation Worker 01
- **Target Files:**
  - `cli/cmdpull/pull_remediation_ssh.go`
  - `cli/cmdpull/pull_remediation.go`
  - `cli/cmdpull/pull_db_sync.go`
- **Status:** APPROVED & READY FOR IMPLEMENTATION

---

## 1. Targeted Problem & Architectural Rationale

In `cli/cmdpull/pull_remediation_ssh.go`, `clearResolvedPullError` evicts resolved error tracking entries from `PullSplitDB` when a repository is successfully remediated.

Historically, error evictions were performed as void operations or discarded with blank assignments (`_ = ...`), swallowing database lock failures and corruption errors. This caused stale error state to linger in `installation.db` and `pull.db`, leading to incorrect telemetry reporting and false positive alerts during subsequent health checks and pipeline scans.

Under **Rule 1 (Error Return Sovereignty)**:
1. `clearResolvedPullError(repoName, repoPath string) error` must return any failure encountered when opening `store.OpenPullSplitDB()` or calling `db.ClearPullErrorsForRepo(...)`.
2. Calling functions in `cli/cmdpull/pull_remediation.go` (`handleSuccessfulRemediation` and `handleSuccessfulRemediationSSH`) must check the error and surface actionable warnings if eviction fails.
3. Call sites in `cli/cmdpull/pull_db_sync.go` (`syncSingleStateToDB`) must handle errors returned by `ClearPullErrorsForRepo` instead of blind blank discards (`_ = ...`).

---

## 2. Component Boundaries & Code Signatures

### 2.1 Refactor `clearResolvedPullError`
- **File:** `cli/cmdpull/pull_remediation_ssh.go`
- **Signature:** `func clearResolvedPullError(repoName, repoPath string) error`
- **Requirements:**
  - Open split database via `store.OpenPullSplitDB()`.
  - Return error immediately if non-nil.
  - Ensure `defer db.Close()`.
  - Return `db.ClearPullErrorsForRepo(repoName, repoPath)` directly.
  - Enforce function length <= 10 lines.

```go
func clearResolvedPullError(repoName, repoPath string) error {
	db, err := store.OpenPullSplitDB()
	if err != nil {
		return err
	}

	defer db.Close()

	return db.ClearPullErrorsForRepo(repoName, repoPath)
}
```

### 2.2 Verify Caller Handling in `cli/cmdpull/pull_remediation.go`
- **Files:** `cli/cmdpull/pull_remediation.go`
- **Functions:**
  - `handleSuccessfulRemediation(f PullFailureSummary)`:
    ```go
    func handleSuccessfulRemediation(f PullFailureSummary) {
    	errClear := clearResolvedPullError(f.RepoName, f.RepoPath)
    	if errClear != nil {
    		fmt.Printf("    %s %s⚠%s Warning clearing pull error: %v\n",
    			constants.TreeBranch, constants.ColorYellow, constants.ColorReset, errClear)
    	}

    	fmt.Printf("    %s %s✓%s Successfully updated.\n",
    		constants.TreeCorner, constants.ColorGreen, constants.ColorReset)
    }
    ```
  - `handleSuccessfulRemediationSSH(f PullFailureSummary)`:
    ```go
    func handleSuccessfulRemediationSSH(f PullFailureSummary) {
    	errClear := clearResolvedPullError(f.RepoName, f.RepoPath)
    	if errClear != nil {
    		fmt.Printf("    %s %s⚠%s Warning clearing pull error: %v\n",
    			constants.TreeBranch, constants.ColorYellow, constants.ColorReset, errClear)
    	}

    	fmt.Printf("    %s %s✓%s Successfully updated via SSH.\n",
    		constants.TreeCorner, constants.ColorGreen, constants.ColorReset)
    }
    ```

### 2.3 Refactor DB Sync Error Handling in `cli/cmdpull/pull_db_sync.go`
- **File:** `cli/cmdpull/pull_db_sync.go`
- **Signature:** `func syncSingleStateToDB(db *store.PullSplitDB, state *PullRepoState, nodeID string, now time.Time) error`
- **Refactoring:**
  - Replace `_ = db.ClearPullErrorsForRepo(...)` with direct return or error propagation.
  - Log diagnostic or return `err` to ensure audit logging captures database write faults.

---

## 3. Strict Coding Guidelines Compliance Checklist

- [ ] **Error Sovereignty:** `clearResolvedPullError` and callers propagate errors; zero silent discards.
- [ ] **15-Line Limit:** All refactored methods <= 12 lines.
- [ ] **Zero Nested Ifs:** Single guard clauses only.
- [ ] **Positive Booleans:** Positive predicate variables used exclusively.
- [ ] **Return Styling:** Blank lines before `return` and after closing braces.
- [ ] **Zero Git / Zero Build / Zero Test:** Strictly spec and plan authoring.
