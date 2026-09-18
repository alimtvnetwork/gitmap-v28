## Quick Install v6.215.0

### Windows (PowerShell 5.1+)
```powershell
irm https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.215.0/install.ps1 | iex
```

### Linux / macOS (Bash)
```bash
curl -fsSL https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.215.0/install.sh | bash
```

## Changelog v6.215.0

- Aggregate all pipeline section and step failure errors together in one consolidated view
- Explicitly state saved error log files (.gitmap/pipeline/pipeline_errors.log, <runId>.log, and DB path)
- Display full execution metadata (which workflow, run ID, branch, commit, when run timestamp and relative age, run duration, and URL)
