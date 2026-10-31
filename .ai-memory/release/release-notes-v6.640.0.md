# GitMap v6.640.0

## What's Changed in v6.640.0

- **Dual Sync Script Versioning:** Full execution parity for fleet sync between V6 unrestricted mode (`38-sync-prompts-skills-scripts.py`) and V7 restricted mode (`38-sync-prompts-skills-scripts-v7.py` / `--v7`), validating morning commit windows (< 08:00 AM on weekdays; 24h weekends), timezone resolution, and bot credentials.
- **GitMap Cat `-t` Tail & Combined Query Windowing:** Added `-t` shorthand for `--tail` in `cli/cmdmacro/cat_options.go` and empowered `gitmap cat` to seamlessly evaluate query searches over tail and line range windows, directly replacing PowerShell `Get-Content ... -Tail N | Select-String` and Node.js file read scripts.
- **GitMap Py Anti-Pattern Detection:** Enhanced `cli/cmdpy/py_cmd.go` with inline code guards warning AI agents against ad-hoc Python filesystem crawling (`os.walk`), sleep loops (`time.sleep`), and file reading (`open().read()`), routing them to high-speed native GitMap commands.
- **Authoritative Anti-Pattern Ban Matrix:** Updated `.agents/skills/gitmap/skill.md` and synced to `coding-guidelines` with side-by-side replacements for the 7 observed agent anti-patterns.
- **Fleet Synchronization:** Synchronized updated skills and sync scripts to `coding-guidelines` and connected repositories (`antigravity-manager`).

### Quick Install One-Liners

**Windows (PowerShell):**
```powershell
irm https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.640.0/install.ps1 | iex
```

**Linux / macOS (Bash):**
```bash
curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.640.0/install.sh | sh
```
