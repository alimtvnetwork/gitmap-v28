# Master Execution Plan: 225-ubuntu-cursor-start-menu-dock-position-and-profile-migration

## Task Identification
- **Slug**: `225-ubuntu-cursor-start-menu-dock-position-and-profile-migration`
- **Lead Agent**: Orchestrator (Antigravity)
- **Protocol**: `@.agents/skills/execute-parent-task-with-n-steps-v6`
- **Target Version**: `v6.487.0` (Minor Bump from `v6.486.0`)
- **Remote Fleet Node**: `u1` (Ubuntu 24.04 / Linux x86_64, user: `a`)

## Architecture Overview
This plan addresses the three core requirements specified by the user:
1. **Ubuntu Start Menu (App Grid) Integration for Cursor**:
   - Inject `'cursor.desktop': <{'position': <18>}>` into GNOME Shell's `org.gnome.shell app-picker-layout` via live D-Bus session so Cursor directly appears on the first page of the Ubuntu Show Applications / Start menu.
   - Run `update-desktop-database` on both system and user application directories.
   - Ensure `cursor.desktop` has Freedesktop categories (`Development;IDE;TextEditor;Utility;`) and execution bits (`chmod 755`).
2. **Configurable Start Menu / Dock Position in `gitmap os`**:
   - Implement `gitmap os dock [bottom|left|right|top] [--node <alias>]` (with aliases `panel`, `start-menu`).
   - If position argument is omitted, print the current dock/start-menu configuration.
   - On Linux (Ubuntu): Configure `org.gnome.shell.extensions.dash-to-dock dock-position` (`'BOTTOM'`, `'LEFT'`, `'RIGHT'`, `'TOP'`) via live D-Bus session.
   - On Windows: Query and report taskbar position.
   - Remote delegation: When `--node <alias>` is provided, delegate over SSH via `cmdssh.RunSSHExec`.
3. **Cursor Profile & AI Memory Export and Migration**:
   - Author migration script: `repo-secrets/05-scripts/sync-cursor-profile.py`.
   - Bundle Cursor user configuration from `%APPDATA%\Cursor\User` (`settings.json`, `keybindings.json`, `snippets`, `globalStorage`) and AI memory from `%USERPROFILE%\.cursor` (`agents`, `ai-tracking`, `extensions`, `plugins`, `projects`, `skills-cursor`).
   - Normalize Windows paths to POSIX `/home/a/...`.
   - Transfer archive to Ubuntu node `u1` and load into `~/.config/Cursor/User/` and `~/.cursor/`.
   - Author comprehensive notes in `repo-secrets/04-ubuntu-migration/cursor-profile-transfer-notes.md`.
4. **Live Verification, Quality Gates & Release Ceremony**:
   - Verify on `u1`: Cursor visible in `app-picker-layout`, dock position toggles cleanly between `bottom`, `left`, `right`, `top`, profile & memory loaded.
   - Verify linters: `check-relative-paths.py`, `check-nested-ifs.py`, `check-boolean-guidelines.py`, `go-format-check.py`.
   - Bump minor version to `v6.487.0`, commit with hyphen format, push, and monitor CI/CD pipeline via `gitmap pe -t`.

## Subtasks Breakdown
| ID | Subtask Name | Primary Files |
|---|---|---|
| 01 | Cursor Start Menu App Picker Layout & Desktop Refresh | `repo-secrets/05-scripts/setup-cursor-ubuntu.py`, `setup-cursor-ubuntu.sh` |
| 02 | `gitmap os dock` Position Command & Remote Delegation | `cli/cmdos/os_dock_cmd.go`, `cli/cmdos/os_dock_linux.go`, `cli/cmdos/os_dock_win.go`, `cli/cmdos/os_cmd.go` |
| 03 | Cursor Profile & Memory Export / Ubuntu Loader Script | `repo-secrets/05-scripts/sync-cursor-profile.py`, `sync-cursor-profile.sh`, `repo-secrets/04-ubuntu-migration/cursor-profile-transfer-notes.md` |
| 04 | Live Verification on `u1`, Quality Gates & Release Ceremony | `repo-secrets/04-ubuntu-migration/cursor-fleet-status.json`, `03-ai-scripts/37-bump-version.py` |

## Execution Checklist
- [x] Phase 1: Planning, Research & Spec Authoring (A=2 Subagents)
- [ ] Phase 2: Worker Execution (H=2 Subagents)
- [ ] Phase 3: Remote Node `u1` Live Verification
- [ ] Phase 4: Linter Quality Gates
- [ ] Phase 5: Minor Version Bump & Release Ceremony
