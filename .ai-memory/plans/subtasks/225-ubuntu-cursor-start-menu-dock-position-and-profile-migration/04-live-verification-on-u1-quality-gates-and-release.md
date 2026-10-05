# Subtask Plan 04: Remote Node u1 Live Verification, Linter Quality Gates & Release Ceremony

- **Spec Reference:** [02-spec/21-app/225-ubuntu-cursor-start-menu-dock-position-and-profile-migration/02-component-and-cli-spec.md](../../../../02-spec/21-app/225-ubuntu-cursor-start-menu-dock-position-and-profile-migration/02-component-and-cli-spec.md)
- **Status:** Queued
- **Target Area:** `repo-secrets/04-ubuntu-migration/cursor-fleet-status.json`, `linter-scripts/`, `03-ai-scripts/`, `version.json`, `package.json`, `changelog.md`, `cli/constants/constants.go`, `cli/cmdpipeline/`

---

## 1. Objective

Orchestrate and execute the comprehensive live verification, quality gate enforcement, and release ceremony for task 225:
1. Conduct live remote verification across Ubuntu node `u1` via SSH to confirm GNOME Shell app grid placement (`cursor.desktop` in `app-picker-layout`), desktop database sync, `gitmap os dock` multi-orientation control, and profile/AI memory migration.
2. Validate zero violations across repository linters: relative path hygiene, nested-if limits, boolean standards, and Go code formatting.
3. Perform the official minor version bump from `6.486.0` to `6.487.0` via `03-ai-scripts/37-bump-version.py`.
4. Stage changes using GitMap's hyphen-separated atomic commit convention:
   `cursor - start menu app grid placement, os dock position configuration, and profile memory migration`
5. Monitor CI/CD pipeline health in real time via `gitmap pe -t` until 100% green.

---

## 2. Implementation Details

### Step 1: Remote Node `u1` Live Verification Workflow
Execute remote commands over SSH via `gitmap ssh exec u1 "<command>"`:
1. **Cursor Desktop Launcher & Permissions:**
   ```bash
   gitmap ssh exec u1 "test -x /usr/share/applications/cursor.desktop || test -f /usr/share/applications/cursor.desktop && echo 'LAUNCHER_EXISTS'"
   ```
   - Assert exit code 0 and output `LAUNCHER_EXISTS`.
2. **Desktop Database Cache Refresh:**
   ```bash
   gitmap ssh exec u1 "update-desktop-database ~/.local/share/applications/ && update-desktop-database /usr/share/applications/ && echo 'DESKTOP_DB_OK'"
   ```
   - Assert exit code 0 and output `DESKTOP_DB_OK`.
3. **GNOME Shell Start Menu / App Grid Placement:**
   ```bash
   gitmap ssh exec u1 "DBUS_SESSION_BUS_ADDRESS=unix:path=/run/user/1000/bus gsettings get org.gnome.shell app-picker-layout"
   ```
   - Assert stdout contains `'cursor.desktop'`.
4. **Dock Favorite-Apps Verification:**
   ```bash
   gitmap ssh exec u1 "DBUS_SESSION_BUS_ADDRESS=unix:path=/run/user/1000/bus gsettings get org.gnome.shell favorite-apps"
   ```
   - Assert stdout contains `'cursor.desktop'`.

### Step 2: `gitmap os dock` Position Control Verification on `u1`
Test dock position query and orientation mutation via the GitMap CLI:
1. **Query Current Position:**
   ```bash
   gitmap os dock --node u1
   ```
   - Assert output reports the active dock orientation.
2. **Set Orientation to Bottom:**
   ```bash
   gitmap os dock bottom --node u1
   ```
   - Assert command exits 0 and reports dock position updated to `BOTTOM`.
3. **Verify D-Bus State on `u1`:**
   ```bash
   gitmap ssh exec u1 "DBUS_SESSION_BUS_ADDRESS=unix:path=/run/user/1000/bus gsettings get org.gnome.shell.extensions.dash-to-dock dock-position"
   ```
   - Assert output equals `'BOTTOM'`.
4. **Set Orientation to Left (Standard Ubuntu Default):**
   ```bash
   gitmap os dock left --node u1
   ```
   - Assert command exits 0 and reports dock position updated to `LEFT`.

### Step 3: Profile & AI Memory Integrity Check on `u1`
Verify that migrated profile configurations and AI memory assets are active on `u1`:
1. **Settings File & Dracula Theme:**
   ```bash
   gitmap ssh exec u1 "test -f ~/.config/Cursor/User/settings.json && grep -q 'Dracula Theme' ~/.config/Cursor/User/settings.json && echo 'SETTINGS_OK'"
   ```
   - Assert exit code 0 and output `SETTINGS_OK`.
2. **Keybindings & Snippets:**
   ```bash
   gitmap ssh exec u1 "test -f ~/.config/Cursor/User/keybindings.json && test -d ~/.config/Cursor/User/snippets && echo 'CONFIG_USER_OK'"
   ```
   - Assert exit code 0 and output `CONFIG_USER_OK`.
