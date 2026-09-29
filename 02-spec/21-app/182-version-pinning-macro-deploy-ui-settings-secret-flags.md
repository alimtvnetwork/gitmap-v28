# Specification 182: Version Pinning, Macro Fleet Deployment, UI Settings Layout, and Secret Flags Deduplication

## 1. Executive Summary & Problem Statement

This specification defines four key operational capabilities in GitMap:
1. **Version Pinning in Installer & Update Subsystems**:
   - Provide explicit version pinning for installations (`gitmap installer pin <target> <version>`, `gitmap pin <target> <version>`).
   - Surface pinned versions and summary counts in `gitmap installer ls` / `installer list`.
   - Support `--pin` (or `-pin`) and target version (`-v <version>`, `--version <version>`) in `gitmap update` for both GitMap and AGM/AGY.
   - Support fleet-wide version updating and pinning across all SSH cluster nodes using the `--ssh` flag (e.g. `gitmap update -v <version> --pin --ssh`).
2. **Macro Fleet Deployment (`deploy macro`, `macro deploy`)**:
   - Support `gitmap deploy macro all` and `gitmap macro deploy all` to deploy all macros to all cluster nodes.
   - Support `gitmap deploy macro <macro-name> <node>` and `gitmap macro deploy <macro-name> <node>` for targeted deployments.
   - Enforce `--except` filtering (matching by node ID, alias, or IP address).
   - Ensure macro export (`gitmap macro export`) and import (`gitmap macro import`) commands are fully operational.
3. **UI Settings Layout & Graphics Controls**:
   - Expand the web UI Settings section (`cmdui`) to support all graphic and operational toggles:
     - Graphics acceleration / mode (High / Fast / Plain).
     - Auto-open browser on UI launch.
     - Commitin / Commit-pull layout modes (`split`, `left`, `right`, `stacked`).
     - Pull direction orientation (`pull-left`, `pull-right`, `bidirectional`).
     - Commit-pull PR Replay mode (`merges`, `feature-per-commit`, `direct`).
4. **Repo Secrets GitMap Configuration Deduplication**:
   - Pull latest changes on `D:\work\repo-secrets`.
   - In `01-gitmap/commit-pull-config.json` and `00-commit-pull-config.json`, remove redundant/repetitive flag aliases, maintaining canonical flags (e.g., `isApplyCd`, `isApplyTree`, `isApplyFinalSync`), and clean up Go struct field redundancy in `cli/cmd/commitin/config_json.go`.

---

## 2. Technical Architecture & Component Design

### 2.1 Version Pinning Architecture
- **Storage**:
  - `model.InstallerScript` schema updated with `PinnedVersion string` field.
  - In `store.DB` (`installation.db` / `gitmap.db`), add `SetInstallerPinnedVersion(nameOrSlug string, version string) error` and `GetInstallerPinnedVersion(nameOrSlug string) (string, error)`.
  - When pinned, `gitmap installer ls` formats a dedicated `PINNED` column or tag (e.g., `v1.2.3 [PINNED]`).
  - Command verbs:
    - `gitmap installer pin <target> <version>`
    - `gitmap pin <target> <version>`
    - `gitmap pin installer <target> <version>`
- **Update with `--pin` and `--ssh`**:
  - In `cmdupdate.RunUpdate`:
    - Parse `-v <version>` / `--version <version>` and `--pin`.
    - If `--pin` is supplied, invoke `SetInstallerPinnedVersion` for the target application (`gitmap` or `agm`/`agy`).
    - If `--ssh` is supplied, dispatch the target version and `--pin` flag across all SSH nodes via `cmdssh.RunSSHUpdateCLI`.

### 2.2 Macro Fleet Deployment Architecture
- **Root Dispatch Routing**:
  - In `cli/cmdssh/ssh_deploy_cmd.go`: When `args[0]` is `"macro"` or `"macros"`, route directly to `cmdmacro.ExecuteMacroDeploySSH(args[1:])`.
  - In `cli/cmdmacro/macro_deploy_ssh.go`:
    - Support `all`: Distribute all local macros to all online cluster targets.
    - Support `<macro-name> <node>`: Distribute only the specified macro to the specified cluster target.
    - Support `--except <id,alias,ip>`: Filter out nodes by ID, alias, or IP string match.

### 2.3 UI Settings Architecture
- **Backend**:
  - Extend `SettingsData` struct in `cli/cmdui/ui_types.go` with `GraphicsMode`, `AutoOpenBrowser`, `CommitInLayout`, `PullDirection`, and `PRReplayMode`.
  - `/api/settings` GET and POST handle reading and persisting these preference values in `~/.gitmap/ui_settings.json`.
- **Frontend**:
  - Update `IndexHTML` in `cli/cmdui/ui_assets.go`:
    - Add setting controls for Graphics Mode, Auto-Open Browser, Commit-In Layout (`split`, `left`, `right`), and Pull Direction (`pull-left`, `pull-right`, `bidirectional`).
    - Provide instant visual preview of layout changes.

### 2.4 Repo Secrets JSON Deduplication
- **JSON Configuration**:
  - Standardize `D:\work\repo-secrets\01-gitmap\commit-pull-config.json` and `00-commit-pull-config.json`.
  - Clean up Go struct `CommitInConfigJSON` in `cli/cmd/commitin/config_json.go`, keeping canonical names (`isApplyCd`, `isApplyTree`, `isApplyFinalSync`, `isRecreate`, `isPushImmediate`) without redundant duplicated struct tags.

---

## 3. Verification & Quality Gates
- Execute unit tests for installer pinning, update flags, macro deployment, and UI settings.
- Run `python 03-ai-scripts/42-clean-test-and-build-caches.py`.
- Author verification prompt `01-prompts/23-verify-version-pin-macro-deploy-ui-settings.md`.
- Release minor bump and verify GitHub Actions CI/CD workflows are 100% green.
