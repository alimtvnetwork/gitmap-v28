## Quick Install v6.243.0

### Windows (PowerShell 5.1+)
```powershell
irm https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.243.0/install.ps1 | iex
```

### Linux / macOS (Bash)
```bash
curl -fsSL https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.243.0/install.sh | bash
```

## Changelog v6.243.0

- Convert join command to universal AppError envelopes
- Introduce CheckHelpOrEmpty centralized DRY help checking engine
- Eliminate repeated help checks across cluster, server-cmd, and servers-clients
