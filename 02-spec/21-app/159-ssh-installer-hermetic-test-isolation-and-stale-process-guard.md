# Spec 159: SSH Fleet Installer Hermetic Test Isolation & Stale Process Lock Guard

## User Request (Verbatim)

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

## System Architecture & Invariants

1. **Hermetic Test Isolation (`main_test.go`):**
   - Unit tests executing within `cli/cmdssh` must never read, mutate, or purge the central user SQLite database located at `C:\Users\Administrator\AppData\Local\gitmap-cli\data\gitmap.db`.
   - `TestMain` establishes an ephemeral temporary directory (`cmdssh_hermetic_*`) and hooks `openSSHDBFunc` and `openSSHDB` to isolated SQLite files for the entire test binary lifetime.

2. **Stale Process Termination Before Binary Streaming:**
   - When rerunning installers or recovering from interrupted non-interactive sessions, Windows locks executable files in `C:\Windows\Temp\<name>.exe`.
   - Before streaming payload bytes via `StreamFileToRemote`, `killStaleInstallerProcess` issues non-blocking process termination (`taskkill /F /IM "<base>*" 2>nul` on Windows, `pkill -9 -f "<fileName>"` on Unix) to release file handles.

3. **Guaranteed Silent Switch Injection:**
   - For unattended OpenSSH sessions, `BuildRemoteInstallerExecCmdWithPayload` unconditionally ensures appropriate silent switches are present in the command string (`/S` for NSIS, `/VERYSILENT` for Inno Setup, `/qn` for MSI), even if positional arguments or custom switches are supplied.

4. **Remote Directory Bootstrap:**
   - Remote directory creation for Windows targets relies on native `cmd.exe /c if not exist "<dir>" mkdir "<dir>"` to avoid PowerShell quote-stripping parser errors.
