# Spec 155: SSH Install-Exec Streaming Upload Protocol, DB Resolution Fallback & Dynamic Remote OS Probing

> **Spec ID:** `SPEC-155`  
> **Version:** `v6.333.0`  
> **Status:** Active  
> **Date:** 2026-09-25  

---

## 0. User Request (Verbatim)

```text
use

gitmap pe and fix the below installer issue by e2e test running and trying to install properly, clear???




PS C:\Users\Administrator\Downloads> gitmap ssh install-exec .\Antigravity.Manager.Tools_4.70.0_x64-setup.exe --dry-run
No matching SSH machines to deploy setup to (filtered by except: "", except-os: "", os: "").
PS C:\Users\Administrator\Downloads> gitmap ssh install-exec .\Antigravity.Manager.Tools_4.70.0_x64-setup.exe --dry-run --os win
No matching SSH machines to deploy setup to (filtered by except: "", except-os: "", os: "win").
PS C:\Users\Administrator\Downloads> gitmap ssh install-exec .\Antigravity.Manager.Tools_4.70.0_x64-setup.exe --dry-run --os win --except unix
No matching SSH machines to deploy setup to (filtered by except: "unix", except-os: "", os: "win").
PS C:\Users\Administrator\Downloads> gitmap ssh install-exec .\Antigravity.Manager.Tools_4.70.0_x64-setup.exe --dry-run --os win --except-os unix
No matching SSH machines to deploy setup to (filtered by except: "", except-os: "unix", os: "win").
PS C:\Users\Administrator\Downloads>




PS Z:\VmDownloads\04. SharedSoft> gitmap ssh install-exec .\Antigravity.Manager.Tools_4.70.0_x64-setup.exe

  ● Deploying installer 'Antigravity.Manager.Tools_4.70.0_x64-setup.exe' across 5 machine(s)

  ALIAS            HOST (IP)            OS         REMOTE PATH                RESULT           DURATION
  ----------------------------------------------------------------------------------------------------
  w1               192.168.1.3:22       linux      /tmp/Antigravity.Manager.Tools_4.70.0_x64-setup.exe ✖ UPLOAD FAILED  0ms
  w2               192.168.1.7:22       windows    C:\Windows\Temp\Antigravity.Manager.Tools_4.70.0_x64-setup.exe ✖ UPLOAD FAILED  0ms
  w3               192.168.1.12:22      windows    C:\Windows\Temp\Antigravity.Manager.Tools_4.70.0_x64-setup.exe ✖ UPLOAD FAILED  0ms
  u1               192.168.1.22:22      linux                                 ○ OFFLINE        0ms
  w4               192.168.1.13:22      windows                               ○ OFFLINE        0ms

PS Z:\VmDownloads\04. SharedSoft>


PS Z:\VmDownloads\04. SharedSoft> gitmap ssh nodes

  ALIAS            ROLE           HOST (IP:PORT)         USER           STATUS     ENROLLED
  ----------------------------------------------------------------------------------------------------
  w4               worker         192.168.1.13:22        Administrator  ● ready    2026-09-23 02:40:20
  w1               worker         192.168.1.3:22         Administrator  ● ready    2026-09-23 02:31:21
  u1               worker         192.168.1.22:22        a              ● ready    2026-09-23 02:12:22
  w3               worker         192.168.1.12:22        Administrator  ● ready    2026-09-23 02:12:22
  w2               worker         192.168.1.7:22         Administrator  ● ready    2026-09-23 02:12:22

  Total: 5 registered node(s)


PS Z:\VmDownloads\04. SharedSoft>
```

---

## 1. Visual Evidence & Reproduction

![SSH Install-Exec Failure Evidence](assets/screenshots/ssh-install-exec-failure-01.png)

### Key Observations from User Run:
1. `C:\Users\Administrator\Downloads` produced `No matching SSH machines to deploy setup to` when using `--dry-run` and OS filters because co-located or isolated databases lacked the registered fleet, or the filter logic failed on unclassified nodes.
2. In `Z:\VmDownloads\04. SharedSoft`, 5 machines were identified:
   - Node `w1` was erroneously classified as `linux` and targeted for `/tmp/...exe`, despite running Windows Server.
   - `w1`, `w2`, `w3` all failed immediately with `✖ UPLOAD FAILED  0ms`.
   - `u1` and `w4` were offline.
