# Subtask 217.1: Antigravity Profile & Preset Forensics Specification

- **Parent Plan:** `pending/217-antigravity-fleet-parity-theme-preset-plugins-and-delegation.md`
- **Spec Reference:** [01-architecture-spec.md](../../../../02-spec/21-app/217-antigravity-fleet-parity-theme-preset-plugins-and-delegation/01-architecture-spec.md)
- **Status:** Pending
- **Target Area:** `~/.gemini/config/config.json`, `~/.gemini/config/plugins/`, `~/.gemini/config/projects/*.json`, `~/.config/Antigravity/User/settings.json`, `~/.antigravity_tools/instances/instances.json`

---

## 1. Objective

Formulate, document, and verify the exhaustive configuration mapping from the master Windows Antigravity default instance (`C:\Users\Administrator\...`) to the canonical Linux fleet node environment (`/home/a/...` on Ubuntu `u1`). This ensures complete visual, permission, plugin, and skill parity, eradicating unconfigured fallbacks and filesystem path leakages.

---

## 2. Forensic Discoveries & Root Cause Analysis

### 2.1 UI Theme Mismatch & Color Seed Isolation
- **Forensic Discovery:** Antigravity's conversation UI does not derive its primary background and accent palette from standard VS Code themes. Instead, it parses `customThemeSeedsDark` inside `~/.gemini/config/config.json`.
- **Root Cause:** When Ubuntu `u1` was initialized, only bare git repository trees were transferred. The master `config.json` containing the Dracula purple seed (`#BD93F9`) and deep charcoal background (`#19191C`) was never deployed, causing Antigravity's webview to load default unbranded CSS variables.

### 2.2 The "Default" Permission Preset Fallback
- **Forensic Discovery:** On Ubuntu `u1`, the permission preset badge displays `Default` instead of an active unattended execution mode, forcing every terminal execution and file modification to prompt for manual user confirmation.
- **Root Cause:**
  1. The underlying `language_server` binary reads permissions from two tiers: `~/.gemini/config/config.json` (global `userSettings`) and `~/.gemini/config/projects/<id>.json` (project-specific `settings`).
  2. In the deployed state on `u1`, `projects/*.json` contained empty object literals: `"settings": {}` and `"permissionGrants": { "allow": [] }`.
  3. In `settings.json`, `"antigravity.turboMode"` was set to `false`.
  4. Without explicit policies defined, `language_server`'s protobuf evaluation logic defaults to `CASCADE_COMMANDS_AUTO_EXECUTION_OFF` and prompts on all actions, rendering the preset as "Default".

### 2.3 Total Plugins & Skills Absence
- **Forensic Discovery:** The directory `/home/a/.gemini/config/plugins/` on Ubuntu `u1` was completely missing or empty (0 plugins, 0 skills).
- **Root Cause:** Antigravity does not automatically download plugins across offline instances without explicit user marketplace interaction. The 4 official core plugins and their 43 nested skills must be bundled and deployed directly from the master instance.

### 2.4 Windows Backslash Path Leakage on Linux
- **Forensic Discovery:** The filesystem on Ubuntu `u1` contained an anomaly directory literally named `/home/a/C:\Users\Administrator\AppData\Roaming\Antigravity`.
- **Root Cause:** `~/.antigravity_tools/instances/instances.json` contained Windows backslash paths:
  `"data_dir": "C:\\Users\\Administrator\\AppData\\Roaming\\Antigravity"`. On Linux, backslashes are valid literal characters in file and directory names. Any tool executing directory operations created the literal Windows path on disk, failing to locate user configurations at `/home/a/.config/Antigravity`.

---

## 3. Configuration Mapping & Normalization Specification

### 3.1 Filesystem Path Transformation Matrix

