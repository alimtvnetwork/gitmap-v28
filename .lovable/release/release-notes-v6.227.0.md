## Quick Install v6.227.0

### Windows (PowerShell 5.1+)
```powershell
irm https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.227.0/install.ps1 | iex
```

### Linux / macOS (Bash)
```bash
curl -fsSL https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.227.0/install.sh | bash
```

## Changelog v6.227.0

- Fix Installer Dry-Run Windows step in release.yml by updating path to cli/scripts/install.ps1
- Fix Installer Smoke Windows seed data download by updating GitHub raw data URLs to cli/data/
- Add local repository fallback for seed files in install.ps1 and install.sh
- Synchronize installer URLs across Go constants, helptext, and scripts
