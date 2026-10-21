# Subtask 05: Headless OpenSSH Wincredman Bypass & Safe-Pull Auto-SSH Conversion

- **Subtask ID**: `05-subtask-wincredman-ssh-pull-fix`
- **Parent Task**: `gitmap-enhancement-and-issue-resolution`
- **Component Spec**: `02-spec/21-app/gitmap-enhancement-and-issue-resolution/02-component-spec.md`
- **Owned Files**:
  - `cli/gitutil/git_env.go`
  - `cli/gitutil/git_env_ssh_test.go`
  - `cli/cmdssh/ssh_pull_fleet_dispatch.go`
  - `cli/cloner/safe_pull.go`
  - `cli/cloner/safe_pull_ssh.go`
- **Status**: COMPLETED

---

## 1. Objectives & Context

During `gitmap pa --ssh` execution across fleet nodes, remote Windows node `w2` (`192.168.1.7`) failed on 14 HTTPS repositories with exit code 128:
`Windows Credential Manager (wincredman) failed to persist credentials`.

This failure stems from the Windows OpenSSH Service Logon Type 3 (Network Logon) security model, which isolates user DPAPI master keys and blocks credential persistence. Because `cli/gitutil/git_env.go` unconditionally appends `GCM_CREDENTIAL_STORE=wincredman` on Windows, headless Git subprocesses abort immediately.

This subtask resolves the issue across three defense layers:
1. **Dynamic Environment Configuration**: Detect headless SSH sessions (`SSH_CLIENT`, `SSH_CONNECTION`, `GITMAP_HEADLESS_SSH`) and configure `GCM_CREDENTIAL_STORE=none` and `GCM_NO_PERSIST=1` in `cli/gitutil/git_env.go`.
2. **Fleet Command Injection**: Inject GCM bypass environment variables in `buildRemoteFleetPullCmd` in `cli/cmdssh/ssh_pull_fleet_dispatch.go`.
3. **Automated Protocol Fallback**: Auto-convert HTTPS remotes to SSH remotes (`git@github.com:...`) upon credential failure in the retry loop of `cli/cloner/safe_pull.go`.

---

## 2. Step-by-Step Implementation Plan

### Step 2.1: Headless SSH Detection & Credential Env (`cli/gitutil/git_env.go`)
- Implement `isHeadlessSSHSession() bool` (<= 8 lines):
  - Check `len(os.Getenv("SSH_CLIENT")) > 0`.
  - Check `len(os.Getenv("SSH_CONNECTION")) > 0`.
  - Check `os.Getenv("GITMAP_HEADLESS_SSH") == "1"`.
- Refactor `AppendGitCredentialEnv(env []string) []string` (<= 15 lines):
  - On non-Windows: return append `GCM_CREDENTIAL_STORE=cache`.
  - On Windows with `isHeadlessSSHSession()`: return append `GCM_CREDENTIAL_STORE=none` and `GCM_NO_PERSIST=1`.
  - On Windows interactive: return append `GCM_CREDENTIAL_STORE=wincredman`.

### Step 2.2: Unit Test Suite for Safe Git Env (`cli/gitutil/git_env_ssh_test.go`)
- Test `AppendGitCredentialEnv` under standard interactive conditions.
- Test `AppendGitCredentialEnv` with simulated `SSH_CLIENT` environment variable.
- Verify `GCM_CREDENTIAL_STORE=none` and `GCM_NO_PERSIST=1` are present in headless mode.
- Verify `GCM_CREDENTIAL_STORE=wincredman` is NOT present in headless mode.

### Step 2.3: Remote Command Injection (`cli/cmdssh/ssh_pull_fleet_dispatch.go`)
- Refactor `buildRemoteFleetPullCmd`:
  - Check if target node OS is Windows via node metadata.
  - For Windows nodes, prepend PowerShell/CMD environment definitions:
    `$env:GCM_CREDENTIAL_STORE='none'; $env:GCM_NO_PERSIST='1'; gitmap pa --json ...`
  - For Unix/Linux nodes, prepend `GITMAP_HEADLESS_SSH=1 gitmap pa --json ...`.
- Guarantee that remote subcommands run cleanly even if remote GitMap binary is un-upgraded.

### Step 2.4: Pull Retry Auto-Conversion to SSH (`cli/cloner/safe_pull.go` & `safe_pull_ssh.go`)
- Implement `hasCredentialPersistenceFailure(output string) bool` in `cli/cloner/safe_pull.go` (<= 10 lines):
  - Check for substrings `wincredman`, `failed to persist`, or GCM exit 128.
- In `executePullAttempt`:
  - If attempt fails and `hasCredentialPersistenceFailure(output)` is true:
    - Query origin URL using `queryRemoteOriginURL(repoDir)`.
    - If URL starts with `https://`:
      - Execute `rewriteRemoteToSSHIfNeeded(repoDir, true)`.
      - Log protocol auto-conversion event via `verbose.Get()`.
      - Trigger an immediate retry using the converted SSH URL without counting against generic retry backoff limits.

---

## 3. Coding Guidelines Invariants Checklist

- [ ] Error Sovereignty: All functions return error directly; zero calls to exit handlers.
- [ ] Boolean Naming: All booleans prefixed with `is*` or `has*` only (`isHeadless`, `hasPersistenceFailure`, `isWindowsNode`).
- [ ] Function Size: All functions <= 15 lines (<= 8 lines preferred).
- [ ] Conditional Flattening: Zero nested `if` statements; early returns used throughout.
- [ ] Relative Path Hygiene: Zero absolute paths; strictly relative paths in code and tests.

---

## 4. Verification & Validation Commands

1. **Unit Tests Execution**:
   ```bash
   go test -v ./cli/gitutil/ -run TestAppendGitCredentialEnv
   go test -v ./cli/cloner/ -run TestSafePullCredentialFallback
   ```
2. **Headless SSH Simulation**:
   ```bash
   pwsh -Command "$env:SSH_CLIENT='192.168.1.3 54321 22'; gitmap pa --repo ../story-telling"
   ```
   Verify Git subprocess does not invoke `wincredman` and exits cleanly.
3. **Fleet Pull Test**:
   ```bash
   gitmap pa --ssh --limit 2
   ```
   Verify remote node `w2` pulls 67 repositories with zero `wincredman` exit 128 failures.
