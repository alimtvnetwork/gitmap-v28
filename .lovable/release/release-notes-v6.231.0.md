## Quick Install v6.231.0

### Windows (PowerShell 5.1+)
```powershell
irm https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.231.0/install.ps1 | iex
```

### Linux / macOS (Bash)
```bash
curl -fsSL https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.231.0/install.sh | bash
```

## Changelog v6.231.0

- Profile installation idempotency with component tree rendering
- Untracked SQLite database auditing for profiles and packages with full stack trace preservation
- Cross-platform VMware CLI support across Linux and Windows
- Resilient VMware shared folder mounting and startup persistence
- Terminal table column alignment fixes for gitmap pull and gitmap status
