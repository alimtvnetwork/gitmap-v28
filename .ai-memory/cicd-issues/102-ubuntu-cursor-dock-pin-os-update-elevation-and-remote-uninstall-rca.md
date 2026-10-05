# RCA-102: Root Cause Analysis - Ubuntu Cursor Desktop Integration, OS Update Auto-Elevation & Remote Dev Tools Uninstallation

## 1. Executive Summary
- **Task Identification**: `224-ubuntu-cursor-dock-pin-projects-os-update-and-remote-uninstall`
- **Release Version**: `v6.486.0` (Commit: `aac37a97`)
- **Fleet Machine**: `u1` (Ubuntu 24.04 LTS / Linux x86_64, user: `a`)
- **Impacted Subsystems**:
  * OS package manager index update engine (`cli/cmdos/`)
  * Cursor IDE desktop integration, dock pinning, and Project Manager workspace synchronization (`repo-secrets/05-scripts/setup-cursor-ubuntu.py`, `cli/cmdcursor/`)
  * Dev tools uninstallation and remote node delegation (`cli/cmd/uninstall.go`, `cli/cmdinstall/uninstall_tools.go`)

---

## 2. Part 1: Symptom & Environmental Manifestation

### Issue A: `gitmap os update` Failed with `✖ FAILED` on Ubuntu
- **Observed Behavior**:
  Running `gitmap os update` on Ubuntu node `u1` produced:
  ```text
  ▶ Updating package indices across detected package managers...
  ▶ Summary of Package Manager Results:
    • apt   : ✖ FAILED
    • snap  : ✔ OK
  ```
- **Error Diagnostic**: The command exited without any error description or actionable remediation advice.

### Issue B: Cursor IDE Missing from Ubuntu Start Menu, Dock & Projects
- **Observed Behavior**:
  After initial provisioning, Cursor was executable in terminal (`cursor --version`) but did not appear in the Ubuntu application launcher / Start menu, was not pinned to the GNOME dock (as shown in user screenshot `media_1791197057349.png`), and opened without any recognized projects from `/home/a/git-work/*`.

### Issue C: Lack of Clean Uninstallation for Dev Tools
- **Observed Behavior**:
  Developer tools (`pwsh`, `agm`, `vim`, `vscode`, `cursor`) on Ubuntu could not be removed directly from the IDE or CLI, and lacked remote SSH delegation (`--node <alias>`) for headless node maintenance.

---

## 3. Part 2: 4-Part Root Cause Analysis (RCA)

### RCA 1: Lack of Privilege Auto-Elevation in Package Index Updates
1. **Root Cause**: `cli/cmdos/os_update_engine.go` executed `exec.Command("apt-get", "update")` directly as the invoking user (`a`). On Debian/Ubuntu, updating apt lists requires root privileges to acquire `/var/lib/apt/lists/lock`. Non-root invocations exit with `status 100` (`Permission denied`).
2. **Missing Diagnostic Context**: `printUpdateSummary` discarded the `CombinedOutput()` and `err`, displaying only `✖ FAILED` without reasons or elevation hints.

### RCA 2: Desktop Launcher, Icon Asset & GNOME Dock Integration Deficiencies
1. **Desktop Entry Permission Denial**: The previous provisioning script attempted to write exclusively to `/usr/share/applications/cursor.desktop`. When executed without root privileges, write access failed and silently aborted without creating the user-level launcher at `~/.local/share/applications/cursor.desktop`.
2. **Icon Asset Mismatch**: The desktop entry specified `Icon=cursor`. However, the official Cursor AppImage packages its icon asset as `co.anysphere.cursor.png` (the 3D cube with arrow/triangle). Because this icon was not extracted or installed into system pixmaps (`/usr/share/pixmaps/`) or the user icon theme (`~/.local/share/icons/hicolor/512x512/apps/`), GNOME failed to display any icon.
3. **GNOME Shell Favorite-Apps Live Session**: Pinning applications to the Ubuntu dock requires updating the GSettings key `org.gnome.shell favorite-apps`. Over an SSH session, GSettings requires an explicit `DBUS_SESSION_BUS_ADDRESS` (at `/run/user/<uid>/bus`). Without this bus address, GSettings commands cannot communicate with the live GNOME Shell session.
4. **Project Registry Isolation**: Cursor uses the VS Code Project Manager extension (`alefragnani.project-manager`), reading workspaces from `~/.config/Cursor/User/globalStorage/alefragnani.project-manager/projects.json`. The repositories residing in `/home/a/git-work/*` were not enumerated or injected into this registry.

### RCA 3: Missing CLI Uninstallation Command Handlers & Remote Delegation
1. **Tool Coverage Gap**: `gitmap uninstall` handled self-uninstallation, Windows apps, and generic package managers, but lacked specialized lifecycle handlers for `pwsh`, `agm`, `vim`, `vscode`, and `cursor`.
2. **Remote Node Flag Reordering Hazard**: When `--node <alias>` was supplied after the tool name (`gitmap uninstall vim --node u1`), `reorderFlagsBeforeArgs` did not recognize `--node` as a value-consuming flag, separating `--node` from `u1` and preventing SSH delegation.

---

## 4. Part 3: Integrated Technical Solutions

