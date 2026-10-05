# 227.02: Settings UI Modernization, Embedded Terminal Widget & Cross-Node Path Suggestion Engine

**Spec ID:** 227  
**Sub-Spec:** 02 (Component, UI & Path Suggestion Specification)  
**Status:** Approved  
**Version:** 1.0.0  
**Updated:** 2026-10-05  
**Subsystems:** `cli/cmdui`, `cli/completion`, `cli/cmdpipeline`, `install-quick.ps1`, `install-quick.sh`  
**Parent Spec:** [01-architecture-spec.md](01-architecture-spec.md)  
**Parent Plan:** [227-gitmap-cursor-fleet-delegate-settings-ui-and-cross-node-path-suggestions.md](../../../.ai-memory/plans/227-gitmap-cursor-fleet-delegate-settings-ui-and-cross-node-path-suggestions.md)  
**Design System Reference:** [02-spec/24-app-ui-design-system/01-index.md](../../24-app-ui-design-system/01-index.md)  

---

## 1. Executive Overview & Design Principles

This specification details the component, visual layout, and behavioral specifications for two core areas in Task 227:
1. **The Modernized GitMap Settings UI & Embedded Interactive Terminal Widget:**
   - Overhauls `#tab-settings` in `cli/cmdui/ui_assets.go` following the visual architecture and token hierarchy established in `02-spec/24-app-ui-design-system/`.
   - Organizes configuration into six ergonomic sub-tab categories: General, Cursor Fleet, Speed & Automation, SSH Nodes, Cache & Storage, and Terminal.
   - Fixes the two-way binding data-loss defect by coupling `loadSettings()` and `saveSettings()` to the `SettingsData.Attributes` key-value store, ensuring all 19 configuration parameters are reliably retrieved and persisted.
   - Integrates an in-browser interactive terminal widget with node delegation (`local` vs. remote SSH nodes), command history navigation, quick diagnostic chips, and execution status badges.
2. **The Cross-Node Path Suggestion Engine:**
   - Enhances Windows (`install-quick.ps1`) and POSIX (`install-quick.sh`) installation scripts with candidate directory discovery and numbered interactive selection menus.
   - Implements cross-node repository suggestions in `cli/completion/` and `cli/cmdpipeline/` by querying local SQLite indexes and cached remote fleet nodes (e.g., `u1:<repo>`).

---

## 2. Design Tokens, CSS Architecture & Component Primitives

The modernized Settings UI strictly consumes the design tokens defined in the GitMap App UI Design System. No hardcoded or unstyled raw colors are permitted.

### 2.1 Color Palette & Semantic Tokens

| Token Variable | Hex / Value | Semantic Role |
|:---|:---|:---|
| `--color-bg-canvas` | `#0f172a` (Slate 900) | Primary view background |
| `--color-bg-card` | `rgba(30, 41, 59, 0.7)` | Frosted card surface (`backdrop-filter: blur(12px)`) |
| `--color-bg-subtle` | `rgba(51, 65, 85, 0.5)` | Input field and inactive chip background |
| `--color-border-card` | `rgba(148, 163, 184, 0.15)` | Subtle 1px structural container boundary |
| `--color-border-focus` | `#06b6d4` (Cyan 500) | Active input focus ring and sub-tab underline |
| `--color-text-primary` | `#f8fafc` (Slate 50) | High-contrast headings and active labels |
| `--color-text-secondary` | `#94a3b8` (Slate 400) | Form helper text and placeholder hints |
| `--color-accent` | `#38bdf8` (Light Blue) | Primary button background and interactive highlights |
| `--color-accent-hover` | `#0284c7` (Sky 600) | Hover state for accent buttons |
| `--color-success` | `#10b981` (Emerald 500) | Status pill: online, exit code 0, saved confirmation |
| `--color-danger` | `#ef4444` (Red 500) | Status pill: error, failed command, exit code != 0 |
| `--color-warning` | `#f59e0b` (Amber 500) | Warning badges, threshold notices |
| `--color-console-bg` | `#0b1120` (Obsidian) | Terminal widget background |
| `--color-console-text` | `#e2e8f0` (Slate 200) | Terminal console stream font color |

