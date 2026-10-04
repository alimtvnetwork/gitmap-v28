# Subtask 217.02: PowerShell Full Parity Sync Script (`sync-antigravity-full-profile.ps1`)

**Subtask Code:** 217.02  
**Parent Task:** `217-antigravity-fleet-parity-theme-preset-plugins-and-delegation`  
**Owner:** Worker 02  
**Target File:** `repo-secrets/04-ubuntu-migration/sync-antigravity-full-profile.ps1`  
**Status:** Ready for Implementation  
**Date:** 2026-10-05  

---

## 1. Context & Objective

During initial migration of workstation repositories to Ubuntu node `u1`, the Antigravity IDE configuration was only partially synchronized. This resulted in:
1. Missing `config.json`, leaving UI theme unbranded and permission presets at the restrictive "Default" prompt mode.
2. Missing plugins (`~/.gemini/config/plugins/` was empty).
3. Missing 43 skills from the agent prompt HUD.
4. Serialized Windows backslashes causing a stray `/home/a/<windows-appdata>\...` folder to be created on Linux.
5. Missing SUID permissions on `chrome-sandbox`.

The objective of this subtask is to author a complete, standalone, idempotent PowerShell script `repo-secrets/04-ubuntu-migration/sync-antigravity-full-profile.ps1` that packages the complete Antigravity profile from Windows, sanitizes all paths for Linux, streams the payload over SSH to `u1`, applies elevated system hardening, and restarts the IDE with full theme, preset, and plugin parity.

---

## 2. Technical Requirements & Parameter Specifications

The script MUST expose the following parameters with positive boolean switches:

```powershell
[CmdletBinding()]
param(
    [Parameter(Position = 0)]
    [string]$TargetHost = "u1",

    [Parameter(Position = 1)]
    [string]$TargetUser = "a",

    [Parameter()]
    [ValidateSet("turbo", "eager", "default")]
    [string]$Preset = "turbo",

    [Parameter()]
    [ValidateSet("dark-dracula", "dark-default", "light-default")]
    [string]$Theme = "dark-dracula",

    [Parameter()]
    [switch]$SyncPlugins = $true,

    [Parameter()]
    [switch]$SyncSkills = $true,

    [Parameter()]
    [switch]$SanitizePaths = $true,

    [Parameter()]
    [switch]$RestartIDE = $true,

    [Parameter()]
    [string]$WindowsGeminiDir = "$HOME\.gemini",

    [Parameter()]
    [string]$WindowsToolsDir = "$HOME\.antigravity_tools",

    [Parameter()]
    [string]$RemoteGeminiDir = "/home/a/.gemini",

    [Parameter()]
    [string]$RemoteToolsDir = "/home/a/.antigravity_tools",

    [Parameter()]
    [string]$RemoteIdeInstallDir = "/home/a/.local/share/antigravity-ide",

    [Parameter()]
    [switch]$DryRun,

    [Parameter()]
    [switch]$Force
)
```

---

## 3. Step-by-Step Implementation Architecture

### Step 1: Pre-Flight Connectivity & Authentication Probe
- Test reachability via `ssh -o BatchMode=yes -o ConnectTimeout=5 $TargetUser@$TargetHost "echo CONNECTED"`.
- If probe fails, display human-readable troubleshooting guidance (check network, verify OpenSSH service on `u1`, check SSH keys) and gracefully terminate.
- Detect sudo password requirement for elevated remote commands (`echo a | sudo -S ...`).

### Step 2: Ingest & Transform Local Windows Configuration
- **Read `config.json`**:
  - Ingest `<user-home>/.gemini/config/config.json`.
  - Ensure `customThemeSeedsDark` is set to Dracula (`#19191C` background, `#BD93F9` primary, `#F8F8F2` foreground).
  - Ensure `userSettings` contains `CASCADE_COMMANDS_AUTO_EXECUTION_EAGER`, `BROWSER_JS_EXECUTION_POLICY_TURBO`, and wide `globalPermissionGrants`.
  - Ensure all 4 plugins are enabled in the `plugins` map:
    - `chrome-devtools-plugin` (v0.21.0)
    - `data-agent-kit-plugin` (v0.7.0)
    - `google-antigravity-sdk` (v0.0.9)
    - `modern-web-guidance-plugin` (v1.0.6)
