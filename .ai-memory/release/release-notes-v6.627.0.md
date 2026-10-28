# GitMap v6.627.0

## What's Changed in v6.627.0

- **Context Menu Icon Isolation & Legacy Purge**: Fixed icon display bug by removing duplicate child entry icon writes, ensuring the top-level GitMap context menu item crisply renders the official application icon, and automatically purging legacy orphaned registry keys (`70_git`, `95_install_wizard`, `96_clean_powershell`).
- **Show-Log Non-Git Directory Guard & Dashboard Fallback**: Resolved the non-git container folder error in `gitmap show-log`; when executed inside a parent directory with child repositories (e.g. `D:\work`), GitMap now displays child repository targets and automatically launches the GitMap Web UI dashboard instead of halting with an abrupt command-line error.
- **Default Browser Decoupling**: Configured system browser launching via native OS openers (`rundll32 url.dll,FileProtocolHandler`, `open`, `xdg-open`) without forcing separate Chrome profile user-data directories, controllable via `GITMAP_BROWSER_MODE` and Web UI settings.
- **Workstation Profiles, Single-Tool Installer & Macro UI**: Added dedicated context menu actions (`Macro ▸ Open Macro Builder UI`, `Install ▸ Open Installer Manager UI`), wired live Profile Tree viewer, single-tool installers/uninstallers, and macro execution/builder into the embedded Fleet Web UI.

### Quick Install One-Liners

**Windows (PowerShell):**
```powershell
irm https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.627.0/install.ps1 | iex
```

**Linux / macOS (Bash):**
```bash
curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.627.0/install.sh | sh
```