### 2.2 Typography & Spacing Hierarchy
- **Primary Font Stack:** `-apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif`
- **Monospace Font Stack (Terminal & Code Inputs):** `"JetBrains Mono", "Fira Code", "Cascadia Code", Menlo, monospace`
- **Card Padding:** `20px` desktop / `16px` mobile
- **Form Row Gap:** `16px` vertical separation between input groups
- **Grid Column Gap:** `24px` horizontal column gutters for two-column desktop layouts

---

## 3. Settings UI Sub-Tab Layout & Form Schema

The monolithic card list in `#tab-settings` is structured into six focused sub-tabs. Operators switch sub-tabs without page reloading.

```
+---------------------------------------------------------------------------------------------------+
|  [ Settings View Header: GitMap Configuration & Node Management ]                                 |
+---------------------------------------------------------------------------------------------------+
|  [Sub-Tabs]                                                                                       |
|  (General)  |  (Cursor Fleet)  |  (Speed & Automation)  |  (SSH Nodes)  |  (Cache)  |  (Terminal) |
+---------------------------------------------------------------------------------------------------+
|  Active Sub-Tab Card Container                                                                    |
|  ... [Sub-Tab Content & Field Rows] ...                                                           |
+---------------------------------------------------------------------------------------------------+
|  Footer Actions:  [Save All Settings Button]  [Status Notification Toast / Badge]                 |
+---------------------------------------------------------------------------------------------------+
```

### 3.1 Sub-Tab 1: General Settings
Controls global runtime behavior and user interface presentation.

| Field Label | Input ID | Bound Key | Type / Options | Default |
|:---|:---|:---|:---|:---|
| Visual Theme | `setting-theme` | `theme` | Select (`dark`, `light`, `oled`) | `dark` |
| Default Remote | `setting-remote` | `defaultRemote` | Text input | `origin` |
| Cluster REST Port | `setting-port` | `clusterPort` | Number input (1024-65535) | `49152` |
| Graphics Mode | `setting-graphics` | `graphicsMode` | Select (`high`, `low`, `off`) | `high` |
| Auto-Open Browser | `setting-auto-open` | `autoOpenBrowser` | Select (`true`, `false`) | `true` |
| Commit Layout | `setting-commitin-layout` | `commitInLayout` | Select (`split`, `unified`) | `split` |
| Pull Direction | `setting-pull-direction` | `pullDirection` | Select (`pull-left`, `pull-right`) | `pull-left` |
| PR Replay Mode | `setting-pr-mode` | `prReplayMode` | Select (`merges`, `fast-forward`) | `merges` |

### 3.2 Sub-Tab 2: Cursor Fleet Settings
Configures fleet-wide Cursor IDE migration and workspace sync options.

| Field Label | Input ID | Bound Key | Type / Options | Default |
|:---|:---|:---|:---|:---|
| Default Delegation Node | `setting-cursor-node` | `attributes["cursor.default_node"]` | Text input (e.g. `u1`) | `u1` |
| Target Workdir Mapping | `setting-cursor-workdir` | `attributes["cursor.target_workdir"]` | Text input | `/home/a/git-work` |
| Default Exclusions | `setting-cursor-exclude` | `attributes["cursor.exclude_patterns"]` | Text input | `projects/temp*,logs` |
| Auto-Backup Target State | `setting-cursor-backup` | `attributes["cursor.auto_backup"]` | Select (`true`, `false`) | `true` |
| Force Overwrite Existing | `setting-cursor-force` | `attributes["cursor.force_overwrite"]` | Select (`true`, `false`) | `false` |

### 3.3 Sub-Tab 3: Speed & Automation Settings
Preserves operational latency, threshold, and alerting webhooks (previously dropped by the save bug).

| Field Label | Input ID | Bound Key | Type / Options | Default |
|:---|:---|:---|:---|:---|
| LAP Default Lookback Hours | `setting-lap-hours` | `attributes["lap.default_hours"]` | Number input (1-168) | `24` |
| Account Switch FF Threshold (%) | `setting-account-switch-threshold` | `attributes["account_switch.threshold"]` | Number input (1-100) | `15` |
| Telegram Bot Token | `setting-telegram-token` | `attributes["telegram.bot_token"]` | Masked text input | `""` |
| Telegram Target Chat ID | `setting-telegram-chat` | `attributes["telegram.chat_id"]` | Text input | `""` |
| Email SMTP Host | `setting-email-host` | `attributes["email.smtp_host"]` | Text input | `""` |
| Email SMTP Port | `setting-email-port` | `attributes["email.smtp_port"]` | Number input | `587` |
| Email Sender Address | `setting-email-sender` | `attributes["email.sender"]` | Text input | `""` |
| Email Recipient Address | `setting-email-recipient` | `attributes["email.recipient"]` | Text input | `""` |

