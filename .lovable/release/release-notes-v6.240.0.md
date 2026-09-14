## Quick Install v6.240.0

### Windows (PowerShell 5.1+)
```powershell
irm https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.240.0/install.ps1 | iex
```

### Linux / macOS (Bash)
```bash
curl -fsSL https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.240.0/install.sh | bash
```

## Changelog v6.240.0

- Resolve duplicate identifiers in store and cmdssh packages
- Eliminate swallowed SQLite migration errors using isBenignAlterError
- Flatten control-flow and eliminate nested if blocks in cluster bootstrap
- Fix undefined PullResult reference and missing resolveJoinAlias argument
- Remove unused fmt import in Linux antigravity deploy
