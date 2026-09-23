# Root Cause Analysis (RCA-18): SSH Exec IP Command Delegation Failure and Windows 'sh' Not Found

## 1. Symptom
When executing `gitmap ssh exec ip` (or related remote fleet commands) across active Windows nodes (`w1`, `w2`, `w3`):
```text
PS repo> gitmap ssh exec ip

  Notice: The following machine(s) are currently OFF or unreachable:
    • [w4 | 192.168.1.13] (machine is off)
    • [u1 | 192.168.1.22] (machine is off)

  ● Commands injected for 3 machine(s) [w1:192.168.1.3, w3:192.168.1.12, w2:192.168.1.7]:
    Command(s): sh -c ip -br a 2>/dev/null || ip a 2>/dev/null || hostname -I 2>/dev/null || ifconfig
    Processing execution across active node(s)...

  ─── [w2 | 192.168.1.7] ───
    Execute error: command execution failed: Process exited with status 1 (output: 'sh' is not recognized as an internal or external command,
operable program or batch file.
)
    'sh' is not recognized as an internal or external command,
    operable program or batch file.

  ─── [w1 | 192.168.1.3] ───
    Execute error: command execution failed: Process exited with status 1 (output: 'sh' is not recognized as an internal or external command,
operable program or batch file.
)
    'sh' is not recognized as an internal or external command,
    operable program or batch file.

  ─── [w3 | 192.168.1.12] ───
    Execute error: command execution failed: Process exited with status 1 (output: 'sh' is not recognized as an internal or external command,
operable program or batch file.
)
    'sh' is not recognized as an internal or external command,
    operable program or batch file.

  Summary of offline machines (2 machine(s) off):
    • [w4 | 192.168.1.13] (off)
    • [u1 | 192.168.1.22] (off)

  ✓ SSH Execution completed across 3 active node(s).
```

---

## 2. Root Cause
1. **Linux-Centric Argument Replacement in `resolveIPCommandArgs`:**
   In `cli/cmdssh/ssh_exec_resolve.go`, `resolveIPCommandArgs()` prematurely intercepted `args = ["ip"]` and replaced it with `[]string{"sh", "-c", "ip -br a 2>/dev/null || ip a 2>/dev/null || hostname -I 2>/dev/null || ifconfig"}`.
2. **Execution Failure on Non-Linux Targets:**
   When dispatched to remote Windows OpenSSH servers, the command invoked `sh`, which does not exist natively on Windows systems, returning `'sh' is not recognized as an internal or external command`.
3. **Bypassed GitMap Delegation & Auto-Installation:**
   Because `args` was mutated before reaching `determineSSHCommand()`, `isIPCommand(args)` evaluated to `false`, preventing delegation to GitMap (`gitmap ip`) and completely bypassing `ensureDelegateInstalled()` which is responsible for ensuring GitMap is installed on remote machines via SSH depending on the OS (PowerShell for Windows, Bash for Linux).

---

## 3. Resolution
1. **Preserve IP Command Arguments in `cli/cmdssh/ssh_exec_resolve.go`:**
   Removed hardcoded `sh -c` shell replacement from `resolveIPCommandArgs()`, returning `args` unmodified so downstream command resolvers see `"ip"`.
2. **Explicit Delegation to GitMap in `cli/cmdssh/ssh_exec_command.go`:**
   Updated `resolveIPCommand()` to return `("", "gitmap ip", true)`, designating `"ip"` as a delegated GitMap command on all platforms.
3. **Cross-Platform Auto-Install Verification:**
   Ensured that when `isDelegate` is `true`:
   - `runSSHWorker` calls `ensureDelegateInstalled(client, c)`.
   - `runSSHWorkerJSON` calls `ensureDelegateInstalledJSON(client, c, results, mu)`.
   - `executeRemoteTargetCommand` calls `ensureDelegateInstalled(client, c)`.
   If GitMap is not found on the remote machine, it automatically installs GitMap over SSH using PowerShell (`irm https://... | iex`) on Windows or Bash (`curl -fsSL https://... | bash`) on Linux/macOS before executing `gitmap ip`.
4. **Enhanced Windows Path Detection in `cli/cmdssh/ssh_exec_install.go`:**
   Updated `getGitmapCheckCmd` to verify `%LOCALAPPDATA%\gitmap-cli\gitmap.exe` in addition to PATH and legacy locations.
5. **Direct Alias Support in `cli/cmdssh/ssh.go`:**
   Updated `dispatchNodeSSH` to support `gitmap ssh ip`, `gitmap ssh ip <target>`, and `gitmap ssh <target> ip`.
6. **E2E Test Suite in `cli/tests/vm_cluster_e2e_test.go`:**
   Added `TestVMClusterE2EIPDelegation` to verify that all online cluster VMs (`w1`, `w2`, `w3`) report valid IPv4 addresses via `gitmap ip` delegation with 0 `'sh'` errors.

---

## 4. Prevention & Learnings
1. **Universal Protocol Parity:**
   Never inject shell-specific commands (`sh`, `bash`, `powershell`) into general command argument resolvers before target OS discovery has occurred.
2. **Prefer Native Go Utilities Over System Binaries:**
   Remote commands that query network state (`ip`, `status`, `health`) should always delegate to the cross-platform GitMap CLI (`gitmap ip`) rather than relying on disparate OS utilities (`ifconfig`, `ip -br a`, `Get-NetIPAddress`).
3. **Mandatory Multi-Platform E2E Verification:**
   All remote SSH execution commands must be verified against heterogeneous cluster environments (Windows OpenSSH and Linux SSH) prior to release.
