# Master Execution Plan: 224-ubuntu-cursor-dock-pin-projects-os-update-and-remote-uninstall

## Task Identification
- **Slug**: `224-ubuntu-cursor-dock-pin-projects-os-update-and-remote-uninstall`
- **Lead Agent**: Orchestrator (Antigravity)
- **Protocol**: `@.agents/skills/execute-parent-task-with-n-steps-v6`
- **Target Version**: `v6.485.0` (Minor Bump from `v6.484.0`)
- **Remote Fleet Node**: `u1` (Ubuntu 24.04 / Linux x86_64)

## Architecture Overview
This parent plan addresses four tightly coupled requirements requested by the user:
1. **Cursor Desktop & Dock Integration on Ubuntu (`u1`)**:
   - Install official `.desktop` file at `/usr/share/applications/cursor.desktop` and `~/.local/share/applications/cursor.desktop`.
   - Deploy official 3D cube icons from AppImage extraction to `/usr/share/pixmaps/co.anysphere.cursor.png`, `/usr/share/pixmaps/cursor.png`, and `~/.local/share/icons/hicolor/512x512/apps/cursor.png`.
   - Update GNOME live dock favorite apps (`org.gnome.shell favorite-apps`) via D-Bus session to pin `cursor.desktop`.
   - Synchronize workspace projects from `/home/a/git-work/*` into Cursor's Project Manager extension (`~/.config/Cursor/User/globalStorage/alefragnani.project-manager/projects.json`).
2. **`gitmap os update` Linux Auto-Elevation & Diagnostics**:
   - Detect non-root execution on Linux (`os.Geteuid() != 0`).
   - For package managers requiring superuser privileges (`apt-get`, `dnf`, `pacman`), automatically prepend `sudo -n` when `sudo` is available.
   - Capture failure output and render detailed diagnostic error messages in `printUpdateSummary` rather than an opaque `✖ FAILED`.
3. **Remote & Local Dev Tools Uninstall Subsystem**:
   - Support `gitmap uninstall <tool> [--node <alias>] [--dry-run] [--force] [--purge]`.
   - Support tools: `pwsh` / `powershell`, `agm` / `antigravity-manager`, `vim`, `vscode` / `code`, `cursor`.
   - If `--node <alias>` is specified, delegate execution over SSH without requiring local presence.
   - Non-destructive by default; support `--dry-run` to print the exact execution steps without modifying packages or files.
4. **Remote Verification, Linter Gates & Release Ceremony**:
   - Remotely verify `gitmap os update` on `u1` to confirm `apt : ✔ OK`.
   - Remotely verify Cursor dock icon, launcher, and `projects.json` on `u1`.
   - Verify `gitmap uninstall` dry-run commands across all 5 tools.
   - Pass all repository linters (`check-relative-paths.py`, `check-nested-ifs.py`, `check-boolean-guidelines.py`, `go-format-check.py`).
   - Bump minor version to `v6.485.0` and monitor pipeline `gitmap pe -t` until 100% green.

## Subtasks Breakdown
| ID | Subtask Name | Primary Files |
|---|---|---|
| 01 | Cursor Desktop Launcher, Dock Pin & Project Sync | `repo-secrets/05-scripts/setup-cursor-ubuntu.py`, `setup-cursor-ubuntu.sh`, `cli/cmdcursor/cursor_sync.go` |
| 02 | `gitmap os update` Linux Sudo Elevation & Diagnostics | `cli/cmdos/os_update_engine.go`, `cli/cmdos/os_update_cmd.go`, `cli/cmdos/os_update_types.go` |
| 03 | Dev Tools Uninstall Subsystem with Remote Delegation | `cli/cmd/uninstall.go`, `cli/cmdinstall/uninstall_tools.go`, `cli/cmdinstall/install_packages.go` |
| 04 | Remote Node `u1` Verification & Release Ceremony | `repo-secrets/04-ubuntu-migration/cursor-fleet-status.json`, `03-ai-scripts/37-bump-version.py` |

## Execution Progress
- [x] Phase 1: Planning, Research & Spec Authoring (A=2 Subagents)
- [ ] Phase 2: Worker Execution (H=2 Subagents)
- [ ] Phase 3: Remote Node `u1` Live Verification
- [ ] Phase 4: Linter Quality Gates
- [ ] Phase 5: Minor Version Bump & Release Ceremony
