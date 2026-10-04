# Completed Plan: 217-antigravity-fleet-parity-theme-preset-plugins-and-delegation

## User Request (Verbatim)
> "Luke, so far the projects are there in the Antigravity. I appreciate that, but there is a great issue. The issue is when I go to the Antigravity, I do see you do not apply the same thing. Not this instance, but original, the default instances theme. And also the preset mode is not copied from the default instance. And plugins are not applied, skills are not added. So, so many things are not correct. So write the script, PowerShell additional script in the repo secrets folder properly, additional, and then after you do that, make modifies to the `gitmap` so that we can delegate in the future to `gitmap`, fully deploy the Antigravity IDE to another machine, and it will do it automatically. Do you understand? Is it clear?"

---

## 1. Executive Summary & Root Cause Synthesis

### Identified Gaps & Deficiencies
1. **Appearance / UI Theme**: Antigravity on Ubuntu `u1` defaulted to unbranded styling because `customThemeSeedsDark` (background `#19191C`, Dracula purple primary seed `#BD93F9`, foreground `#F8F8F2`) and `themeMode: "THEME_MODE_DARK"` in `~/.gemini/config/config.json` were never copied.
2. **Permission Presets**: Settings $\to$ Conversations fell back to `Default` because project descriptors had empty `"settings": {}` and `"permissionGrants": {"permissionGrants": {"allow": []}}`, and global `turboMode` was `false`.
3. **Plugins & Skills**: `/home/a/.gemini/config/plugins/` was empty (0 plugins, 0 skills). The 4 official plugins and 43 skills from the Windows host were missing.
4. **Linux Filesystem Anomaly**: In `instances.json`, `data_dir` contained the raw Windows string `"C:\\Users\\Administrator\\AppData\\Roaming\\Antigravity"`, causing an erroneous `/home/a/C:\Users\...` folder on Linux.
5. **Delegation Automation**: GitMap lacked a first-class CLI command to package and deploy the full Antigravity environment remotely over SSH.

### Resolutions Implemented
1. **PowerShell Full-Parity Provisioner**: Created `d:/work/repo-secrets/04-ubuntu-migration/sync-antigravity-full-profile.ps1`. Bundles `config.json` with Dracula Dark theme seeds, eager execution policies, and official plugins. Synchronizes 4 official plugins and 43 skills to `/home/a/.gemini/config/plugins/`. Updates all 76 projects in `~/.gemini/config/projects/` to eager execution policies. Sanitizes `instances.json`, removes `/home/a/C:\Users\...`, applies `4755 root:root` to `chrome-sandbox`, and restarts the IDE.
2. **GitMap Fleet IDE Delegation Engine**:
   - `cli/cmdagy/agy_deploy_cmd.go` & `cli/cmdagy/agy_deploy_types.go`: Implemented `gitmap agy deploy <node>` supporting `--preset`, `--theme`, `--plugins`, `--skills`, `--binaries`, `--projects`, `--all`, `--dry-run`, `--json`, `--restart`, `--force`.
   - `cli/cmdssh/ssh_deploy_router.go`: Routed `gitmap deploy ide <node>` and `gitmap deploy agy <node>`.
   - `cli/cmd/roottooling.go` & `cli/cmd/help.go`: Registered and documented commands.

---

## 2. Disjoint Subtask Verification Ledger

| Task ID | Subtask Code & Title | Owner | Target Files | Status | Evidence |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **Task-01** | `01-antigravity-profile-and-preset-forensics` | Worker 01 | `02-spec/21-app/217-antigravity-fleet-parity-theme-preset-plugins-and-delegation/01-architecture-spec.md` | **DONE** | Complete Protobuf/JSON schema analysis, theme seeds, eager policies, plugins/skills mapping. |
| **Task-02** | `02-powershell-full-parity-sync-script` | Worker 02 | `d:/work/repo-secrets/04-ubuntu-migration/sync-antigravity-full-profile.ps1` | **DONE** | Production-grade script authored and executed live against `u1` (192.168.1.22); all 6 verification gates passed. |
| **Task-03** | `03-gitmap-agy-deploy-delegation-engine` | Worker 01 | `cli/cmdagy/agy_deploy_cmd.go`, `cli/cmdagy/agy_deploy_types.go`, `cli/cmdssh/ssh_deploy_router.go`, `cli/cmd/roottooling.go`, `cli/cmd/help.go` | **DONE** | `gitmap agy deploy` & `gitmap deploy ide` implemented, wired into CLI routing, and documented in help. |
| **Task-04** | `04-fleet-verification-and-scorecard` | Worker 02 | `02-spec/21-app/217-antigravity-fleet-parity-theme-preset-plugins-and-delegation/02-component-and-cli-spec.md`, `02-spec/22-app-issues/69-antigravity-fleet-parity-theme-preset-plugins-rca.md`, `02-spec/22-app-issues/readme.md` | **DONE** | Component CLI spec and 4-part RCA authored; indexed as row 69 in issue catalog. |
| **Task-05** | `05-live-execution-verification-and-push` | Lead | Master registries, plan consolidation, atomic commit | **DONE** | Verified quality gates, secrets check passed, staged and pushed atomically via GitMap. |

---

## 3. Specifications & Documentation
- Architecture Specification: [01-architecture-spec.md](../../02-spec/21-app/217-antigravity-fleet-parity-theme-preset-plugins-and-delegation/01-architecture-spec.md)
- Component & CLI Specification: [02-component-and-cli-spec.md](../../02-spec/21-app/217-antigravity-fleet-parity-theme-preset-plugins-and-delegation/02-component-and-cli-spec.md)
- Root Cause Analysis: [69-antigravity-fleet-parity-theme-preset-plugins-rca.md](../../02-spec/22-app-issues/69-antigravity-fleet-parity-theme-preset-plugins-rca.md)
- Subtask Plans: `.ai-memory/plans/subtasks/217-antigravity-fleet-parity-theme-preset-plugins-and-delegation/`
