## Quick Install v6.218.0

### Windows (PowerShell 5.1+)
```powershell
irm https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.218.0/install.ps1 | iex
```

### Linux / macOS (Bash)
```bash
curl -fsSL https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.218.0/install.sh | bash
```

## Changelog v6.218.0

- Add gitmap pipelines alias and errorlogs/errors subcommands with negative index inspection (-N)
- Implement incremental SQLite caching for last-failed-logs avoiding duplicate log downloads
- Correlate GitHub Actions job API data to guarantee all failing sections are combined
