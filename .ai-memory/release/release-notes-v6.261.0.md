## Quick Install v6.261.0

### Windows (PowerShell 5.1+)
```powershell
irm https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.261.0/install.ps1 | iex
```

### Linux / macOS (Bash)
```bash
curl -fsSL https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.261.0/install.sh | bash
```

## Changelog v6.261.0

- feat(pipeline): bounded pattern-aware stack trace extraction and job-scoped isolation
- fix(pipeline): halt context capture on exit code lines to eliminate runner cleanup noise
- fix(pipeline): eliminate duplicated section reports and filter tool installation progress tickers
- feat(ssh): multi-command discovery, machine join parity, and terminal help verification
- fix(ci): resolve nested ifs, ssh help exit panic, absolute path linters, and error management regex