### 3.4 Sub-Tab 4: SSH Nodes & Machine Identity
Manages SSH credentials and machine identifiers across nodes.

| Field Label | Input ID | Bound Key | Type / Options | Default |
|:---|:---|:---|:---|:---|
| Auto-Deploy SSH Keys | `setting-auto-deploy-keys` | `isAutoDeployKey` | Select (`true`, `false`) | `true` |
| Local Machine Alias | `setting-machine-alias` | `attributes["machine.alias"]` | Text input | `""` |
| Datacenter / Zone Tag | `setting-machine-dc` | `attributes["machine.datacenter"]` | Text input | `local` |
| Node Discovery Interval (sec) | `setting-discovery-interval` | `attributes["cluster.discovery_sec"]` | Number input | `30` |

### 3.5 Sub-Tab 5: Cache & Storage Maintenance
Maintains Split-DB files and system cache cleaning.

| Feature / Action | ID / Endpoint | Action Behavior |
|:---|:---|:---|
| Telemetry Retention (Days) | `setting-retention-days` / `attributes["db.retention_days"]` | Number input (default `30`) |
| Vacuum Split-DB Databases | `btn-vacuum-splitdb` / `POST /api/db/vacuum` | Triggers SQLite `VACUUM` across `gitmap.db`, `installation.db`, `commands.db` |
| Purge Temporary Bundles | `btn-purge-cache` / `POST /api/cache/purge` | Deletes staging archives in `temp/` |

---

## 4. Two-Way Binding Engine Implementation

### 4.1 Client-Side Hydration (`loadSettings()`)
When the Settings view initializes:
1. `GET /api/settings` retrieves the envelope payload containing standard properties and the `attributes` dictionary.
2. Direct properties (`theme`, `defaultRemote`, `clusterPort`, `graphicsMode`, `autoOpenBrowser`, `commitInLayout`, `pullDirection`, `prReplayMode`, `isAutoDeployKey`) populate corresponding DOM elements.
3. The engine iterates over `data.attributes`:
   ```javascript
   if (data.attributes) {
     for (const [key, val] of Object.entries(data.attributes)) {
       const el = document.querySelector(`[data-key="${key}"]`);
       if (el) { el.value = val; }
     }
   }
   ```

### 4.2 Client-Side Serialization (`saveSettings()`)
When the operator clicks "Save All Settings":
1. The engine constructs the base `SettingsData` object.
2. The engine scans all inputs with attribute `data-key` and packs them into `attributes`:
   ```javascript
   const attributes = {};
   document.querySelectorAll('[data-key]').forEach(el => {
     const k = el.getAttribute('data-key');
     if (k) { attributes[k] = el.value.trim(); }
   });
   payload.attributes = attributes;
   ```
3. `POST /api/settings` commits the payload to `~/.gitmap/ui_settings.json` wrapped in the canonical `jsonenvelope.TypeUISettings` envelope.
4. An unobtrusive toast notification indicates `✔ Settings saved successfully`.

---

## 5. Sub-Tab 6: Embedded Interactive Terminal Widget

### 5.1 Terminal Component Wireframe & DOM Structure

```
+---------------------------------------------------------------------------------------------------+
|  Terminal Header                                                                                  |
|  [>_ Console]   Node: [ Local Host (Host OS)  v ]    [ Status: Ready ]    [ Clear ]  [ Copy Output ]|
+---------------------------------------------------------------------------------------------------+
|  Quick Command Chips:                                                                             |
|  [ gitmap status ]  [ gitmap nodes list ]  [ gitmap cur delegate --node u1 --dry-run ]  [ gitmap db]|
+---------------------------------------------------------------------------------------------------+
|  Terminal Console Viewport (Height: 360px, Overflow-Y: Auto, JetBrains Mono 13px)                 |
|  gitmap:~$ gitmap nodes list                                                                      |
|  ----------------------------------------------------------------                                 |
|  #1  alias: u1   host: 192.168.1.105   os: linux (ubuntu)   online: true                          |
|  ----------------------------------------------------------------                                 |
|  [Exit Code: 0 | Execution Time: 42ms]                                                            |
|  gitmap:~$ _                                                                                      |
+---------------------------------------------------------------------------------------------------+
|  Command Input Row                                                                                |
|  [ u1:~$ ]  [ Enter command...                                         ]  [ Run Button (Enter) ]  |
+---------------------------------------------------------------------------------------------------+
```