- **Sanitize `instances.json`**:
  - Ingest `<user-home>/.antigravity_tools/instances/instances.json`.
  - Replace `data_dir` Windows path with Linux path: `/home/a/.config/Antigravity`.
  - Replace `executable_path` with `/home/a/.local/share/antigravity-ide/antigravity`.
  - Strip all backslashes and drive letters.
- **Transform `projects/*.json`**:
  - Read all project descriptors in `<user-home>/.gemini/config/projects/`.
  - Inject `"settings": { "fileAccessPolicy": "AGENT_SETTING_POLICY_ALLOW", "autoExecutionPolicy": "CASCADE_COMMANDS_AUTO_EXECUTION_EAGER" }`.
  - Inject `"permissionGrants": { "allow": ["read_file(/home/a/git-work)", "write_file(/home/a/git-work)", "command(*)"] }`.
  - Normalize `folderUri` to `file:///<user-home>/git-work/<repo>`.

### Step 3: Bundle Plugins & Skills Directory Tree
- Verify source folder `<user-home>/.gemini/config/plugins` exists.
- Ingest all 4 plugin directories, manifests (`plugin.json`, `gemini-extension.json`), and all 43 skill subdirectories containing `SKILL.md`.
- Create a temporary staging archive (or in-memory tar.gz stream) containing:
  - `config.json`
  - `instances.json`
  - `projects/`
  - `plugins/` (with nested `skills/`)

### Step 4: Stream, Unpack & Elevate on Remote Node
- Stream the archive over SSH into a dedicated remote receiver script executing on `u1`.
- Unpack into:
  - `/home/a/.gemini/config/config.json`
  - `/home/a/.gemini/config/plugins/`
  - `/home/a/.gemini/config/projects/`
  - `/home/a/.antigravity_tools/instances/instances.json`
- Ensure file permissions and ownership: `chown -R a:a /home/a/.gemini /home/a/.antigravity_tools`.

### Step 5: Elevated Remote Hygiene & Path Sanitization
- Execute elevated bash cleanup via SSH:
  ```bash
  # Eradicate stray Windows path directories
  find /home/a -maxdepth 1 \( -name 'C:*' -o -name 'c:*' -o -name 'C:\\*' \) -print0 2>/dev/null | while IFS= read -r -d '' d; do
      echo a | sudo -S rm -rf "$d"
  done

  # Secure chrome-sandbox with SUID root
  SANDBOX="/home/a/.local/share/antigravity-ide/chrome-sandbox"
  if [ -f "$SANDBOX" ]; then
      echo a | sudo -S chown root:root "$SANDBOX"
      echo a | sudo -S chmod 4755 "$SANDBOX"
  fi
  ```

### Step 6: Antigravity IDE Restart & Parity Verification
- If `-RestartIDE` is true:
  - Gracefully terminate existing instances: `pkill -x antigravity 2>/dev/null || true`.
  - Relaunch IDE or verify background runner readiness.
- Emit structured summary table indicating:
  - Remote host reachable: Yes
  - Theme deployed: `dark-dracula` (`#BD93F9`)
  - Preset deployed: `turbo` (`CASCADE_COMMANDS_AUTO_EXECUTION_EAGER`)
  - Plugins synced: 4 (`chrome-devtools`, `data-agent-kit`, `google-antigravity-sdk`, `modern-web-guidance`)
  - Skills synced: 43
  - Stray Windows paths purged: Yes
  - Chrome sandbox SUID configured: Yes

---

## 4. Verification & Acceptance Checklist

- [ ] Script is saved at `repo-secrets/04-ubuntu-migration/sync-antigravity-full-profile.ps1`.
- [ ] Contains all specified parameters with defaults matching `u1` (host `u1`, user `a`, preset `turbo`, theme `dark-dracula`).
- [ ] Uses positive booleans (`$SyncPlugins`, `$SyncSkills`, `$SanitizePaths`, `$RestartIDE`).
- [ ] Correctly packages `config.json`, 4 plugins, 43 skills, and transformed project descriptors.
- [ ] Eliminates stray Windows path folders on Linux and configures SUID 4755 on `chrome-sandbox`.
- [ ] Powershell syntax passes validation with zero syntax errors.
