# Subtask 06: Git Credential Store Helper Error Return Sovereignty

> **Subtask ID:** Subtask-06  
> **Parent Plan:** `.ai-memory/plans/292-improvements-planning-and-execution.md`  
> **Target Path:** `.ai-memory/plans/subtasks/292-improvements-planning-and-execution/subtask-06.md`  
> **Status:** Completed  

---

## 1. Objectives & Scope

1. Eliminate void signatures and error discards (`_ = exec.Command(...).Run()`) in `cli/cmdfixauth/fixcredential.go`.
2. Convert `unsetLegacyCredentialHelper` and `applyCredentialStoreType` to return explicit `*apperror.AppError` types.
3. Update `applyWindowsCredentialFix` to verify every intermediate step and bubble any configuration errors directly up to `RunFixCredential(args []string) *apperror.AppError`.
4. Ensure audit log integrity by preventing `printFixCredentialSuccess` from falsely claiming success when `git config` operations fail.
5. Enforce strict coding guidelines across `cli/cmdfixauth/`:
   - Function limit: $\le 15$ lines per function.
   - Conditionals: Zero nested `if` statements with guard clauses.
   - Positive booleans: Variables named strictly with `is` or `has` (`isWindowsPlatform`, `hasArgs`).
   - Vertical spacing: Blank line after closing brace `}` and blank line before every `return`.
   - File length: Keep files safely below the 500-line limit.

---

## 2. Detailed Technical Plan

### 2.1 Refactor `unsetLegacyCredentialHelper` in `cli/cmdfixauth/fixcredential.go`
- Locate `unsetLegacyCredentialHelper()`:
  - Update signature from `func unsetLegacyCredentialHelper()` to `func unsetLegacyCredentialHelper() *apperror.AppError`.
  - Execute `git config --global --unset-all credential.helper`.
  - Handle exit code: In Git, unsetting a non-existent key returns exit status 5 (which is benign). If another error occurs (permission denied, corrupt config, lock file), wrap and return `apperror.WrapSimple(err, "unset credential.helper")`.
  - Return `nil` when unset succeeds or key was not present:
    ```go
    func unsetLegacyCredentialHelper() *apperror.AppError {
        cmd := exec.Command("git", "config", "--global", "--unset-all", "credential.helper")
        if err := cmd.Run(); err != nil && !isBenignUnsetExit(err) {
            return apperror.WrapSimple(err, "git config unset credential.helper")
        }

        return nil
    }
    ```

### 2.2 Refactor `applyCredentialStoreType` in `cli/cmdfixauth/fixcredential.go`
- Locate `applyCredentialStoreType(storeType constants.CredentialStoreType)`:
  - Update signature to `func applyCredentialStoreType(storeType constants.CredentialStoreType) *apperror.AppError`.
  - Execute `git config --global credential.credentialStore <storeType>`.
  - Check returned error:
    ```go
    func applyCredentialStoreType(storeType constants.CredentialStoreType) *apperror.AppError {
        cmd := exec.Command("git", "config", "--global", "credential.credentialStore", string(storeType))
        if out, err := cmd.CombinedOutput(); err != nil {
            return apperror.WrapSimple(fmt.Errorf("%w: %s", err, strings.TrimSpace(string(out))), "git config set credentialStore")
        }

        return nil
    }
    ```

### 2.3 Refactor `applyWindowsCredentialFix` & Error Bubbling
- Locate `applyWindowsCredentialFix(args []string) *apperror.AppError`:
  - Enforce stepwise error checks with guard clauses:
    ```go
    if appErr := unsetLegacyCredentialHelper(); appErr != nil {
        return appErr
    }

    if appErr := setCredentialManagerHelper(); appErr != nil {
        return appErr
    }

    storeType := resolveTargetCredentialStore(args)
    if appErr := applyCredentialStoreType(storeType); appErr != nil {
        return appErr
    }

    printFixCredentialSuccess(storeType)

    return nil
    ```
- Verify that `RunFixCredential` propagates `applyWindowsCredentialFix` return directly to the caller.

---

## 3. Verification Criteria

- [x] All `_ = exec.Command(...).Run()` discards are eliminated from `cli/cmdfixauth/fixcredential.go`.
- [x] `unsetLegacyCredentialHelper` and `applyCredentialStoreType` return `*apperror.AppError`.
- [x] `applyWindowsCredentialFix` halts immediately upon any error and bubbles the `*apperror.AppError`.
- [x] `printFixCredentialSuccess` is only called if all Git configuration commands succeed.
- [x] All functions in `cli/cmdfixauth/` satisfy the $\le 15$ lines limit.
- [x] Zero nested `if` statements throughout.
- [x] Mandatory blank lines maintained after `}` and before `return`.
- [x] Strictly relative Git paths used throughout all references.