### 5.2 Interactive Behavioral Specifications
- **Node Selection:** Selecting a remote node (e.g. `u1`) automatically changes the input prompt prefix to `u1:~$` and sets `nodeAlias: "u1"` in the request payload.
- **Quick Action Chips:** Clicking any chip automatically fills the input field, focuses it, and optionally triggers execution if double-clicked.
- **Command History:**
  - Up arrow (`ArrowUp`) cycles backwards through the execution history stack.
  - Down arrow (`ArrowDown`) cycles forward.
  - History is persisted in browser `localStorage` up to 50 items.
- **Console Scroll & Styling:**
  - Standard output is displayed in clean slate white (`#e2e8f0`).
  - Standard error output is rendered in soft crimson (`#f87171`).
  - Upon command completion, an inline metadata badge displays execution duration and return code:
    `[Exit Code: 0 | 48ms]` (green) or `[Exit Code: 1 | 215ms]` (red).
  - Auto-scrolls to the bottom upon appending new output.

### 5.3 Backend Handler (`POST /api/terminal/exec`)
The server executes commands within strict security boundaries:
1. **Context Deadline Guard:** Bounded by `time.Duration(timeoutSec) * time.Second` (clamped between 1 and 120 seconds, default 30 seconds). If expired, the underlying process group is killed and exit code `124` is returned with an explanatory timeout error.
2. **Execution Dispatch:**
   - When `nodeAlias == "" || nodeAlias == "local"`:
     - On Windows: `powershell.exe -NoProfile -NonInteractive -Command "<command>"`
     - On Linux/macOS: `sh -c "<command>"`
   - When `nodeAlias` is a remote alias: routes command through `cmdssh.ExecuteRemoteCommand(nodeAlias, command, timeoutSec)`.
3. **Response Enveloping:** Streams combined output, exit code, and duration back in `TerminalExecResp` JSON.

---

## 6. Cross-Node Path Suggestion Engine

### 6.1 Installer Candidate Path Probing

#### Windows (`install-quick.ps1`)
When run interactively (or without explicit `-InstallDir`), the script executes proactive candidate probing before requesting user input:
1. Probing Steps:
   - Probes `D:\gitmap`, `C:\gitmap`, `E:\gitmap`, and `$HOME\gitmap`.
   - Inspects existing `powershell.json` deployed in candidate paths.
   - Checks `$env:PATH` for an active `gitmap.exe`.
   - Inspects drives with > 5 GB free disk space.
2. Menu Presentation:
   ```
     gitmap quick installer — installation directory selection
     ----------------------------------------------------------
     [1] D:\gitmap (Recommended - Existing Installation / High-Speed SSD)
     [2] C:\gitmap (System Drive - 42 GB Free)
     [3] E:\Tools\gitmap (Alternate Volume - 120 GB Free)
     [4] Custom path...

     Select an option [1-4] or enter custom path (Default: [1]):
   ```
3. Input Handling:
   - If the operator presses `Enter`, candidate `[1]` is automatically selected.
   - If the operator enters a single digit `1-3`, the matching candidate path is selected.
   - If the operator enters an absolute filesystem path, that path is used.
   - In non-interactive mode (`-Interactive` omitted), option `[1]` is selected silently.

#### Linux / macOS (`install-quick.sh`)
When executed in interactive mode:
1. Probing Steps:
   - If `id -u == 0` (root): candidates are `/usr/local/bin` (Recommended) and `/opt/gitmap`.
   - If non-root: candidates are `${HOME}/.local/bin` (Recommended) and `${HOME}/gitmap`.
   - Checks `which gitmap` to discover existing installations.
