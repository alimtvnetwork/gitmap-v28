# Specification 287: Component & Golang Error Return Specification

## 1. Component Boundaries: `cli/cmdpull/`

This specification governs the error management architecture in `cli/cmdpull/`, specifically eliminating void return signatures and swallowed errors across SSH remediation and pull error recovery.

---

## 2. Functions Refactoring Blueprint

### 2.1 `updateRepoTransportInDB`
- **Location:** `cli/cmdpull/pull_remediation_ssh.go`
- **Current Signature:** `func updateRepoTransportInDB(repoPath, sshURL string)` (void, swallows DB errors)
- **Target Signature:** `func updateRepoTransportInDB(repoPath, sshURL string) error`
- **Implementation:**
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

### 2.2 `clearResolvedPullError`
- **Location:** `cli/cmdpull/pull_remediation_ssh.go`
- **Current Signature:** `func clearResolvedPullError(repoName, repoPath string)` (void, swallows DB errors)
- **Target Signature:** `func clearResolvedPullError(repoName, repoPath string) error`
- **Implementation:**
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

### 2.3 `applySSHTransportRemediation`
- **Location:** `cli/cmdpull/pull_remediation_ssh.go`
- **Current Signature:** `func applySSHTransportRemediation(repoPath string)` (void, swallows errors)
- **Target Signature:** `func applySSHTransportRemediation(repoPath string) error`
- **Implementation:**
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

### 2.4 `remediateCredentialAndConvertToSSH`
- **Location:** `cli/cmdpull/pull_remediation_ssh.go`
- **Current Signature:** `func remediateCredentialAndConvertToSSH(repoPath, repoName string)` (void, swallows exec errors)
- **Target Signature:** `func remediateCredentialAndConvertToSSH(repoPath, repoName string) error`
- **Implementation:**
  ```go
  func remediateCredentialAndConvertToSSH(repoPath, repoName string) error {
      storeType := gitutil.ResolveWindowsCredentialStore()
      errHelper := exec.Command(constants.GitBin, constants.GitConfigCmd, "--global", "credential.helper", "manager").Run()
      if errHelper != nil {
          return errHelper
      }
      errStore := exec.Command(constants.GitBin, constants.GitConfigCmd, "--global", "credential.credentialStore", string(storeType)).Run()
      if errStore != nil {
          return errStore
      }

      return applySSHTransportRemediation(repoPath)
  }
  ```

---

## 3. Caller Handling Blueprint in `cli/cmdpull/pull_remediation.go`

1. In `RemediatePullFailure`:
   ```go
   if isCredentialOrHTTPSFailure(f.ErrorText) || isSameOwnerHTTPSRepo(f.RepoPath) {
       if remErr := remediateCredentialAndConvertToSSH(f.RepoPath, f.RepoName); remErr != nil {
           fmt.Printf("    %s Warning: credential remediation failed: %v\n", constants.TreeBranch, remErr)
       }
   }
   ```
2. In `tryRetryAfterCredentialRemediation`:
   ```go
   _ = remediateCredentialAndConvertToSSH(f.RepoPath, f.RepoName)
   ```
   Replace blank discard with structured check:
   ```go
   if remErr := remediateCredentialAndConvertToSSH(f.RepoPath, f.RepoName); remErr != nil {
       fmt.Printf("    %s Warning: retry remediation encountered error: %v\n", constants.TreeBranch, remErr)
   }
   ```
3. In `clearResolvedPullError` call sites:
   Check error and output diagnostic warning if database clearing fails.

---

## 4. Coding Invariants

1. **Leaf Function Error Sovereignty**: Leaf functions returning `error` must return it directly (`return err`), never calling `cliexit.HandleError` or `os.Exit`.
2. **Positive Booleans**: Use `hasChanged`, `isSuccess`, `hasError`.
3. **Function Size**: Maximum 15 lines per function (<= 8 lines preferred).
4. **Spacing**: Blank line after `}` and blank line before `return`.
