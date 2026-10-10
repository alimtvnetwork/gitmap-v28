# GitMap v6.531.0

## What's Changed in v6.531.0

### GitMap Repo-Cache Scan, Clone, and LS Ecosystem
- **`gitmap scan . --rc` (and `gitmap scan --rc`):**
  - Export and auto-merge discovered repositories directly into `repo-cache/01-gitmap/gitmap.json`.
  - Normalizes repository URLs (HTTPS/SSH) and performs smart in-place metadata updates (branch, head SHA, commit timestamps).
- **`gitmap scan . --rc --separate` (aliases `--seprate`, `--sep`, `--s`):**
  - Dynamically scans `repo-cache/` and allocates the next sequential numbered manifest (`02-gitmap.json`, `03-gitmap.json`, etc.).
- **`gitmap clone/cfr/cfrp --rc` (and positional `rc`):**
  - Inspects discovered manifests (prioritizing `01-gitmap/gitmap.json` first, then sequential manifests).
  - High-performance Pure Go SQLite Split-DB cache at `.gitmap/data/repocache/sql.db` (<5ms response time with mtime & SHA256 cache invalidation).
  - Clean short TreeView rendering displaying concise `owner/repo` paths instead of bloated URLs.
  - Interactive selection menu or non-interactive batch cloning via `-y` / `--yes`.
- **`gitmap ls rc` (and `gitmap list rc`):**
  - High-signal manifest summary with protocol breakdown (`[SSH]` vs `[Public HTTPS]`).
  - Copy-pasteable import command recommendations.

### Quick Install One-Liners

**Windows (PowerShell):**
```powershell
irm https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.531.0/install.ps1 | iex
```

**Linux / macOS (Bash):**
```bash
curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.531.0/install.sh | sh
```
