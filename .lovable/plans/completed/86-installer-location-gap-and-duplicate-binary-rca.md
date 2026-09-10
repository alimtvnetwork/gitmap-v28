# Plan 86: Installer Location Alignment, Terminal Spacing Gaps, Duplicate Binary Migration & Release Orchestration (Completed)

> **Task Origin & Problem Initiation:**  
> Initiated from user session report and screenshot (`media_1789063713788.png`):
> 1. `gitmap update` update summary (`→ Source: ...`) ended with zero line gap, flush against the PowerShell command prompt.
> 2. `gitmap update` deposited binaries into `%LOCALAPPDATA%\gitmap\gitmap.exe` instead of `%LOCALAPPDATA%\gitmap-cli\gitmap.exe`.
> 3. Two conflicting binaries existed simultaneously on `PATH`, causing subsequent commands to run the stale older version (`v6.199.0`) instead of the updated version (`v6.208.1`).
> 4. User explicitly commanded: "add some gaps after the installation, it seems like the update is not performing properly it is not updating to proper location thus causing issues pease check and find the RCA and fix it and release it" and "at the end of the update it shoyld let us know simila r info like git hash, gitrepo, verson, exe location, db location etc".

---

## 1. Overview & Context

This specification resolves the installer divergence, duplicate binary conflicts, terminal output spacing, and orchestrates the requested release:
1. **Pass Active Install Directory in Remote Updater (`gitmap/cmd/updateremoteinstall.go`):**
   - Detect the running executable's directory via `os.Executable()`.
   - Pass `-InstallDir <dir>` to PowerShell installer or `--dir <dir>` to bash installer.
2. **Align Default Install Directory with Single Source of Truth (`install.ps1` & `gitmap/scripts/install.ps1`):**
   - Update `Resolve-InstallDir` to use `$script:AppSubdir` (defaulting to `gitmap-cli` per `gitmap/constants/deploy-manifest.json`).
   - Implement `Repair-LegacyLayout` to automatically detect legacy `%LOCALAPPDATA%\gitmap`, migrate data files, eliminate colliding binaries, and prune stale PATH entries.
3. **Visual Spacing & Line Gaps:**
   - Update `constants.MsgUpdateSummaryDetail` in `gitmap/constants/constants_messages.go` to append `\n\n`, providing a clean visual line gap before the shell prompt.
   - Ensure clean exit spacing across all installer summary paths.
4. **Binary Identity Summary & Commands (`gitmap binary` / `gitmap info`):**
   - Registered `gitmap binary` and `gitmap info` CLI commands to output the complete identity block (name, Git URL, version, commit SHA, database path, installed executable path, build date).
   - Integrated automatic identity block invocation at the conclusion of `gitmap update` and `install.ps1`.
5. **Active Environment Sanitization:**
   - Removed stale duplicate binary `%LOCALAPPDATA%\gitmap\gitmap.exe` on the host machine and pruned obsolete PATH references.
6. **Quality Gates & Release Execution:**
   - 100% test pass rate and 0 linter violations across all CI/CD quality gates.
   - Release orchestration via release bumper, updating `version.json`, `changelog.md`, git tags, and release assets.

---

## 2. Invariants & Rules Enforced

1. **Strict Relative Paths:** All documentation links and file paths strictly relative to the repository root. Zero absolute paths or `file:///` URIs.
2. **Coding Guidelines Adherence:** All Go functions <= 15 lines, blank line before every return, affirmative booleans (`is*`, `has*`), zero nested `if` blocks.
3. **Zero Error Swallowing:** All runtime errors wrapped in `*apperror.AppError`.
4. **Single Source of Truth:** Deployment subfolder strictly follows `gitmap/constants/deploy-manifest.json` (`gitmap-cli`).
5. **Explicit Release Authorization:** Release executed per explicit user instruction ("and release it").

---

## 3. Implementation Subtasks

### Subtask 01: Remote Updater Directory Passing, Summary Spacing & Identity Reporting
- Modified `gitmap/cmd/updateremoteinstall.go`:
  - Implemented `resolveCurrentInstallDir()` and `buildRemoteInstallerCmd(scriptPath, installDir string)`.
  - Passed `-InstallDir` (Windows) or `--dir` (POSIX) to remote installer invocation.
  - Added `printPostUpdateIdentity()` and `executeInstalledBinaryIdentity()` to invoke `binary` on the updated executable.
- Modified `gitmap/constants/constants_messages.go`:
  - Updated `MsgUpdateSummaryDetail` to end with `\n\n`.
- Modified `gitmap/cmd/rootutility.go`:
  - Registered `gitmap binary` and `gitmap info` CLI commands.
- Added unit tests in `gitmap/cmd/updateremoteinstall_test.go`.

### Subtask 02: PowerShell Installer SSoT Alignment & Legacy Layout Migration
- Updated `install.ps1` and `gitmap/scripts/install.ps1`:
  - Updated parameter documentation to state default is `$env:LOCALAPPDATA\gitmap-cli`.
  - Added `Load-DeployManifest` to load `deploy-manifest.json` from GitHub / local default.
  - Set `Resolve-InstallDir` default to `$script:AppSubdir` (`gitmap-cli`).
  - Implemented `Repair-LegacyLayout` to migrate data files from `gitmap\data` to `gitmap-cli\data`, delete duplicate `gitmap\gitmap.exe`, `gitmap\gm.exe`, `gitmap\gitmap.ps1`, and remove `gitmap` from User PATH and session `$env:PATH`.
  - Added post-install call to `& $binPath binary` for detailed identity summary.

### Subtask 03: Update-Cleanup Legacy Duplicate Removal
- Enhanced `gitmap/cmd/updatecleanup_extra.go`:
  - Implemented `cleanupLegacyDeployDir(ctx updateCleanupContext) int` to scan and remove legacy sibling `gitmap\gitmap.exe` when running inside `gitmap-cli`.
  - Registered in `gitmap/cmd/updatecleanup.go`.
- Added unit tests in `gitmap/cmd/updatecleanup_extra_test.go`.

### Subtask 04: Local Working Tree & System Sanitization
- Executed `Repair-LegacyLayout` on host machine.
- Removed stale `C:\Users\Alim\AppData\Local\gitmap\gitmap.exe` and auxiliary files.
- Pruned stale `C:\Users\Alim\AppData\Local\gitmap` from User PATH and session PATH.
- Deployed newly compiled binary to `C:\Users\Alim\AppData\Local\gitmap-cli\gitmap.exe`.
- Verified `Get-Command gitmap -All` yields only canonical `gitmap-cli` location.

### Subtask 05: Verification & Quality Gates
- Executed `go test ./...` on `cmd` and `macro` (100% pass).
- Executed `python linter-scripts/check-nested-ifs.py` (0 violations).
- Executed `python linter-scripts/check-boolean-guidelines.py` (0 violations).
- Executed `python linter-scripts/check-relative-paths.py` (0 violations).
- Executed `python linter-scripts/check-error-management.py` (0 violations).

### Subtask 06: Release Orchestration
- Synchronized version, changelog, and release assets.
