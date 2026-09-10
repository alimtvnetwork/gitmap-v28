## Quick Install v6.209.0

### Windows (PowerShell 5.1+)
```powershell
irm https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.209.0/install.ps1 | iex
```

### Linux / macOS (Bash)
```bash
curl -fsSL https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.209.0/install.sh | bash
```

## Changelog v6.209.0

- Pass active install directory from gitmap update to remote installers
- Align PowerShell installer default directory with SSoT (gitmap-cli)
- Add Repair-LegacyLayout to prune legacy gitmap directory and stale PATH
- Add trailing visual spacing gaps to post-update and installer summaries
- Register gitmap binary and gitmap info commands with full identity reporting