2. Menu Presentation:
   ```
     gitmap quick installer — installation directory selection
     ----------------------------------------------------------
     [1] ~/.local/bin (Recommended - Standard User Binary Directory)
     [2] ~/gitmap (Standalone GitMap Suite Directory)
     [3] Custom path...

     Select an option [1-3] or enter custom path (Default: [1]):
   ```
3. Fallback:
   - Non-interactive piping (`curl ... | bash` or `eval`) defaults to candidate `[1]` immediately without stalling `stdin`.

---

### 6.2 CLI Cross-Node Completion & Suggestion Aggregator

#### Data Model (`cli/completion/cross_node_suggest.go`)
```go
type CrossNodeCandidate struct {
    NodeAlias  string // "local" or SSH node alias (e.g. "u1")
    RepoSlug   string // e.g. "alimtvnetwork/gitmap-v28" or "gitmap-v28"
    FullPath   string // local path or POSIX remote path (/home/a/git-work/...)
    IsRemote   bool
}
```

#### Query & Union Logic
1. Local repositories are queried from `store.DB` (`Repo` and `WorkDir` tables).
2. Remote fleet node repositories are discovered from `cmdssh.FetchAllSSHConnections()` and cached inventory records.
3. For remote nodes, items are formatted with the node prefix:
   - `<nodeAlias>:<repoSlug>` (e.g., `u1:gitmap-v28`, `u1:movie-cli-v8`)
   - `<nodeAlias>:<remotePath>` (e.g., `u1:/home/a/git-work/gitmap-v28`)
4. In `cli/completion/dynamic.go`:
   - `Dynamic()` integrates `DynamicRepoSupplier` and includes cross-node candidates when completing `cd`, `clone`, `cfr`, `cur`, `cursor`, `nodes`, and `ssh`.

#### Diagnostic Formatting in `cli/cmdpipeline/pipeline_repo_suggest.go`
When a requested target repo is not found locally, `PrintRepoTargetNotFoundDiagnostic` queries both local and remote inventories:
```
  ✖ Repository not found: "gitmap-v28"

  Did you mean (Local):
    • gitmap pe gitmap-v28
    • gitmap cd gitmap-v28

  Did you mean (Remote Fleet Node u1):
    • gitmap nodes run u1 -- gitmap pe gitmap-v28
    • gitmap cur delegate --node u1
```

---

## 7. Acceptance Criteria & Verification Matrix

| ID | Specification Requirement | Verification Procedure |
|:---|:---|:---|
| **AC-227-08** | Sub-tab navigation in `#tab-settings` switches active views without DOM refetch | DOM click event test across all 6 tabs |
| **AC-227-09** | All 19 configuration fields round-trip through `loadSettings()` and `saveSettings()` without data loss | Browser test asserting equality before/after save |
| **AC-227-10** | Embedded terminal executes local commands via `POST /api/terminal/exec` and returns stdout and exit code | HTTP POST test with command `gitmap --version` |
| **AC-227-11** | Embedded terminal times out and terminates long-running commands exceeding `timeoutSec` | HTTP POST test with sleep command exceeding limit |
| **AC-227-12** | `install-quick.ps1` probes candidate paths and supports single-digit selection | PowerShell Pester / CLI simulated interactive input |
| **AC-227-13** | `install-quick.sh` provides candidate menu and handles non-interactive execution safely | Bash smoke test with piped input |
| **AC-227-14** | CLI completion suggests `u1:<repo>` when completing cross-node paths | Unit test in `cli/completion/cross_node_suggest_test.go` |

---

## 8. Cross-References

- [01-architecture-spec.md](01-architecture-spec.md) — Architectural Overview & Flag Specifications
- [Subtask 02: Settings UI](../../../.ai-memory/plans/subtasks/227-gitmap-cursor-fleet-delegate-settings-ui-and-cross-node-path-suggestions/02-settings-ui-design-and-embedded-terminal.md)
- [Subtask 03: Path Suggestions](../../../.ai-memory/plans/subtasks/227-gitmap-cursor-fleet-delegate-settings-ui-and-cross-node-path-suggestions/03-cross-node-path-suggestions-engine.md)
- [Subtask 04: Verification & Release](../../../.ai-memory/plans/subtasks/227-gitmap-cursor-fleet-delegate-settings-ui-and-cross-node-path-suggestions/04-testing-linters-and-release.md)
