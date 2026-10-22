# Subtask Plan 03: Eliminate Swallowed Error in `remediateMissingRepoTransport`

- **Subtask ID:** `subtask-03`
- **Subtask Name:** `03-pull-missing-remediate-transport-error`
- **Parent Spec:** `02-spec/21-app/289-improvements-planning-and-execution/01-architecture-spec.md`
- **Assigned Agent Role:** Spec Author 01 / Implementation Worker 01
- **Target Files:**
  - `cli/cmdpull/pull_missing_remediate.go`
- **Status:** Completed

---

## 1. Targeted Problem & Architectural Rationale

In `cli/cmdpull/pull_missing_remediate.go`, when a repository directory is missing from disk during `gitmap pull`, `RemediateMissingRepo` executes a fallback `git clone` using the recorded remote URL.

However:
1. After successfully cloning the repository back onto disk via `ExecCloneRepo`, the repository's active transport and remote URL were not synchronized with `gitmap.db`, leaving the database transport out of sync with disk state.
2. In interactive prompt workflows (`PromptMissingRepoAction`), user input errors (`reader.ReadString('\n')`) and database lookups (`LookupRepoRemoteAndID`) silently swallowed errors, failing to report database connection drops or read failures.
3. Under **Rule 1 (Error Return Sovereignty)** and **Rule 4 (Positive Booleans)**, we must introduce a dedicated `remediateMissingRepoTransport(repoPath, remoteURL string) error` helper, handle its return value explicitly, and ensure zero swallowed errors during transport recording.

---

## 2. Component Boundaries & Code Signatures

### 2.1 Add `remediateMissingRepoTransport` Helper
- **File:** `cli/cmdpull/pull_missing_remediate.go`
- **Signature:** `func remediateMissingRepoTransport(repoPath, remoteURL string) error`
- **Requirements:**
  - Open database via `store.OpenDefault()`; return `err` immediately if non-nil.
  - Ensure `defer db.Close()`.
  - Classify transport positively: `constants.TransportHTTPS` if prefix matches `constants.GitPrefixHTTPS`, else `constants.TransportSSH`.
  - Call `db.UpdateRepoTransport(repoPath, remoteURL, transport)` and return the error.
  - Function length <= 15 lines.

```go
func remediateMissingRepoTransport(repoPath, remoteURL string) error {
	db, err := store.OpenDefault()
	if err != nil {
		return err
	}

	defer db.Close()

	transport := constants.TransportSSH
	if strings.HasPrefix(strings.ToLower(remoteURL), string(constants.GitPrefixHTTPS)) {
		transport = constants.TransportHTTPS
	}

	return db.UpdateRepoTransport(repoPath, remoteURL, transport)
}
```

### 2.2 Integrate Transport Sync in `RemediateMissingRepo`
- **File:** `cli/cmdpull/pull_missing_remediate.go`
- **Changes in `RemediateMissingRepo`:**
  - After `ExecCloneRepo(remoteURL, f.RepoPath)` succeeds:
    ```go
    if errTrans := remediateMissingRepoTransport(f.RepoPath, remoteURL); errTrans != nil {
    	fmt.Printf("    %s %s⚠%s Notice: cloned successfully, but failed to record transport in DB: %v\n",
    		constants.TreeBranch, constants.ColorYellow, constants.ColorReset, errTrans)
    }
    ```
  - Eliminates silent loss of transport records upon clone recovery.

### 2.3 Refactor Input Error Handling in `PromptMissingRepoAction`
- **File:** `cli/cmdpull/pull_missing_remediate.go`
- **Changes in `PromptMissingRepoAction`:**
  - Check error returned by `reader.ReadString('\n')`:
    ```go
    line, readErr := reader.ReadString('\n')
    if readErr != nil {
    	return false, false
    }
    ```
  - Eliminates blank discard `line, _ := reader.ReadString('\n')`.

---

## 3. Strict Coding Guidelines Compliance Checklist

- [ ] **Error Sovereignty:** `remediateMissingRepoTransport` returns `error`; `readErr` is inspected and handled.
- [ ] **15-Line Limit:** Helper function is <= 14 lines.
- [ ] **Zero Nested Ifs:** Flattened conditionals with early return.
- [ ] **Positive Booleans:** Uses `hasUpdated` and `hasError` positive patterns.
- [ ] **Return Styling:** Standard vertical whitespace before and after returns.
- [ ] **Zero Git / Zero Build / Zero Test:** Strictly spec and plan authoring.