### Solution 1: Linux `sudo -n` Auto-Elevation & Error Diagnostics
- Updated `cli/cmdos/os_update_types.go` and `cli/cmdos/os_update_engine.go` to mark `NeedsSudo: true` for `apt`, `dnf`, and `pacman`.
- Added `isElevationRequired(tc UpdateToolchain) bool`: detects `runtime.GOOS == "linux" && os.Geteuid() != 0 && hasSudoBinary()`. When true, transparently prepends `sudo -n` to the command execution.
- Refactored `printUpdateSummary` in `cli/cmdos/os_update_cmd.go` to print diagnostic output:
  ```text
  ▶ Summary of Package Manager Results:
    • apt   : ✖ FAILED
      ↳ Error: exit status 100
      ↳ Hint: Package manager requires root. Run with sudo or configure passwordless sudo.
  ```

### Solution 2: Complete Cursor Desktop, Icon, Dock & Project Synchronization
- Updated `repo-secrets/05-scripts/setup-cursor-ubuntu.py`:
  1. **Icon Deployment**: Extracted official `co.anysphere.cursor.png` from `/opt/cursor/squashfs-root/` and deployed to:
     - `/usr/share/pixmaps/co.anysphere.cursor.png` & `/usr/share/pixmaps/cursor.png`
     - `~/.local/share/icons/hicolor/512x512/apps/co.anysphere.cursor.png` & `cursor.png`
     - `~/.local/share/pixmaps/cursor.png`
     - Invoked `gtk-update-icon-cache` to refresh desktop caches.
  2. **Dual Desktop Launcher Installation**: Installed standardized `cursor.desktop` (with `StartupWMClass=Cursor` and `Actions=new-empty-window`) to both `/usr/share/applications/cursor.desktop` and `~/.local/share/applications/cursor.desktop`.
  3. **Live GNOME Dock Pinning**: Detected the user's active D-Bus session (`/run/user/1000/bus`), retrieved current `favorite-apps`, appended `'cursor.desktop'`, and updated via `gsettings set org.gnome.shell favorite-apps`.
  4. **Workspace Synchronization**: Enumerated all 39 git repositories in `/home/a/git-work/*`, formatted Project Manager entries (`name`, `rootPath`, `paths: []`, `tags: []`, `enabled: true`), and wrote to both `~/.config/Cursor/User/globalStorage/alefragnani.project-manager/projects.json` and portable backup `~/.config/Cursor/User/projects.json`.

### Solution 3: Dev Tools Remote & Local Uninstallation Subsystem
- Added `--node <alias>` flag to `uninstallFlags` and registered `"--node": true, "-node": true` in `cli/cmd/releaseargs.go:knownValueFlags`.
- In `cli/cmd/uninstall.go`: intercepted `flags.node != ""` to delegate over SSH:
  ```go
  fmt.Printf("● Delegating uninstall to remote node '%s'...\n", flags.node)
  return cmdssh.RunSSHExec([]string{flags.node, remoteCmd})
  ```
- Implemented `cli/cmdinstall/uninstall_tools.go`:
  * Supports `pwsh`, `agm`, `vim`, `vscode`, `cursor`.
  * Fully supports `--dry-run` to print simulated commands without mutating packages or files.
  * Safety guards prevent any modification to user workspaces (`/home/a/git-work`, `d:/work`, current working directory).

---

## 5. Part 4: Live Verification & Proof Matrix

| Check / Requirement | Target Machine | Command Executed | Observed Result | Status |
|:---|:---|:---|:---|:---:|
| `gitmap os update` Auto-Elevation | `u1` (Ubuntu) | `gitmap ssh exec u1 "gitmap os update"` | `apt : ✔ OK`, `snap : ✔ OK` | **VERIFIED** |
| Official Cursor Icon Installed | `u1` (Ubuntu) | `ls -la /usr/share/pixmaps/co.anysphere.cursor.png` | 10316 bytes PNG asset present | **VERIFIED** |
| Desktop Launcher Installed | `u1` (Ubuntu) | `ls -la ~/.local/share/applications/cursor.desktop` | 442 bytes desktop file present | **VERIFIED** |
| GNOME Dock Pinning Active | `u1` (Ubuntu) | `gsettings get org.gnome.shell favorite-apps` | Includes `'cursor.desktop'` | **VERIFIED** |
| Cursor Projects Synchronized | `u1` (Ubuntu) | `cat ~/.config/Cursor/User/projects.json` | 39 workspace repositories synced | **VERIFIED** |
| Cursor Version Executable | `u1` (Ubuntu) | `cursor --version` | `3.23.12 (x64)` | **VERIFIED** |
| Dev Tools Uninstall Dry-Run | `u1` (Ubuntu) | `gitmap uninstall cursor --dry-run` | Prints full simulation, 0 files modified | **VERIFIED** |
| Remote Uninstall Delegation | Local -> `u1` | `gitmap uninstall vim --node u1 --dry-run` | Successfully delegates to `u1` via SSH | **VERIFIED** |
| Relative Path Compliance | Local | `python linter-scripts/check-relative-paths.py` | PASS (0 violations across 8733 files) | **VERIFIED** |
| Nested If Compliance | Local | `python linter-scripts/check-nested-ifs.py` | PASS (0 violations across 8 files) | **VERIFIED** |
| Boolean Guidelines Compliance | Local | `python linter-scripts/check-boolean-guidelines.py` | PASS (0 violations across 8 files) | **VERIFIED** |
| Go Formatting Cleanliness | Local | `python .github/scripts/go-format-check.py` | PASS (0 unformatted across 3964 files) | **VERIFIED** |
