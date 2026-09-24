# 42: SSH Install-Exec 0ms Upload Failure & Remote OS Misclassification RCA

> **Issue ID:** `APP-ISSUE-42`  
> **Status:** Investigated & Resolved  
> **Date:** 2026-09-25  
> **Severity:** High  
> **Target Subsystem:** `cli/cmdssh/`, `cli/db/`, `cli/store/`  
> **Reference Spec:** [Spec 155: SSH Install-Exec Streaming Upload Protocol, DB Resolution Fallback & Dynamic Remote OS Probing](../21-app/155-ssh-install-exec-streaming-upload-and-os-resolution.md)  

---

## 1. Reproduction

When attempting to deploy a Windows setup binary (`Antigravity.Manager.Tools_4.70.0_x64-setup.exe`, 17.3 MB) across registered cluster fleet nodes:

```powershell
PS Z:\VmDownloads\04. SharedSoft> gitmap ssh install-exec .\Antigravity.Manager.Tools_4.70.0_x64-setup.exe

  ● Deploying installer 'Antigravity.Manager.Tools_4.70.0_x64-setup.exe' across 5 machine(s)

  ALIAS            HOST (IP)            OS         REMOTE PATH                RESULT           DURATION
  ----------------------------------------------------------------------------------------------------
  w1               192.168.1.3:22       linux      /tmp/Antigravity.Manager.Tools_4.70.0_x64-setup.exe ✖ UPLOAD FAILED  0ms
  w2               192.168.1.7:22       windows    C:\Windows\Temp\Antigravity.Manager.Tools_4.70.0_x64-setup.exe ✖ UPLOAD FAILED  0ms
  w3               192.168.1.12:22      windows    C:\Windows\Temp\Antigravity.Manager.Tools_4.70.0_x64-setup.exe ✖ UPLOAD FAILED  0ms
  u1               192.168.1.22:22      linux                                 ○ OFFLINE        0ms
  w4               192.168.1.13:22      windows                               ○ OFFLINE        0ms
```

Additionally, running the command from a directory like `C:\Users\Administrator\Downloads` with `--os win` returned:
`No matching SSH machines to deploy setup to (filtered by except: "", except-os: "", os: "win").`

---

## 2. Root Cause Analysis (RCA)

### RCA 1: Inlined Base64 Command String Buffer Overflow (0ms Upload Failure)
In `cli/cmdssh/ssh_install_exec.go`, file data was passed to `buildRemoteWriteCmd`:
```go
writeCmd := buildRemoteWriteCmd(remoteDestPath, data, isWindowsOS(c.OS))
shell := determineFallbackShell(c.OS)
_, errWrite := crypto.RunCommand(client, writeCmd, shell)
```
In `buildRemoteWriteCmd`:
The entire file data (17,286,424 bytes) was converted to Base64 (23,048,568 bytes) and inlined into a single command argument:
`powershell -NoProfile -Command "$d=[Convert]::FromBase64String('<23 MB BASE64 STRING>'); ..."`
In OpenSSH, channel `exec` request packets have a protocol payload limit (~32KB - 256KB), and Windows `CreateProcess` has a hard limit of 32,767 characters. The command line was rejected instantly by the OS and SSH channel, exiting in 0ms. Furthermore, `renderInstallExecResultsTable` only printed `r.Stdout`, hiding `r.Error` from the user.

### RCA 2: Unconditional `OS: "linux"` Hardcoding in SSH Connection Merging
In `cli/db/sshconnection.go`:
```go
func scanAndAppendSSHHostRow(rows *sql.Rows, seen map[string]bool, merged []SSHConnection) []SSHConnection {
    ...
    return append(merged, SSHConnection{
        Alias: alias,
        IPAddress: ip,
        Username: user,
        EncryptedPassword: encPass,
        OS: "linux", // Hardcoded fallback for any node in ssh_hosts table
    })
}
```
Because the `ssh_hosts` table lacked an explicit `OS` column, all hosts were merged with `OS: "linux"`. Node `w1` (running Windows Server 2022) was stamped as `linux`, targeted with `/tmp/...`, and ignored during `--os win` filtering. In `ssh_install_exec.go`, `probeRemoteOSType` was guarded by `if c.OS == ""`, preventing runtime correction.

### RCA 3: Lack of Global User Data Directory Fallback
When invoked from directories outside the canonical installation, `store.BinaryDataDir()` looked for `./data` relative to the executable. When running local development or unlinked binaries with empty local `./data`, zero machines were found.

---

## 3. Code Fix & Remediation

1. **Pure SSH Streaming (`StreamFileToRemote`):**
   - Replaced inlined command-line Base64 strings with native SSH session `StdinPipe()` streaming.
   - Utilizes POSIX Tar streaming (`tar.exe -xf - -C <dir>` on Windows, `tar -xf - -C <dir>` on Linux) supported by default in Windows 10/11/Server 2019/2022 (`bsdtar`) and all Unix environments.
   - Provides direct Stdin stream fallback (`cat > <destPath>` on Linux, PowerShell console stream consumer on Windows).
2. **Active Runtime Remote OS Probing:**
   - In `executeInstallerOnSingleNode`, always probe the remote OS upon connection (`probeRemoteOSType(client)`).
   - If probed OS differs from cached record (e.g. `w1` returning `windows`), update in-memory connection and write back to database via `db.InsertOrUpdateSSHConnection`.
3. **Canonical AppData / Global DB Fallback:**
   - In `store.OpenDefault()`, if the co-located database has 0 registered nodes, fallback to `%LOCALAPPDATA%\gitmap-cli\data\gitmap.db` (or `~/.gitmap/data/gitmap.db`).
4. **Transparent Failure Surfacing:**
   - Update `renderInstallExecResultsTable` to display `r.Error` whenever `r.Status` denotes failure.

---

## 4. Prevention & Quality Gates

- Added E2E tests in `cli/tests/e2e/ssh_exec_os_filter_and_install_exec_tempe2e_test.go` and `cli/cmdssh/ssh_install_test.go` to test multi-megabyte streaming transfers and verify zero command-line string inlining.
- Enforced strict remote OS probing and persistent cache refresh across fleet operations.
