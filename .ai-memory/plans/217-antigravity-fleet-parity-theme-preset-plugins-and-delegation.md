# Master Plan: 217-antigravity-fleet-parity-theme-preset-plugins-and-delegation

## User Request (Verbatim)
> "Luke, so far the projects are there in the Antigravity. I appreciate that, but there is a great issue. The issue is when I go to the Antigravity, I do see you do not apply the same thing. Not this instance, but original, the default instances theme. And also the preset mode is not copied from the default instance. And plugins are not applied, skills are not added. So, so many things are not correct. So write the script, PowerShell additional script in the repo secrets folder properly, additional, and then after you do that, make modifies to the `gitmap` so that we can delegate in the future to `gitmap`, fully deploy the Antigravity IDE to another machine, and it will do it automatically. Do you understand? Is it clear?"

---

## 1. Executive Summary & Root Cause Synthesis

### User Symptoms & Observations
1. **Theme Mismatch**: Antigravity on Ubuntu `u1` is rendered in raw default unbranded styling rather than the default Windows instance's custom theme (`customThemeSeedsDark`: `#19191C` background, `#BD93F9` Dracula purple primary seed).
2. **Permission Preset Default Fallback**: In Antigravity Settings $\to$ Conversations (and per-project conversation prompt bar), Permission Preset displays as `Default` instead of the expected unattended execution mode (`CASCADE_COMMANDS_AUTO_EXECUTION_EAGER`, `BROWSER_JS_EXECUTION_POLICY_TURBO`, `ARTIFACT_REVIEW_MODE_TURBO`, and wide `globalPermissionGrants`).
3. **Plugins Missing**: `/home/a/.gemini/config/plugins/` is completely empty (0 plugins installed).
4. **Skills Missing**: All 43 plugin skills (`chrome-devtools`, `data-agent-kit`, `google-antigravity-sdk`, `modern-web-guidance`) are absent from the UI.
5. **Windows Path Leak on Linux**: An anomaly directory `/home/a/<windows-appdata>\Antigravity` was created on Linux because `data_dir` in `instances.json` serialized Windows backslashes.

### Forensic Root Causes
- **Antigravity Custom Theme & Settings Architecture**:
  - Primary UI theme seeds, execution policies, permission grants, and plugin registries live in `~/.gemini/config/config.json`.
  - When instances or remote workstations are provisioned, only bare skeleton directories were initialized; `config.json` and `plugins/` were never deployed.
- **Project Permission Overrides**:
  - `~/.gemini/config/projects/*.json` previously contained empty `settings: {}` and `permissionGrants: { allow: [] }`. To enable seamless preset parity, project descriptors must explicitly configure `fileAccessPolicy: AGENT_SETTING_POLICY_ALLOW` and `autoExecutionPolicy: CASCADE_COMMANDS_AUTO_EXECUTION_EAGER`.
- **GitMap Remote Delegation Gap**:
  - GitMap CLI lacked a dedicated command to bundle and deploy the full Antigravity environment (binaries, permissions preset, theme, plugins, skills, projects) over SSH in a single call.

---

## 2. Strict Workflow Constraints

- Total step budget: `N = 300` (Phase 1: Steps 1-150, Phase 2: Steps 151-300).
- Concurrency: `A = 2` autonomous subagents, `H = 2` hands per agent. Solo execution is strictly banned.
- Search engine: GitMap exclusive (`gitmap aum search`, `gitmap find`, `gitmap cat`, `gitmap ps`). Total ban on `Select-String`, `rg`, `grep`, `git grep`.
- Build/Test policy: Zero builds (`go build` ban) and zero full test suites (`go test` ban). Targeted fast syntax/linters only.
- Booleans: Positive prefixes only (`is`, `has`).
- Commit: Single atomic GitMap commit via `gitmap cpf "<module> - <summary>"` or `gitmap cpb "<module> - <summary>"` (hyphen separated, no colon in argument).

---

## 3. Disjoint Subtask Breakdown & Ownership Matrix

| Task ID | Subtask Code & Title | Owner | Target Files | Scope & Deliverable |
| :--- | :--- | :--- | :--- | :--- |
| **Task-01** | `01-antigravity-profile-and-preset-forensics` | Spec Writer 01 / Worker 01 | `02-spec/21-app/217-antigravity-fleet-parity-theme-preset-plugins-and-delegation/01-architecture-spec.md` | Formalize the complete protobuf/JSON schema for Antigravity themes (`customThemeSeedsDark`), permission presets (`CASCADE_COMMANDS_AUTO_EXECUTION_EAGER`), plugin manifests, and skills directories. |
| **Task-02** | `02-powershell-full-parity-sync-script` | Spec Writer 02 / Worker 02 | `repo-secrets/04-ubuntu-migration/sync-antigravity-full-profile.ps1` | Author standalone, idempotent PowerShell script in `repo-secrets` packaging default theme, turbo presets, 4 plugins, 43 skills, and deploying to `u1` via SSH with path sanitization. |
| **Task-03** | `03-gitmap-agy-deploy-delegation-engine` | Spec Writer 01 / Worker 01 | `cli/cmdagy/agy_deploy_cmd.go`, `cli/cmdagy/agy_deploy_types.go`, `cli/cmdssh/ssh_deploy_router.go`, `cli/cmd/roottooling.go`, `cli/cmd/help.go` | Implement `gitmap agy deploy <node>` (`--preset`, `--theme`, `--plugins`, `--skills`, `--all`, `--json`) and wire into `gitmap deploy ide <node>` and `gitmap migrate host`. |
| **Task-04** | `04-fleet-verification-and-scorecard` | Spec Writer 02 / Worker 02 | `02-spec/21-app/217-antigravity-fleet-parity-theme-preset-plugins-and-delegation/02-component-and-cli-spec.md`, `02-spec/22-app-issues/69-antigravity-fleet-parity-theme-preset-plugins-rca.md`, `02-spec/22-app-issues/readme.md` | Author component CLI spec, 4-part RCA document, and index in app-issues catalog. |
| **Task-05** | `05-live-execution-verification-and-push` | Lead Orchestrator | `repo-secrets`, `gitmap` codebase, master registries | Execute live sync on `u1`, verify Antigravity IDE UI state, run linters, update registries, and commit/push atomically via GitMap. |

---

## 4. Execution Plan by Waves

- **Wave 1 (Phase 1 Spec Generation)**:
  - Spec Writer 01: Authors Architecture Spec and Subtasks 01 & 03.
  - Spec Writer 02: Authors Component CLI Spec, RCA Document, and Subtasks 02 & 04.
- **Wave 2 (Phase 2 Code Execution)**:
  - Worker 01: Implements Task-03 (`gitmap agy deploy` delegation engine in Go).
  - Worker 02: Implements Task-02 (`sync-antigravity-full-profile.ps1` in `repo-secrets`) and updates issue catalog.
- **Wave 3 (Phase 3 Live Execution & Push)**:
  - Lead: Runs `sync-antigravity-full-profile.ps1` to provision `u1`, verifies state, updates registries, and executes atomic GitMap commit and push.
