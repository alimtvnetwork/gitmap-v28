# Issue 09: Update and Install Location Divergence & Duplicate Binary Root Cause Analysis

## 1. Problem Description

When users execute `gitmap update`, the update command succeeds and reports a new version (e.g. `v6.208.1`), but running `gitmap` in subsequent terminal sessions continues executing an older stale version (e.g. `v6.199.0`).

Inspection revealed:
1. Two separate conflicting installations existed simultaneously on `PATH`:
   - `%LOCALAPPDATA%\gitmap-cli\gitmap.exe` (stale older binary)
   - `%LOCALAPPDATA%\gitmap\gitmap.exe` (newly downloaded binary)
2. In terminal sessions, the update completion message:
   `→ Source: https://.../install.ps1`
   was flush against the shell prompt (`PS D:\...>`) without any visual newline gap.

## 2. Root Cause Analysis

### RCA 1: Hardcoded Legacy Default in PowerShell Installers
- In `install.ps1` and `gitmap/scripts/install.ps1`:
  ```powershell
  function Resolve-InstallDir([string]$dir) {
      if ($dir -ne "") { return $dir }
      return Join-Path $env:LOCALAPPDATA "gitmap"
  }
  ```
- **The Defect:** `Resolve-InstallDir` defaulted to `$env:LOCALAPPDATA\gitmap` instead of `$env:LOCALAPPDATA\gitmap-cli`.
- In v3.6.0 and v3.14.0, Gitmap migrated all deployments to `gitmap-cli` to prevent Windows folder/binary naming collisions (spec `04-generic-cli/22-data-folder-deploy-and-cleanup.md`). The Single Source of Truth (`gitmap/constants/deploy-manifest.json`), Go constants (`constants.GitMapCliSubdir`), and `run.ps1` all specify `gitmap-cli`.
- However, `install.ps1` and `gitmap/scripts/install.ps1` were never updated to load `deploy-manifest.json` and defaulted to the obsolete legacy folder `gitmap`.

### RCA 2: Remote Update Runner Does Not Pass Active Directory
- In `gitmap/cmd/updateremoteinstall.go`:
  ```go
  func runRemoteInstaller(scriptPath string) error {
      cmd = exec.Command("powershell",
          "-ExecutionPolicy", "Bypass",
          "-NoProfile", "-NoLogo",
          "-File", scriptPath)
      ...
  }
  ```
- **The Defect:** `runRemoteInstaller` executed the downloaded installer script without passing the active running executable's directory (`-InstallDir` on Windows or `--dir` on Unix).
- As a consequence, `install.ps1` fell back to its hardcoded default (`AppData\Local\gitmap`), extracting the updated binary into `gitmap` and leaving the active `gitmap-cli` binary completely untouched.

### RCA 3: Missing Trailing Newline in Update Summary Message
- In `gitmap/constants/constants_messages.go`:
  ```go
  MsgUpdateSummaryDetail = "\n  ✓ Successfully updated from v%s to v%s\n  → Source: %s\n"
  ```
- **The Defect:** The string ended with a single `\n`, causing the interactive shell prompt to render immediately on the subsequent line without any visual breathing room.

## 3. Resolution & Safeguards

1. **Pass Active Install Directory in `updateremoteinstall.go`:**
   - Detect the running executable's directory via `os.Executable()`.
   - Pass `-InstallDir <dir>` (Windows) or `--dir <dir>` (POSIX) to the remote installer runner.

2. **Align Default Install Directory with Single Source of Truth:**
   - Update `install.ps1` and `gitmap/scripts/install.ps1` to load `deploy-manifest.json` (defaulting to `gitmap-cli`).
   - Ensure default target is always `%LOCALAPPDATA%\gitmap-cli`.

3. **Layout Repair & Duplicate Binary Migration:**
   - Implement `Repair-LegacyLayout` in `install.ps1` and `gitmap/scripts/install.ps1`.
   - Automatically detect legacy `%LOCALAPPDATA%\gitmap` deployments, migrate any unique data files, remove obsolete duplicate binaries, and prune stale PATH/profile entries.

4. **Add Terminal Gap:**
   - Update `constants.MsgUpdateSummaryDetail` in `constants_messages.go` to terminate with `\n\n`.

5. **Legacy Cleanup in `update-cleanup`:**
   - Enhance `gitmap update-cleanup` to identify and remove legacy colliding binaries when running inside `gitmap-cli`.
