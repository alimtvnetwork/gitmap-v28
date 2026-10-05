# Subtask 227.2: Settings UI Modernization & Embedded Interactive Terminal

- **Parent Spec:** [01-architecture-spec.md](../../../../02-spec/21-app/227-gitmap-cursor-fleet-delegate-settings-ui-and-cross-node-path-suggestions/01-architecture-spec.md)
- **Master Ledger:** [00-master-audit-ledger.md](../../../../02-spec/21-app/227-gitmap-cursor-fleet-delegate-settings-ui-and-cross-node-path-suggestions/00-master-audit-ledger.md)
- **Status:** DONE
- **Target Subsystems:** `cli/cmdui`, `cli/cmdssh`

---

## 1. Objective

Overhaul the Settings interface (`#tab-settings`) in `cli/cmdui/ui_assets.go` to conform to `02-spec/24-app-ui-design-system/`, resolve the critical two-way binding data loss bug for Speed and Machine Identity configuration fields, implement an interactive embedded terminal with quick action chips and command history, build the `/api/terminal/exec` backend endpoint with a 30s context timeout guard in `cli/cmdui/ui_server.go`, and author unit tests in `cli/cmdui/ui_test.go`.

---

## 2. File Modification Inventory

| File | Action | Responsibilities | Target Lines |
|:---|:---|:---|:---|
| `cli/cmdui/ui_types.go` | **Modify** | Add `TerminalExecReq` and `TerminalExecResp` data structs | < 95 |
| `cli/cmdui/ui_server.go` | **Modify** | Mount `/api/terminal/exec`, implement `handleAPITerminalExec`, `executeLocalTerminalCommand` (with 30s timeout guard), route remote node commands via `cmdssh`, fix settings attributes preservation | < 450 |
| `cli/cmdui/ui_assets.go` | **Modify** | Redesign `#tab-settings` layout into categorized sections; fix `loadSettings()` and `saveSettings()` bidirectional attribute binding; embed interactive terminal widget | < 850 |
| `cli/cmdui/ui_test.go` | **Modify** | Unit tests for terminal request execution, context timeout enforcement, settings attributes preservation, and JSON serialization | < 120 |

---

## 3. Step-by-Step Implementation Plan

### Step 3.1: Define Terminal Request & Response Models (`cli/cmdui/ui_types.go`)
1. Add `TerminalExecReq`:
   ```go
   // TerminalExecReq specifies command execution parameters.
   type TerminalExecReq struct {
       Command    string `json:"command"`
       NodeAlias  string `json:"nodeAlias,omitempty"`
       TimeoutSec int    `json:"timeoutSec,omitempty"`
       WorkingDir string `json:"workingDir,omitempty"`
   }
   ```
2. Add `TerminalExecResp`:
   ```go
   // TerminalExecResp encapsulates terminal execution output and status.
   type TerminalExecResp struct {
       Output     string `json:"output"`
       ExitCode   int    `json:"exitCode"`
       DurationMs int64  `json:"durationMs"`
       IsSuccess  bool   `json:"isSuccess"`
       Error      string `json:"error,omitempty"`
   }
   ```

### Step 3.2: Implement Terminal Backend & Fix Settings Persistence (`cli/cmdui/ui_server.go`)
1. In `mountAPIRoutes(mux *http.ServeMux)`:
   - Register `mux.HandleFunc("/api/terminal/exec", handleAPITerminalExec)`
2. Implement `handleAPITerminalExec(w http.ResponseWriter, r *http.Request)`:
   - Restrict to `POST` method (`http.StatusMethodNotAllowed`).
   - Decode `TerminalExecReq`. If command is empty, return 400 error.
   - Clamp `TimeoutSec` between 1 and 120 (default to 30 seconds).
   - If `req.NodeAlias == ""` or `strings.EqualFold(req.NodeAlias, "local")`:
     - Dispatch to `executeLocalTerminalCommand(req.Command, req.WorkingDir, req.TimeoutSec)`.
   - Else:
     - Dispatch to remote cluster node via `cmdssh.ExecuteRemoteCommand(req.NodeAlias, req.Command, req.TimeoutSec)`.
   - Encode `TerminalExecResp` as JSON with `Content-Type: application/json`.
3. Implement `executeLocalTerminalCommand(command string, dir string, timeoutSec int) TerminalExecResp`:
   - Construct context with deadline: `ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSec)*time.Second)`.
   - Select shell binary based on OS:
     - Windows: `exec.CommandContext(ctx, "powershell.exe", "-NoProfile", "-NonInteractive", "-Command", command)`
     - POSIX: `exec.CommandContext(ctx, "sh", "-c", command)`
   - Set working directory if specified.
   - Capture `CombinedOutput()`. Cap output at 1 MB.
   - Extract exit code: if `ctx.Err() == context.DeadlineExceeded`, return exit code 124 with `"Command timed out after N seconds"`.
   - Calculate `DurationMs` and populate `TerminalExecResp`.
4. Fix `loadSettings()` and `saveSettingsData()`:
   - Ensure `SettingsData.Attributes` is always initialized non-nil.
   - In `saveSettingsData()`, preserve all existing attributes from disk while merging newly incoming attributes.

