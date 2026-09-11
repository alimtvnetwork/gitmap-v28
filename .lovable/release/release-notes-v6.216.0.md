## Quick Install v6.216.0

### Windows (PowerShell 5.1+)
```powershell
irm https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.216.0/install.ps1 | iex
```

### Linux / macOS (Bash)
```bash
curl -fsSL https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.216.0/install.sh | bash
```

## Changelog v6.216.0

- Filter out superseded failures when a workflow has succeeded in a newer run
- Eliminate fallback that dredged up ancient failed runs when current pipeline is green
- Clear local last_error.log and pipeline_errors.log when pipeline runs pass cleanly