3. **AI Skills Verification:**
   ```bash
   gitmap ssh exec u1 "ls -1 ~/.cursor/skills-cursor/ | wc -l"
   ```
   - Assert returned count >= 28.
4. **Project Storage Workspace Mappings:**
   ```bash
   gitmap ssh exec u1 "ls -d ~/.cursor/projects/home-a-git-work-* | wc -l"
   ```
   - Assert returned count matches migrated projects.
5. **Project Manager Workspace Index:**
   ```bash
   gitmap ssh exec u1 "grep -c 'rootPath' ~/.config/Cursor/User/globalStorage/alefragnani.project-manager/projects.json"
   ```
   - Assert count matches all 49 repositories located under `/home/a/git-work/*`.

### Step 4: Fleet Status Ledger Record Persistence
Update `repo-secrets/04-ubuntu-migration/cursor-fleet-status.json` with verification metadata:
```json
{
  "nodeAlias": "u1",
  "cursorInstalled": true,
  "cursorVersion": "3.23.12",
  "iconDeployed": true,
  "desktopEntryDeployed": true,
  "dockPinned": true,
  "startMenuGridPlaced": true,
  "dockPosition": "BOTTOM",
  "profileMigrated": true,
  "aiMemoryLoaded": true,
  "skillsCount": 28,
  "projectsCount": 49,
  "status": "HEALTHY",
  "lastVerified": "2026-10-05T..."
}
```

### Step 5: Linter Quality Gates Execution Matrix
Execute each linter sequentially and verify 0 violations:
1. **Relative Path Hygiene:**
   ```bash
   python linter-scripts/check-relative-paths.py
   python 03-ai-scripts/49-verify-privacy-and-relative-paths.py
   ```
   - Verifies zero absolute Windows/Linux drive letters or paths in specifications, code, and plans.
2. **Nested If Guard Clauses:**
   ```bash
   python linter-scripts/check-nested-ifs.py
   ```
   - Verifies shallow control flow and maximum nesting depth <= 2 across all modified Go code.
3. **Boolean Logic Guidelines:**
   ```bash
   python linter-scripts/check-boolean-guidelines.py
   ```
   - Verifies positive boolean naming prefixes (`is*`, `has*`, `can*`, `should*`).
4. **Go Code Formatting:**
   ```bash
   python scripts/format_go.py
   ```
   - Enforces standard `gofmt -w` formatting across `cli/`.
5. **Local CI/CD Runner:**
   ```bash
   python 03-ai-scripts/06-cicd-local-runner.py
   ```
   - Verifies all automated quality gates pass concurrently.

### Step 6: Minor Version Bump Ceremony (`v6.487.0`)
1. Invoke canonical version bump script:
   ```bash
   python 03-ai-scripts/37-bump-version.py --tier minor --scope "cursor - start menu app grid placement, os dock position configuration, and profile memory migration"
   ```
2. Verify updated files:
   - `version.json`: `Version` updated to `6.487.0`.
   - `package.json`: `version` synchronized to `6.487.0`.
   - `readme.md`: Header and version badges synchronized.
   - `changelog.md`: Appended `## [6.487.0]` section documenting:
     * Cursor GNOME Shell Start Menu (App Grid) placement in `app-picker-layout` on Ubuntu.
     * `gitmap os dock` position configuration command with `--node` remote delegation.
     * Complete Cursor profile, settings, extensions, and AI memory migration to Ubuntu fleet node `u1`.
     * `cursor-profile-transfer-notes.md` migration documentation in `repo-secrets`.
   - `cli/constants/constants.go`: Internal version string synchronized.

### Step 7: Atomic Commit & Pipeline Telemetry Monitoring
1. Staging and commit:
   - Commit title:
     `cursor - start menu app grid placement, os dock position configuration, and profile memory migration`
2. Remote pipeline tracking:
   ```bash
   gitmap pe -t
   ```
   - Monitor real-time progress until the pipeline displays `✔ 100% green` across all jobs.

---

## 3. Acceptance Criteria

- [ ] Remote verification confirms `cursor.desktop` present in `app-picker-layout` on `u1`.
- [ ] Remote verification confirms `update-desktop-database` refreshed application caches.
- [ ] Remote verification confirms `gitmap os dock` queries and updates dock position on `u1`.
- [ ] Remote verification confirms `settings.json`, Dracula theme, and keybindings loaded in `~/.config/Cursor/User/`.
- [ ] Remote verification confirms 28 AI skills and workspace project state present in `~/.cursor/`.
- [ ] Remote verification confirms 49 projects indexed in `projects.json` on `u1`.
- [ ] `cursor-fleet-status.json` updated with timestamped health status.
- [ ] All 5 quality gates pass with zero violations.
- [ ] Version successfully bumped to `v6.487.0` with synchronized manifests.
- [ ] Commit message strictly adheres to the hyphen-separated format.
- [ ] CI/CD telemetry monitored via `gitmap pe -t` passes 100% green.
