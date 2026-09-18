## Quick Install v6.232.0

### Windows (PowerShell 5.1+)
```powershell
irm https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.232.0/install.ps1 | iex
```

### Linux / macOS (Bash)
```bash
curl -fsSL https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.232.0/install.sh | bash
```

## Changelog v6.232.0

- Commit-scoped offset resolution for gitmap pipeline errors (-1, -2, -3)
- Added gitmap pipeline history command with tree visualization
- Added gitmap pipeline logs command with clipboard and file export
- Enriched repository SQLite database telemetry (.gitmap/data/pipeline.db)
- Added gitmap storage and gitmap os storage command suite