3. The underlying `r.Error` was completely omitted from the console results table.

---

## 2. Technical Architecture & Protocols

### 2.1 Pure SSH Channel Streaming Protocol (Replacing Inlined Base64 Argument Execution)

#### Problem:
Previous implementation converted files into a Base64 string and embedded it into an inlined command line:
`powershell -NoProfile -Command "$d=[Convert]::FromBase64String('<BASE64>'); ..."`
For an installer like `Antigravity.Manager.Tools_4.70.0_x64-setup.exe` (17.3 MB), the resulting argument string exceeds 23 million characters. Windows `CreateProcess` (32,767 character limit) and SSH channel exec packet sizes reject the command immediately, causing a 0ms `UPLOAD FAILED` error.

#### Streaming Architecture:
1. **Primary Transport — Pure Stdin Tar Stream:**
   - Client creates an SSH session and opens `stdin := session.StdinPipe()`.
   - On the remote host:
     - Windows: Executes `cmd.exe /c "tar.exe -xf - -C \"<destDir>\""` (supported natively via `bsdtar` in Windows 10/11/Server 2019/2022).
     - Unix/Linux/macOS: Executes `mkdir -p "<destDir>" && tar -xf - -C "<destDir>"`.
   - Client streams standard POSIX Tar headers and payload chunks over `stdin` directly from disk.
   - Preserves exact permissions (`0755`), file sizes, and handles arbitrary payloads (from 1 MB to 10+ GB) without memory amplification.

2. **Fallback Transport — Direct Stdin Pipe Streamer:**
   - If `tar` is absent on remote host:
     - Linux/Unix: `mkdir -p "$(dirname '<destPath>')" && cat > '<destPath>' && chmod +x '<destPath>'`
     - Windows: PowerShell stream consumer reading `[System.Console]::OpenStandardInput()` and piping to `[System.IO.File]::Create('<destPath>')`.

### 2.2 Dynamic Remote OS Probing & Persistent Registry Synchronization

1. **Active Probing upon Connect:**
   - Instead of trusting stale or default `"linux"` entries from `ssh_hosts` table, `executeInstallerOnSingleNode` invokes `probeRemoteOSType(client)`.
   - If probed OS differs from cached record (e.g. `w1` returning `windows` instead of `linux`), `c.OS` is dynamically updated.
   - The verified OS is persisted back to `SSHConnection` via `db.InsertOrUpdateSSHConnection`.
2. **OS Filter Ingestion:**
   - In `FilterSSHConnectionsByOS`, flexible tokens (`win`, `windows`, `win32`, `win64`, `linux`, `ubuntu`, `debian`, `unix`, `darwin`, `mac`, `macos`) are matched case-insensitively.

### 2.3 Global User Database Fallback for Working Directory Independence

1. If `store.BinaryDataDir()` returns a co-located `./data` directory that does not contain `gitmap.db` or has 0 registered nodes, `store.OpenDefault()` automatically falls back to the canonical system user installation path:
   - Windows: `%LOCALAPPDATA%\gitmap-cli\data\gitmap.db`
   - Linux/macOS: `~/.gitmap/data/gitmap.db`
2. Guarantees that invoking `gitmap ssh install-exec` from `C:\Users\Administrator\Downloads`, temporary folders, or arbitrary directories discovers the 5 registered cluster fleet nodes.

### 2.4 Error Surfacing in Terminal Results Table

When `r.Status` indicates a failure (`✖ UPLOAD FAILED`, `✗ FAILED (exitCode)`), `renderInstallExecResultsTable` MUST print `r.Error` under the node row so developers and users have immediate visibility into the exact failure reason.

---

## 3. Verification & Acceptance Criteria

1. **AC-SPEC155-01 (Streaming Upload):** Payloads > 15MB upload smoothly over SSH channel without command-line buffer overflow or 0ms failure.
2. **AC-SPEC155-02 (OS Probing):** Node `w1` (Windows Server) is detected as `windows`, deploying to `C:\Windows\Temp` instead of `/tmp/`.
3. **AC-SPEC155-03 (Directory Independence):** Running `gitmap ssh install-exec --dry-run` from any directory discovers all 5 nodes.
4. **AC-SPEC155-04 (Error Reporting):** Upload and execution failures display descriptive error messages in table output.
5. **AC-SPEC155-05 (Live E2E Verification):** Successful installation and verification on live online fleet node.
