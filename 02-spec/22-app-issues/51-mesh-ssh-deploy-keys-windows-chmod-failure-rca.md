# RCA 51: Mesh SSH Deploy Keys Windows Chmod Command Failure & Misleading Output

- **Issue ID:** APP-ISSUE-051
- **Severity:** High
- **Category:** Remote SSH Fleet & OS Compatibility
- **Status:** Resolved
- **Date:** 2026-09-29
- **Affected Commands:** `gitmap ssh deploy keys [all]`, `gitmap ssh keys deploy`, `gitmap machine`

---

## 1. Problem Statement & Reproduction

### 1.1 Observed Output
When running `gitmap ssh deploy keys all` against a mixed or Windows fleet:
```text
PS C:\Users\Administrator> gitmap ssh deploy keys all

  🔑 Mesh SSH Public Key Deployment
    • Unique Keys Identified: 1 (local: 2, remote: 0)
    • Target Fleet Nodes:     6 (succeeded: 0, except="")

    #    ALIAS                IP ADDRESS       STATUS       KEYS ADDED
    ──────────────────────────────────────────────────────────────────────
    1    w3                   node-w3     failed already synced
    2    w1                   node-w1      failed already synced
    3    w2                   node-w2      offline +0 key(s)
    4    w4                   node-w4     offline +0 key(s)
    5    u1                   node-u1     offline +0 key(s)
    6    main                 node-main     failed already synced

  ✓ Mesh public key synchronization complete! All nodes now trust cluster keys passwordlessly.
```

### 1.2 Observed Symptoms
1. **Misleading Status Column:** Every online Windows node displayed `failed already synced` smashed into one column.
2. **0 Remote Keys Collected:** `remote: 0` was reported despite remote nodes having valid SSH public keys.
3. **0 Nodes Succeeded:** `succeeded: 0` despite the footer declaring:
   `✓ Mesh public key synchronization complete! All nodes now trust cluster keys passwordlessly.`
4. **No Error Persistence:** Running `gitmap errors` or `gitmap e` yielded zero records because errors were never logged to `store.LogInternalErrorRecord` / `gitmap-errors.db`.

---

## 2. Root Cause Analysis (4 Whys)

### Why 1: Why did the command fail on online nodes `w3`, `w1`, and `main`?
`syncAuthorizedKeysOnClient` in `cli/cmdssh/ssh_deploy_keys_apply.go` executed a hardcoded POSIX shell command:
`prepCmd := "mkdir -p ~/.ssh && chmod 700 ~/.ssh && touch ~/.ssh/authorized_keys && chmod 600 ~/.ssh/authorized_keys && cat ~/.ssh/authorized_keys"`
On Windows OpenSSH targets, neither `chmod` nor `touch` exist. The Windows shell exited with status 1:
`'chmod' is not recognized as an internal or external command, operable program or batch file.`

### Why 2: Why were 0 remote public keys identified (`remote: 0`)?
`gatherRemotePublicKeys` in `cli/cmdssh/ssh_deploy_keys_gather.go` executed:
`cmd := "cat ~/.ssh/id_ed25519.pub ~/.ssh/id_rsa.pub ~/.ssh/id_ecdsa.pub ~/.ssh/*.pub 2>/dev/null"`
On Windows, `2>/dev/null` is invalid redirection syntax in PowerShell and cmd.exe, returning:
`The system cannot find the path specified.`
Because errors were silently skipped with `continue`, 0 remote keys were collected.

### Why 3: Why did the table print `failed already synced`?
In `cli/cmdssh/ssh_deploy_keys_render.go`:
```go
statusStr := formatNodeStatus(r)
addedStr := fmt.Sprintf("+%d key(s)", r.KeysAdded)
if r.KeysAdded == 0 && r.IsOnline {
    addedStr = "already synced"
}
```
`addedStr` did not check `r.ErrorMsg == ""`. When execution failed with 0 keys added, `statusStr` was `"failed"` and `addedStr` was `"already synced"`. Furthermore, ANSI color escape codes embedded directly in `statusStr` broke the column padding, concatenating the two fields.

### Why 4: Why was success printed when 0 nodes succeeded?
`printDeployKeysFooter` unconditionally printed the success message without evaluating `s.NodesSucceeded > 0 && len(failedNodes) == 0`.

---

## 3. Corrective Actions & Implementation

1. **OS-Adaptive Remote Key Discovery:**
   - In `cli/cmdssh/ssh_deploy_keys_gather.go`, detect `isWindowsOS(c.OS)`.
   - On Windows, read public keys using Windows-native commands (`type %USERPROFILE%\.ssh\*.pub` or PowerShell `Get-Content`).
   - On Unix, read using `cat ~/.ssh/*.pub 2>/dev/null`.

2. **OS-Adaptive Remote Key Application:**
   - In `cli/cmdssh/ssh_deploy_keys_apply.go`, detect `isWindowsOS(c.OS)`.
   - On Windows:
     - Read existing keys from `%USERPROFILE%\.ssh\authorized_keys` and `$env:ProgramData\ssh\administrators_authorized_keys`.
     - Append missing keys safely using PowerShell script (`buildWindowsAuthKeyScript` semantics), applying correct ACLs via `icacls`.
   - On Unix:
     - Use POSIX `mkdir -p ~/.ssh && chmod 700 ~/.ssh && chmod 600 ~/.ssh/authorized_keys`.

3. **Error Logging to `gitmap-errors.db` (`gitmap e` / `gitmap error`):**
   - When a node fails during `deployKeysToRemoteNode`, record an `InternalErrorRecord` via `store.LogInternalErrorRecord` detailing:
     - `Command`: `gitmap ssh deploy keys`
     - `ErrorCode`: `E_SSH_DEPLOY_KEYS_FAILED`
     - `ErrorType`: `SSH_EXECUTION`
     - `Message`: node-specific failure summary
     - `Details`: full remote output and error trace
   - Users can now inspect errors at any time via `gitmap e` or `gitmap errors`.

4. **Correct Table & Summary Rendering:**
   - In `cli/cmdssh/ssh_deploy_keys_render.go`:
     - If `r.ErrorMsg != ""`, print `error: <reason>`, NOT `"already synced"`.
     - Strip ANSI codes when calculating column widths to preserve alignment.
     - Only show green success footer if `s.NodesSucceeded > 0 && len(failed) == 0`. If all or some failed, emit a warning/error banner and advise running `gitmap e` for diagnostics.

---

## 4. Prevention & Quality Invariants

- **Invariant 1:** SSH deployment routines must NEVER assume POSIX utilities exist on target nodes without checking `isWindowsOS(node.OS)`.
- **Invariant 2:** Table row formatters must verify error status before printing informational state (e.g. "already synced").
- **Invariant 3:** All remote node failures during fleet operations must be recorded to the central error database for `gitmap e` observability.
