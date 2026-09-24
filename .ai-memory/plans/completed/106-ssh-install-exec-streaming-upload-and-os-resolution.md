# Completed Plan 106: SSH Install-Exec Streaming Upload Protocol, DB Resolution Fallback & Dynamic Remote OS Probing

Spec Reference: [02-spec/21-app/155-ssh-install-exec-streaming-upload-and-os-resolution.md](../../02-spec/21-app/155-ssh-install-exec-streaming-upload-and-os-resolution.md)  
Issue Reference: [02-spec/22-app-issues/42-ssh-install-exec-upload-failed-and-os-misclassification-rca.md](../../02-spec/22-app-issues/42-ssh-install-exec-upload-failed-and-os-misclassification-rca.md)  
Execution Summary: Completed in 2 orchestration loops across 3 subtask domains with verified live fleet deployment and E2E testing.

## User Request (Verbatim)

```text
gitmap pe and fix the below installer issue by e2e test running and trying to install properly, clear???

PS C:\Users\Administrator\Downloads> gitmap ssh install-exec .\Antigravity.Manager.Tools_4.70.0_x64-setup.exe --dry-run
No matching SSH machines to deploy setup to (filtered by except: "", except-os: "", os: "").
PS C:\Users\Administrator\Downloads> gitmap ssh install-exec .\Antigravity.Manager.Tools_4.70.0_x64-setup.exe --dry-run --os win
No matching SSH machines to deploy setup to (filtered by except: "", except-os: "", os: "win").
PS C:\Users\Administrator\Downloads> gitmap ssh install-exec .\Antigravity.Manager.Tools_4.70.0_x64-setup.exe --dry-run --os win --except unix
No matching SSH machines to deploy setup to (filtered by except: "unix", except-os: "", os: "win").
PS C:\Users\Administrator\Downloads> gitmap ssh install-exec .\Antigravity.Manager.Tools_4.70.0_x64-setup.exe --dry-run --os win --except-os unix
No matching SSH machines to deploy setup to (filtered by except: "", except-os: "unix", os: "win").
```

---

## Consolidated Subtasks & Outcomes

### Subtask 01: Pure SSH Stdin Streaming Protocol (`StreamFileToRemote`)
- **Target Files:** `cli/cmdssh/ssh_exec_stream.go`, `cli/cmdssh/ssh_install_exec.go`
- **Delivered:**
  - Implemented `StreamFileToRemote(client *ssh.Client, localPath, remotePath string, mode os.FileMode) error` using POSIX Tar archive streaming piped directly into remote tar extract process (`tar.exe -xf - -C <dir>` on Windows, `tar -xf - -C <dir>` on Unix/Linux).
  - Implemented fallback via raw StdinPipe streaming for hosts lacking tar utility.
  - Replaced catastrophic Base64 command-line string inlining (which produced 23 million character command strings that blew past Windows 32,767 character command limits and failed in 0ms).
  - Successfully streamed 17.3 MB installer to remote host (`w1` Windows Server 2025) and executed silent installation without error.

### Subtask 02: Dynamic Remote OS Probing & AppData Database Fallback
- **Target Files:** `cli/cmdssh/ssh_install_exec.go`, `cli/store/location.go`, `cli/db/sshconnection.go`
- **Delivered:**
  - Updated `store/location.go` to provide persistent fallback to user AppData data directory (`C:\Users\Administrator\.gitmap` or `~/.gitmap`) when GitMap is invoked outside git repositories (e.g. from `C:\Users\Administrator\Downloads`).
  - Added schema migration for `ssh_hosts` table to guarantee `OS` column persistence and avoid falling back to `"linux"`.
  - Added dynamic remote OS probing upon SSH connection in `ssh_install_exec.go`: if target node has blank or unknown OS in SQLite, it runs `which-os` probe and updates SQLite cache instantly.

### Subtask 03: Live Fleet Testing & E2E Verification
- **Target Files:** `cli/tests/e2e/ssh_install_exec_tempe2e_test.go`
- **Delivered:**
  - Authored isolated temporary E2E test verifying tar streaming and live deployment with guard tags (`//go:build tempe2e`).
  - Verified live deployment of `Antigravity.Manager.Tools_4.70.0_x64-setup.exe` to Windows node `w1` (`192.168.0.125`).
  - Surfaced rich error messages and status reporting in `renderInstallExecResultsTable`.
