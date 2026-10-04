# Windows to Ubuntu AGM Migration Guide

## 1. Overview & Architecture

This guide details the procedure for exporting Antigravity Manager (AGM) settings, registered Google accounts, databases, and configuration from a Windows workstation and migrating them seamlessly to an Ubuntu Linux VM (`u1` at `192.168.1.22`).

### Path & Directory Mappings
| Component | Windows Source Path | Ubuntu Target Path |
| :--- | :--- | :--- |
| AGM Config Root | `%USERPROFILE%\.antigravity_tools` | `~/.antigravity_tools` |
| Accounts Index | `%USERPROFILE%\.antigravity_tools\accounts.json` | `~/.antigravity_tools/accounts.json` |
| Accounts Profiles | `%USERPROFILE%\.antigravity_tools\accounts\*.json` | `~/.antigravity_tools/accounts/` |
| SQLite Databases | `user_tokens.db`, `security.db` | `~/.antigravity_tools/*.db` |
| IDE Credentials | `%USERPROFILE%\.gemini\oauth_creds.json` | `~/.gemini/oauth_creds.json` |
| Global Storage | `%APPDATA%\Antigravity\User\globalStorage` | `~/.config/Antigravity/User/globalStorage` |
| Binary Launcher | `%LOCALAPPDATA%\Programs\Antigravity` | `~/.local/share/antigravity-ide` |

---

## 2. Windows Settings Packaging

From PowerShell on Windows:
```powershell
$ExportDir = "$env:TEMP\agm_windows_export"
New-Item -ItemType Directory -Force -Path $ExportDir
Copy-Item "$env:USERPROFILE\.antigravity_tools\accounts.json" -Destination $ExportDir
Copy-Item "$env:USERPROFILE\.antigravity_tools\accounts" -Recurse -Destination $ExportDir
Copy-Item "$env:USERPROFILE\.antigravity_tools\gui_config.json" -Destination $ExportDir -ErrorAction SilentlyContinue
Copy-Item "$env:USERPROFILE\.antigravity_tools\update_settings.json" -Destination $ExportDir -ErrorAction SilentlyContinue
Copy-Item "$env:USERPROFILE\.antigravity_tools\*.db" -Destination $ExportDir -ErrorAction SilentlyContinue

tar -czf "$env:TEMP\agm_windows_export.tar.gz" -C $ExportDir .
```

---

## 3. Remote Transfer via SCP & Extraction

Transfer the archive to the Ubuntu node:
```bash
scp agm_windows_export.tar.gz a@192.168.1.22:/tmp/agm_windows_export.tar.gz
```

On Ubuntu node `u1`:
```bash
mkdir -p /home/a/.antigravity_tools/accounts
tar -xzf /tmp/agm_windows_export.tar.gz -C /home/a/.antigravity_tools
chmod 700 /home/a/.antigravity_tools
chmod 600 /home/a/.antigravity_tools/*.json /home/a/.antigravity_tools/*.db 2>/dev/null
chmod 600 /home/a/.antigravity_tools/accounts/*.json 2>/dev/null
```

---

## 4. Key Root Cause Analyses (RCAs) Solved

### A. Launcher Script Infinite Recursion Bug (`Argument list too long`)
- **Symptom:** When AGM triggers account switching, it executes `/home/a/.local/bin/antigravity`, which aborts with:
  `/home/a/.local/bin/antigravity: line 10: /home/a/.local/bin/antigravity: Argument list too long`.
- **Root Cause:** The default launcher script has:
  ```bash
  DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
  EXEC="$DIR/antigravity"
  exec "$EXEC" --no-sandbox "$@"
  ```
  When the launcher script is installed in `/home/a/.local/bin/` or `/usr/local/bin/`, `$DIR/antigravity` points to the script itself, causing infinite self-recursion until Linux hits `E2BIG`.
- **Solution:** Deployed `scripts/antigravity-launcher.sh` which explicitly checks for the real ELF executable in `/home/a/.local/share/antigravity-ide/antigravity` or `/opt/antigravity`.

### B. Linux Secret Service Keyring Timeout
- **Symptom:** In headless SSH sessions without an unlocked desktop keyring, `agm switch` times out for 10 seconds attempting to write to `secret-tool store --collection=login`.
- **Root Cause:** `Antigravity-Manager/src-tauri/src/modules/integration.rs` aborts the switch if `secret-tool` fails, even though file-based sync to `~/.gemini/oauth_creds.json` already succeeded.
- **Solution:** Ensured `libsecret-tools` is installed and `~/.gemini/` file credentials take precedence.

---

## 5. Verification

Verify account switching on Ubuntu node `u1`:
```bash
/usr/bin/agm switch rokixshohag1@gmail.com
```

Verify token synchronization:
```bash
cat /home/a/.gemini/oauth_creds.json | jq .
```
Result: Returns active OAuth bearer tokens, refresh tokens, and Google user ID verified.
