# Subtask Plan 01: Error Return Sovereignty in `updateRepoTransportInDB`

- **Subtask ID:** `subtask-01`
- **Subtask Name:** `01-pull-remediation-ssh-error-sovereignty`
- **Parent Spec:** `02-spec/21-app/289-improvements-planning-and-execution/01-architecture-spec.md`
- **Assigned Agent Role:** Spec Author 01 / Implementation Worker 01
- **Target Files:**
  - `cli/cmdpull/pull_remediation_ssh.go`
  - `cli/cmdpull/pull_remediation.go`
- **Status:** APPROVED & READY FOR IMPLEMENTATION

---

## 1. Targeted Problem & Architectural Rationale

In `cli/cmdpull/pull_remediation_ssh.go`, functions handling database updates during remote transport remediation must strictly adhere to **Rule 1 (Error Return Sovereignty)**.

When converting repository transport from HTTPS to SSH during automated pull failure recovery:
1. `updateRepoTransportInDB(repoPath, sshURL string) error` opens the default database and calls `db.UpdateRepoTransport(repoPath, sshURL, constants.TransportSSH)`.
2. Any database connectivity failure (`store.OpenDefault()`) or execution error (`db.UpdateRepoTransport(...)`) must be returned immediately rather than swallowed.
3. Callers up the chain (`applySSHTransportRemediation` and `remediateCredentialAndConvertToSSH`) must return explicit `error` types and propagate them back to the pull remediation orchestrator in `cli/cmdpull/pull_remediation.go`.

---

## 2. Component Boundaries & Code Signatures

### 2.1 Refactor `updateRepoTransportInDB`
- **File:** `cli/cmdpull/pull_remediation_ssh.go`
- **Signature:** `func updateRepoTransportInDB(repoPath, sshURL string) error`
- **Requirements:**
  - Check error returned by `store.OpenDefault()`; return immediately if non-nil.
  - Defer `db.Close()`.
  - Return `db.UpdateRepoTransport(repoPath, sshURL, constants.TransportSSH)` directly.
  - Function length: <= 10 lines.

```go
func updateRepoTransportInDB(repoPath, sshURL string) error {
	db, err := store.OpenDefault()
	if err != nil {
		return err
	}

	defer db.Close()

	return db.UpdateRepoTransport(repoPath, sshURL, constants.TransportSSH)
}
```

### 2.2 Refactor `applySSHTransportRemediation`
- **File:** `cli/cmdpull/pull_remediation_ssh.go`
- **Signature:** `func applySSHTransportRemediation(repoPath string) error`
- **Requirements:**
  - Call `ApplyTransportFlag(repoPath, true, false)`.
  - Inspect error; return immediately if non-nil.
  - Guard with positive boolean: `if !hasChanged { return nil }`.
  - Return `updateRepoTransportInDB(repoPath, newURL)`.

```go
func applySSHTransportRemediation(repoPath string) error {
	hasChanged, oldURL, newURL, err := ApplyTransportFlag(repoPath, true, false)
	if err != nil {
		return err
	}

	if !hasChanged {
		return nil
	}

	fmt.Printf("    %s Converted remote to SSH: %s -> %s\n", constants.TreeBranch, oldURL, newURL)

	return updateRepoTransportInDB(repoPath, newURL)
}
```

### 2.3 Refactor `remediateCredentialAndConvertToSSH`
- **File:** `cli/cmdpull/pull_remediation_ssh.go`
- **Signature:** `func remediateCredentialAndConvertToSSH(repoPath, repoName string) error`
- **Requirements:**
  - Resolve store type via `gitutil.ResolveWindowsCredentialStore()`.
  - Execute `configureGlobalCredentialStore(storeType)`. If error is non-nil, return `errStore`.
  - Return `applySSHTransportRemediation(repoPath)` directly.

```go
func remediateCredentialAndConvertToSSH(repoPath, repoName string) error {
	storeType := gitutil.ResolveWindowsCredentialStore()

	errStore := configureGlobalCredentialStore(storeType)
	if errStore != nil {
		return errStore
	}

	return applySSHTransportRemediation(repoPath)
}
```

### 2.4 Verify Upstream Caller Propagation in `cli/cmdpull/pull_remediation.go`
- **Callers:**
  - `applyPreemptiveCredentialFix(f PullFailureSummary)`: Captures `errRem := remediateCredentialAndConvertToSSH(f.RepoPath, f.RepoName)` and logs a structured warning with `constants.TreeBranch` and `constants.ColorYellow`.
  - `tryRetryAfterCredentialRemediation(f PullFailureSummary)`: Captures `errRem := remediateCredentialAndConvertToSSH(...)` and emits diagnostic notice if non-nil.

---

## 3. Strict Coding Guidelines Compliance Checklist

- [ ] **Error Sovereignty:** All three functions return explicit `error` types; zero `cliexit.HandleError` calls.
- [ ] **15-Line Limit:** Every modified function is <= 12 lines.
- [ ] **Zero Nested Ifs:** All branches flattened with early returns.
- [ ] **Positive Booleans:** Uses `hasChanged` instead of negative polarity flags.
- [ ] **Return Styling:** Blank line before all return statements; blank line after closing brace `}`.
- [ ] **Zero Git / Zero Build / Zero Test:** No compilation or testing executed during planning.
