## Quick Install v6.184.0

### Windows (PowerShell 5.1+)
```powershell
irm https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.184.0/install.ps1 | iex
```

### Linux / macOS (Bash)
```bash
curl -fsSL https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.184.0/install.sh | bash
```

## Changelog v6.184.0

- Enhanced GitMap CLI footer when inside a git repository to display repo name, remote Git URL, active branch, latest branch, open PR count, and comprehensive branch status
- Resolved real-time git branch tracking and upstream sync counters (ahead/behind/up to date)
- Automated SSoT version bumping pipeline via 03-ai-scripts/29-release-bumper.py
- Bumped version to v6.184.0 across all Single Source of Truth manifests