| Component | Source Path (Windows Default Instance) | Target Path (Ubuntu Fleet Node `u1`) |
| :--- | :--- | :--- |
| **Global Config** | `C:\Users\Administrator\.gemini\config\config.json` | `/home/a/.gemini/config/config.json` |
| **Official Plugins** | `C:\Users\Administrator\.gemini\config\plugins\` | `/home/a/.gemini/config/plugins/` |
| **Project Descriptors** | `C:\Users\Administrator\.gemini\config\projects\*.json` | `/home/a/.gemini/config/projects/*.json` |
| **VS Code User Settings** | `C:\Users\Administrator\AppData\Roaming\Antigravity\User\settings.json` | `/home/a/.config/Antigravity/User/settings.json` |
| **Instance Registry** | `C:\Users\Administrator\.antigravity_tools\instances\instances.json` | `/home/a/.antigravity_tools/instances/instances.json` |
| **Workspace Repos** | `D:\work\<repo>` | `/home/a/work/<repo>` |

### 3.2 Global Config JSON Specification (`config.json`)
The Linux target `/home/a/.gemini/config/config.json` must be written with the following exact payload:

```json
{
  "plugins": {
    "chrome-devtools-plugin": {
      "enabled": true,
      "installedFrom": {
        "id": "antigravity-plugins-official/plugin/chrome-devtools-plugin",
        "marketplace": "antigravity-plugins-official",
        "version": "0.21.0"
      }
    },
    "data-agent-kit-plugin": {
      "enabled": true,
      "installedFrom": {
        "id": "antigravity-plugins-official/plugin/data-agent-kit-plugin",
        "marketplace": "antigravity-plugins-official",
        "version": "0.7.0"
      }
    },
    "google-antigravity-sdk": {
      "enabled": true,
      "installedFrom": {
        "id": "antigravity-plugins-official/plugin/google-antigravity-sdk",
        "marketplace": "antigravity-plugins-official",
        "version": "0.0.9"
      }
    },
    "modern-web-guidance-plugin": {
      "enabled": true,
      "installedFrom": {
        "id": "antigravity-plugins-official/plugin/modern-web-guidance-plugin",
        "marketplace": "antigravity-plugins-official",
        "version": "1.0.6"
      }
    }
  },
  "userSettings": {
    "artifactReviewMode": "ARTIFACT_REVIEW_MODE_TURBO",
    "autoExecutionPolicy": "CASCADE_COMMANDS_AUTO_EXECUTION_EAGER",
    "browserJsExecutionPolicy": "BROWSER_JS_EXECUTION_POLICY_TURBO",
    "conversationWidth": "CONVERSATION_WIDTH_WIDE",
    "customThemeSeedsDark": {
      "background": "#19191C",
      "foregroundOverride": "#F8F8F2",
      "primary": "#BD93F9"
    },
    "customThemeSeedsLight": {
      "background": "#EAECF0",
      "foregroundOverride": "#202021",
      "primary": "#8839EF"
    },
    "enableTerminalSandbox": false,
    "globalPermissionGrants": {
      "allow": [
        "command(git status && git remote -v && git log --oneline -5 && git config user.name && git config user.email)",
        "read_file(/home/a/work)",
        "write_file(/home/a/work)",
        "execute_url(*)",
        "read_url(prnt.sc)",
        "read_url(*)"
      ]
    },
    "nonWorkspaceFileAccessPolicy": "AGENT_SETTING_POLICY_ALLOW",
    "queuedMessageDeliveryStrategy": "MESSAGE_DELIVERY_STRATEGY_WHEN_IDLE",
    "remoteControlEnabled": true,
    "remoteControlHostname": "u1",
    "themeMode": "THEME_MODE_DARK",
    "useAiCredits": false,
    "verboseAgentChat": false
  }
}
```

### 3.3 Project Descriptors Normalization Specification
Every file matching `/home/a/.gemini/config/projects/*.json` must be updated to inject explicit auto-execution policies:
1. `settings.autoExecutionPolicy` $\to$ `"CASCADE_COMMANDS_AUTO_EXECUTION_EAGER"`
2. `settings.fileAccessPolicy` $\to$ `"AGENT_SETTING_POLICY_ALLOW"`
3. `settings.sandboxMode` $\to$ `false`
4. Path remap in `projectResources`: Replace `file:///d%3A/work/` with `file:///home/a/work/`.

### 3.4 4 Plugins & 43 Skills Packaging Specification
The deployment engine must archive and sync:
- `~/.gemini/config/plugins/chrome-devtools-plugin/`
- `~/.gemini/config/plugins/data-agent-kit-plugin/`
- `~/.gemini/config/plugins/google-antigravity-sdk/`
- `~/.gemini/config/plugins/modern-web-guidance-plugin/`
- Verify that exactly 43 `SKILL.md` / `skill.md` files exist across the extracted directories.

### 3.5 Instance Registry Sanitization Specification
In `/home/a/.antigravity_tools/instances/instances.json`:
1. `data_dir` must strictly be set to `"/home/a/.config/Antigravity"`.
2. Any anomaly directories matching `/home/a/C:\*` or `/home/a/C:*` must be removed via safe recursive purge.

---

## 4. Verification & Acceptance Criteria

- **Zero Windows Backslash Directories:** Executing `ls -d /home/a/C:* 2>/dev/null` on `u1` yields exit code 1 (no files found).
- **Theme Palette Verified:** Reading `/home/a/.gemini/config/config.json` confirms `primary: "#BD93F9"` and `background: "#19191C"`.
- **4 Plugins Present:** `ls /home/a/.gemini/config/plugins` lists all four plugin directory names.
- **43 Skills Count:** Running `find /home/a/.gemini/config/plugins -name "*skill*.md" | wc -l` yields exactly `43`.
- **Project Settings Injected:** `grep -L "CASCADE_COMMANDS_AUTO_EXECUTION_EAGER" /home/a/.gemini/config/projects/*.json` yields 0 files.
- **Preset Status:** Antigravity IDE UI displays active unattended execution preset instead of "Default".
