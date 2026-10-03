# App Issue 66: OpenSSH Server Capability Installation Crash When Windows Update (wuauserv) Service is Disabled RCA

**Issue ID:** 66  
**Date:** 2026-10-03  
**Status:** Resolved  
**Affected Subsystem:** `cli/cmdssh` (`ssh_daemon_enable.go`, `ssh_daemon_enable_test.go`)  
**Reference Commit:** `02-spec/22-app-issues/66-ssh-enable-wuauserv-disabled-capability-crash-rca.md` (PR #66 / Task 66)

---

## 1. Reproduction

When executing `gitmap ssh enable` on a Windows host where OpenSSH Server is not yet installed (or capability detection falls through) and the Windows Update service (`wuauserv`) is disabled (a common configuration in hardened enterprise, developer, and cloud virtual machines):

```text
PS C:\Users\Administrator> gitmap ssh enable
  [1/4] Installing OpenSSH Server Windows Capability...
gitmap ssh: execute failed: [E9000:EXECUTION] execution: Add-WindowsCapability failed: [E9000:EXECUTION] execution: powershell failed: exit status 1 (Add-WindowsCapability : The service cannot be started, either because it is disabled or because it has no enabled devices associated
with it.
At line:1 char:1
+ Add-WindowsCapability -Online -Name OpenSSH.Server~~~~0.0.1.0
+ ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
    + CategoryInfo          : NotSpecified: (:) [Add-WindowsCapability], COMException
    + FullyQualifiedErrorId : Microsoft.Dism.Commands.AddWindowsCapabilityCommand) (at=cmdssh/ssh_daemon_enable.go:155) (at=cmdssh/ssh_daemon_enable.go:164)
Stack Trace:
    at github.com/alimtvnetwork/gitmap-v28/cli/cmdssh.runWindowsPowerShellScript (cmdssh/ssh_daemon_enable.go:164)
    at github.com/alimtvnetwork/gitmap-v28/cli/cmdssh.ensureWindowsCapability (cmdssh/ssh_daemon_enable.go:196)
    at github.com/alimtvnetwork/gitmap-v28/cli/cmdssh.enableSSHWindows (cmdssh/ssh_daemon_enable.go:205)
    at github.com/alimtvnetwork/gitmap-v28/cli/cmdssh.runEnableSSHWindows (cmdssh/ssh_daemon_enable.go:325)
    at github.com/alimtvnetwork/gitmap-v28/cli/cmdssh.dispatchSSHEnableOS (cmdssh/ssh_daemon_enable.go:346)
    at github.com/alimtvnetwork/gitmap-v28/cli/cmdssh.runSSHEnableCLI (cmdssh/ssh_daemon_enable.go:368)
    at github.com/alimtvnetwork/gitmap-v28/cli/cmdssh.RunSSHEnableCLI (cmdssh/exports.go:192)
    at github.com/alimtvnetwork/gitmap-v28/cli/cmdssh.dispatchDaemonSSH (cmdssh/ssh.go:236)
    at github.com/alimtvnetwork/gitmap-v28/cli/cmdssh.dispatchPrimarySSH (cmdssh/ssh.go:250)
    at github.com/alimtvnetwork/gitmap-v28/cli/cmdssh.dispatchSSH (cmdssh/ssh.go:293)
    at github.com/alimtvnetwork/gitmap-v28/cli/cmdssh.runSSH (cmdssh/ssh.go:21)
    at github.com/alimtvnetwork/gitmap-v28/cli/cmdssh.RunSSH (cmdssh/exports.go:26)
    at github.com/alimtvnetwork/gitmap-v28/cli/cmd.runSSH (cmd/clihelpers.go:340)
    at github.com/alimtvnetwork/gitmap-v28/cli/cmd.toolingDevEntries.func10 (cmd/roottooling.go:59)
    at github.com/alimtvnetwork/gitmap-v28/cli/cmd.runDispatchTable (cmd/rootdispatch.go:24)
    at github.com/alimtvnetwork/gitmap-v28/cli/cmd.dispatchTooling (cmd/roottooling.go:18)
    at github.com/alimtvnetwork/gitmap-v28/cli/cmd.dispatch (cmd/root.go:506)
    at github.com/alimtvnetwork/gitmap-v28/cli/cmd.runDispatch (cmd/root.go:145)
    at github.com/alimtvnetwork/gitmap-v28/cli/cmd.Run (cmd/root.go:112)
    at main.main (cli/main.go:7)
PS C:\Users\Administrator>
```

### Symptoms:
1. `gitmap ssh enable` failed with an unhandled fatal application error (`[E9000:EXECUTION]`) and printed a 20-line internal Go stack trace.
2. The root failure was a Windows DISM `COMException`: `The service cannot be started, either because it is disabled or because it has no enabled devices associated with it.`
3. Even if OpenSSH server binaries (`sshd.exe`) or the `sshd` Windows service were already present, capability detection strictly called DISM `(Get-WindowsCapability -Online -Name OpenSSH.Server~~~~0.0.1.0).State`, which failed when Windows Update services were disabled or unresponsive.
4. When installation failed, the user was left without actionable next steps or remediation guidance.

---

## 2. Root Cause Analysis (4-Part RCA)

### Symptom:
Executing `gitmap ssh enable` on Windows fails during Step 1 (`Add-WindowsCapability`) with `COMException: The service cannot be started...` and dumps a 20-line internal Go call stack.

### Direct Cause:
1. **DISM Dependency on Windows Update (`wuauserv`)**: Windows DISM and the `Add-WindowsCapability` cmdlet rely directly on the background Windows Update service (`wuauserv`) to stage and download optional features. When `wuauserv.StartType == 'Disabled'`, `Add-WindowsCapability` immediately errors with `COMException`.
2. **Fragile Single-Signal Capability Check**: `isWindowsCapabilityInstalled()` relied solely on `Get-WindowsCapability -Online`. If this command failed, errored out, or DISM was restricted by policy, GitMap incorrectly assumed OpenSSH was not installed and attempted reinstallation.
3. **Missing Temporary Service Activation & Restoration**: `buildWindowsAddCapabilityCmd()` directly invoked `Add-WindowsCapability` without checking if `wuauserv` was disabled, temporarily enabling it, and safely restoring it in a `finally` block.
4. **Unreported Error Type Causing Stack Trace**: Capability errors and service configuration errors were wrapped as raw `apperror.NewExecutionError` rather than clean `apperror.ErrorTypeAbort` with `reported: true`. This triggered `root.go`'s fallback stack trace printer.

### Compound Failure:
1. Hardened Developer Workstations: Many software engineers explicitly disable `wuauserv` to avoid unexpected background reboots or high CPU usage.
2. Degradation of Developer Trust: A command-line automation tool should either handle standard OS service quirks automatically or report clean, actionable instructions rather than dumping raw PowerShell parser errors and Go runtime internals.

### Impact:
Users on Windows machines with disabled Windows Update could not configure or start the SSH daemon via `gitmap ssh enable`.

---

## 3. Fix & Remediation

1. **Multi-Signal OpenSSH Detection (`isWindowsCapabilityInstalled`)**:
   - Added direct file inspection `isWindowsSSHDFilePresent()`: checks for `%SystemRoot%\System32\OpenSSH\sshd.exe` in 0ms without invoking PowerShell subprocesses.
   - Enhanced PowerShell capability query to check `Get-WindowsCapability`, `Get-Service -Name sshd`, and local filesystem paths.
   - Avoids unnecessary capability reinstallation when OpenSSH is already installed.

2. **Automatic `wuauserv` Lifecycle Management (`buildWindowsAddCapabilityCmd`)**:
   - Inspects `wuauserv` startup configuration prior to capability installation.
   - If `wuauserv` is disabled, temporarily sets `StartupType` to `Manual` and starts the service.
   - Wraps `Add-WindowsCapability` inside a `try / finally` block that guarantees `wuauserv` is stopped and restored back to `Disabled` upon completion or error.

3. **Clean Abort Handling Without Stack Trace**:
   - Implemented `printWindowsCapabilityFailureNotice` and `newSSHCapabilityAbortError`: presents clear, formatted diagnosis and 3 actionable remediation steps (Administrator privileges, manual `wuauserv` start, or `winget install Microsoft.OpenSSH.Beta`).
   - Uses `apperror.ErrorTypeAbort` with `reported: true` to suppress internal Go stack traces.
   - Implemented `handleSSHDServiceError` for subsequent service startup and firewall configuration steps.

4. **Decomposed Functions for Guideline Compliance**:
   - Split Windows SSH service setup into `configureSSHWindowsService` and `enableSSHWindows`, ensuring all functions remain under 15 lines.

---

## 4. Verification

1. **Unit Testing (`cli/cmdssh/ssh_daemon_enable_test.go`)**:
   - Verified capability check and installation command string builders.
   - Verified flag parsing (`--port`, `--force`, `--help`).
   - Verified error handling (`isForce`, abort error classification).
2. **Local Linter Verification**:
   - `python .github/scripts/check-legacy-refs.py .` -> Exit 0.
   - `python .github/scripts/go-format-check.py --check-only` -> Exit 0.
   - `python .github/scripts/misspell-changed.py` -> Exit 0.
   - `python .github/scripts/tests/test_ci_scripts.py` -> Exit 0.
3. **Live Execution Verification**:
   - Tested capability command on Windows with `wuauserv` disabled: service is temporarily enabled, capability is installed/verified, and service is restored back to disabled.
