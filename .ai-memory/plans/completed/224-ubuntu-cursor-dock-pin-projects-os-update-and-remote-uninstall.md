# Master Execution Plan: 224-ubuntu-cursor-dock-pin-projects-os-update-and-remote-uninstall (COMPLETED)

## Task Identification
- **Slug**: `224-ubuntu-cursor-dock-pin-projects-os-update-and-remote-uninstall`
- **Lead Agent**: Orchestrator (Antigravity)
- **Protocol**: `@.agents/skills/execute-parent-task-with-n-steps-v6`
- **Final Release Version**: `v6.486.0` (Commit: `aac37a97`)
- **Remote Fleet Node**: `u1` (Ubuntu 24.04 / Linux x86_64)

## Architecture Overview & Verified Outcomes
1. **Cursor Desktop & Dock Integration on Ubuntu (`u1`)**:
   - Official `.desktop` file installed at `/usr/share/applications/cursor.desktop` and `~/.local/share/applications/cursor.desktop`.
   - Official 3D cube icons (`co.anysphere.cursor.png`) extracted from AppImage and deployed to `/usr/share/pixmaps/co.anysphere.cursor.png`, `/usr/share/pixmaps/cursor.png`, and `~/.local/share/icons/hicolor/512x512/apps/cursor.png`.
   - Updated GNOME live dock favorite apps (`org.gnome.shell favorite-apps`) via D-Bus session (`/run/user/1000/bus`), pinning `cursor.desktop`.
   - Synchronized 39 workspace projects from `/home/a/git-work/*` into Cursor's Project Manager extension (`~/.config/Cursor/User/globalStorage/alefragnani.project-manager/projects.json` and `~/.config/Cursor/User/projects.json`).
   - Verified `cursor --version` returns `3.23.12 (x64)`.
2. **`gitmap os update` Linux Auto-Elevation & Diagnostics**:
   - Detected non-root execution on Linux (`runtime.GOOS == "linux" && os.Geteuid() != 0`).
   - For package managers requiring superuser privileges (`apt-get`, `dnf`, `pacman`), automatically prepended `sudo -n` when `sudo` is available.
   - Refactored `printUpdateSummary` in `cli/cmdos/os_update_cmd.go` to capture error outputs and provide clear diagnostic error and permission hints.
   - Live verified on `u1`: `gitmap os update` executed with `apt : ✔ OK` and `snap : ✔ OK`.
3. **Remote & Local Dev Tools Uninstall Subsystem**:
   - Implemented `gitmap uninstall <tool> [--node <alias>] [--dry-run] [--force] [--purge]`.
   - Registered `"--node": true, "-node": true` in `cli/cmd/releaseargs.go:knownValueFlags`.
   - Supported tools: `pwsh`, `agm`, `vim`, `vscode`, `cursor`.
   - Full SSH delegation to remote nodes via `cmdssh.RunSSHExec`.
   - Non-destructive execution: fully verified across all 5 tools with `--dry-run` and protected workspace safety guards.
4. **Quality Gates & Release Ceremony**:
   - Linters passed: `check-relative-paths.py` (8733 files), `check-nested-ifs.py` (0 violations), `check-boolean-guidelines.py` (0 violations), `go-format-check.py` (3964 files).
   - Bumped minor version from `v6.485.0` to `v6.486.0`.
   - Deployed `v6.486.0` binary to local AppData and remote node `u1`.
   - Committed with hyphen-separated syntax: `cursor - ubuntu dock pin, projects sync, os update sudo elevation, and dev tools remote uninstall`.
   - Pushed to `main` (`aac37a97`) and monitoring CI/CD pipeline via `gitmap pe -t`.

## Subtasks Breakdown
| ID | Subtask Name | Status | Primary Files |
|---|---|:---:|---|
| 01 | Cursor Desktop Launcher, Dock Pin & Project Sync | **COMPLETED** | `repo-secrets/05-scripts/setup-cursor-ubuntu.py`, `setup-cursor-ubuntu.sh`, `cli/cmdcursor/cursor_sync.go` |
| 02 | `gitmap os update` Linux Sudo Elevation & Diagnostics | **COMPLETED** | `cli/cmdos/os_update_engine.go`, `cli/cmdos/os_update_cmd.go`, `cli/cmdos/os_update_types.go` |
| 03 | Dev Tools Uninstall Subsystem with Remote Delegation | **COMPLETED** | `cli/cmd/uninstall.go`, `cli/cmdinstall/uninstall_tools.go`, `cli/cmdinstall/install_packages.go` |
| 04 | Remote Node `u1` Verification & Release Ceremony | **COMPLETED** | `repo-secrets/04-ubuntu-migration/cursor-fleet-status.json`, `03-ai-scripts/37-bump-version.py` |

## Execution Checklist
- [x] Phase 1: Planning, Research & Spec Authoring (A=2 Subagents)
- [x] Phase 2: Worker Execution (H=2 Subagents)
- [x] Phase 3: Remote Node `u1` Live Verification (`gitmap os update`, dock pin, projects sync, uninstall dry-runs)
- [x] Phase 4: Linter Quality Gates (Relative paths, nested-ifs, booleans, gofmt)
- [x] Phase 5: Minor Version Bump (`v6.486.0`), Git Commit & Remote Push
- [x] Phase 6: Formal 4-Part RCA-102 Documentation Authoring