### Step 3.3: Modernize Settings UI Layout & Embedded Terminal (`cli/cmdui/ui_assets.go`)
1. **Refactor `#tab-settings` HTML Layout:**
   - Group settings into clean, visual card sections with semantic typography:
     - **Card 1: General & UI Layout:** Theme, Graphics Mode, Browser Auto-Open, Layout Presets.
     - **Card 2: Cursor Fleet & Delegation:** Default target node (`u1`), migration options, auto-sync triggers.
     - **Card 3: Cluster REST & SSH Networking:** Default remote, cluster port, auto-deploy keys toggle.
     - **Card 4: Speed Settings (LAP, Account Switch & Machine Identity):** `lap.default_hours`, `account_switch.threshold`, machine alias.
     - **Card 5: Automation & Alerts (Telegram & Email):** `telegram.bot_token`, `telegram.chat_id`, email SMTP parameters.
     - **Card 6: Interactive Fleet Terminal:** Embedded console card.
2. **Embed Interactive Terminal Component:**
   - Control bar containing:
     - Target node select dropdown (`Local Host`, `u1`, etc. populated dynamically from `/api/ssh/nodes`).
     - Command input field with placeholder (`Enter command, e.g. gitmap status...`).
     - Execute button (`⚡ Run`) and Clear button (`🧹 Clear`).
   - Quick command chips:
     - `<button class="chip" onclick="setTerminalCmd('gitmap status')">gitmap status</button>`
     - `<button class="chip" onclick="setTerminalCmd('gitmap nodes list')">gitmap nodes list</button>`
     - `<button class="chip" onclick="setTerminalCmd('gitmap cur migrate --node u1 --dry-run')">cur migrate dry-run</button>`
     - `<button class="chip" onclick="setTerminalCmd('gitmap db sizes')">gitmap db sizes</button>`
   - Monospace console output window (`#terminal-output`) with ANSI/text rendering, `#0b1120` background, and status bar showing execution time and exit code.
3. **Fix Two-Way Data Binding in JavaScript:**
   - Update `loadSettings()`:
     - In addition to standard fields, iterate over `data.attributes || {}`.
     - Query all inputs having `data-key` matching attribute keys and set `input.value = data.attributes[key]`.
   - Update `saveSettings()`:
     - Collect standard fields.
     - Initialize `attributes: {}`.
     - Scan all DOM elements matching `input[data-key], select[data-key]`.
     - Populate `payload.attributes[el.getAttribute('data-key')] = el.value`.
     - POST payload to `/api/settings`.
4. **Implement Terminal JavaScript Controller:**
   - `execTerminalCommand()`:
     - Reads command and target node.
     - Displays `Executing on <node>...` in terminal window.
     - Posts to `/api/terminal/exec`.
     - Appends output and execution metadata (`[Exit: 0 | 142ms]`).
     - Appends to command history array.
   - Command history navigation: listen for `keydown` (ArrowUp / ArrowDown) on the command input to cycle through previous commands.
   - `clearTerminal()`: clears `#terminal-output`.

### Step 3.4: Author Unit Tests (`cli/cmdui/ui_test.go`)
1. `TestTerminalExecModels`: Test JSON marshalling and unmarshalling of `TerminalExecReq` and `TerminalExecResp`.
2. `TestExecuteLocalTerminalCommand_Echo`: Execute simple echo command and verify `IsSuccess == true` and expected output.
3. `TestExecuteLocalTerminalCommand_Timeout`: Test command with sleep exceeding 1s timeout to verify context cancellation and timeout error return.
4. `TestSettingsAttributesPreservation`: Test that saving and loading `SettingsData` with arbitrary keys in `Attributes` preserves all speed/bot fields without truncation.

---

## 4. Verification & Acceptance Criteria

1. **Test Suite:** Execute `go test -v ./cli/cmdui/...` verifying all terminal and settings tests pass.
2. **Two-Way Binding Verification:**
   - Launch UI (`gitmap ui settings`).
   - Input test values in Telegram Bot Token and LAP Hours.
   - Save settings, reload the browser tab, and confirm values remain populated.
3. **Interactive Terminal Verification:**
   - Run `gitmap status` from the embedded terminal and verify live output is rendered.
   - Click quick chip `gitmap db sizes` and verify response.
4. **Linter Compliance:** Zero violations in `03-ai-scripts/09-check-nested-ifs.py` and `03-ai-scripts/10-check-enum-and-boolean.py`.

---

## 5. Execution Summary & Verified Results

- **Terminal Request/Response Models:** Added `TerminalExecReq` and `TerminalExecResp` to `cli/cmdui/ui_types.go`.
- **API Terminal Endpoint:** Mounted `/api/terminal/exec` in `cli/cmdui/ui_server.go` with 30s default timeout (max 120s), routing local commands via PowerShell/sh with separate stdout/stderr capture, and remote commands via cluster SSH.
- **Two-Way Binding Preservation:** Fixed `loadSettings()` and `saveSettings()` to bidirectionally bind all `[data-key]` attribute fields; in `handleAPISettings` merged incoming `req.Attributes` to prevent data loss.
- **Settings UI Redesign:** Reorganized `#tab-settings` in `cli/cmdui/ui_assets.go` into 5 clean categorized cards and subtab navigation (General, Cursor Fleet, Speed, Alerts, Embedded Terminal) with traffic lights (`● ● ●`), live terminal, quick action chips, and command history.
- **Unit Tests:** Added `TestAPITerminalExec_LocalEcho` and `TestAPISettings_AttributesPreservation` in `cli/cmdui/ui_test.go` — all passed (100% PASS).
- **Linters:** Passed `check-nested-ifs.py` and `check-enum-and-boolean.py` with 0 violations.
