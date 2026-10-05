# Subtask Plan 04: Remote Node u1 Live Verification, Linter Quality Gates & Release Ceremony

- **Spec Reference:** [02-spec/21-app/224-ubuntu-cursor-dock-pin-projects-os-update-and-remote-uninstall/02-component-and-cli-spec.md](../../../../02-spec/21-app/224-ubuntu-cursor-dock-pin-projects-os-update-and-remote-uninstall/02-component-and-cli-spec.md)
- **Status:** Queued
- **Target Area:** `repo-secrets/04-ubuntu-migration/cursor-fleet-status.json`, `linter-scripts/`, `03-ai-scripts/`, `version.json`, `package.json`, `changelog.md`, `cli/cmdpipeline/`

---

## 1. Objective

Orchestrate and execute the comprehensive verification, quality gate enforcement, and release ceremony for task 224:
1. Conduct live remote verification across Ubuntu node `u1` via SSH to confirm `gitmap os update` auto-elevation, Cursor desktop launcher/dock favorites, Project Manager repository indexing, and dry-run dev tools uninstallation.
2. Validate zero regressions across repository linters: relative path hygiene, nested-if limits, boolean standards, and Go code formatting.
3. Perform the official minor version bump from `v6.484.0` to `v6.485.0` via `03-ai-scripts/37-bump-version.py`.
4. Stage changes using GitMap's hyphen-separated atomic commit convention:
   `cursor - ubuntu dock pin, projects sync, os update sudo elevation, and dev tools remote uninstall`
5. Monitor CI/CD pipeline health in real time via `gitmap pe -t` until 100% green.

---

## 2. Implementation Details

### Step 1: Remote Node `u1` Live Verification Workflow
Execute remote commands over SSH via `gitmap ssh exec u1 "<command>"`:
1. **OS Update Auto-Elevation:**
   ```bash
   gitmap ssh exec u1 "gitmap os update"
   ```
   - Assert stdout contains: `apt : ✔ OK` (confirming `sudo -n` bypassed exit status 100).
2. **Cursor Desktop Integration:**
   ```bash
   gitmap ssh exec u1 "test -f /usr/share/applications/cursor.desktop && test -f /usr/share/pixmaps/cursor.png && echo 'INTEGRATION_OK'"
   ```
   - Assert exit code 0 and output `INTEGRATION_OK`.
3. **GNOME Dock Favorite-Apps Pinning:**
   ```bash
   gitmap ssh exec u1 "DBUS_SESSION_BUS_ADDRESS=unix:path=/run/user/1000/bus gsettings get org.gnome.shell favorite-apps"
   ```
   - Assert the returned list string contains `'cursor.desktop'`.
4. **Workspace Projects Synchronization:**
   ```bash
   gitmap ssh exec u1 "grep -c 'rootPath' ~/.config/Cursor/User/globalStorage/alefragnani.project-manager/projects.json"
   ```
   - Assert project count matches all 49 repositories located under `/home/a/git-work/*`.

### Step 2: Dev Tools Uninstall Dry-Run Verification on `u1`
Execute dry-run uninstallation commands on node `u1` across all 5 target tools without modifying system packages or deleting files:
1. `gitmap uninstall pwsh --node u1 --dry-run`
2. `gitmap uninstall agm --node u1 --dry-run`
3. `gitmap uninstall vim --node u1 --dry-run`
4. `gitmap uninstall vscode --node u1 --dry-run`
5. `gitmap uninstall cursor --node u1 --dry-run`

Verify each command:
- Reports `[DryRun]` action items.
- Confirms zero files deleted or packages altered.
- Exits with status code 0.

### Step 3: Fleet Status Ledger Record Persistence
Update `repo-secrets/04-ubuntu-migration/cursor-fleet-status.json` with verification metadata:
```json
{
  "nodeAlias": "u1",
  "cursorInstalled": true,
  "dockPinned": true,
  "projectsSynced": 49,
  "osUpdateSudoElevation": true,
  "uninstallDryRunVerified": true,
  "status": "HEALTHY",
  "lastVerified": "2026-10-05T..."
}
```

### Step 4: Linter Quality Gates Execution Matrix
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
   - Verifies all 21 automated quality gates pass concurrently.

### Step 5: Minor Version Bump Ceremony (`v6.485.0`)
1. Invoke canonical version bump script:
   ```bash
   python 03-ai-scripts/37-bump-version.py --tier minor --scope "cursor - ubuntu dock pin, projects sync, os update sudo elevation, and dev tools remote uninstall"
   ```
2. Verify updated files:
   - `version.json`: `Version` updated to `6.485.0`.
   - `package.json`: `version` synchronized to `6.485.0`.
   - `readme.md`: Header and version badges synchronized.
   - `changelog.md`: Appended `## [6.485.0]` section documenting:
     * Cursor Ubuntu dock pinning, official icons, and project manager sync.
     * `gitmap os update` sudo elevation and failure output reporting.
     * Dev tools remote and local uninstall subsystem (`pwsh`, `agm`, `vim`, `vscode`, `cursor`) with `--node` and `--dry-run`.
   - `cli/constants/constants.go`: Internal version string synchronized.

### Step 6: Atomic Commit & Pipeline Telemetry Monitoring
1. Staging and commit (executed during phase 5):
   - Commit title:
     `cursor - ubuntu dock pin, projects sync, os update sudo elevation, and dev tools remote uninstall`
2. Remote pipeline tracking:
   ```bash
   gitmap pe -t
   ```
   - Monitor real-time progress until the pipeline displays `✔ 100% green` across all jobs.

---

## 3. Acceptance Criteria

- [ ] Remote verification confirms `apt : ✔ OK` for `gitmap os update` on `u1`.
- [ ] Remote verification confirms `cursor.desktop` in `org.gnome.shell favorite-apps` on `u1`.
- [ ] Remote verification confirms 49 projects indexed in `projects.json` on `u1`.
- [ ] Remote dry-run uninstallation succeeds across all 5 target tools on `u1` without data loss.
- [ ] `cursor-fleet-status.json` updated with timestamped health status.
- [ ] All 5 quality gates pass with zero violations.
- [ ] Version successfully bumped to `v6.485.0` with synchronized manifests.
- [ ] Commit message strictly adheres to the hyphen-separated format.
- [ ] CI/CD telemetry monitored via `gitmap pe -t` passes 100% green.
