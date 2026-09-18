## Quick Install v6.259.0

### Windows (PowerShell 5.1+)
```powershell
irm https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.259.0/install.ps1 | iex
```

### Linux / macOS (Bash)
```bash
curl -fsSL https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.259.0/install.sh | bash
```

## Changelog v6.259.0

- feat(pipeline): parallel section and job log downloads for faster single-commit pipeline error inspection
- feat(pipeline): two-pass non-mutating line execution and filtering algorithms for error log processing
- feat(pipeline): resilient historical fallback retrieval for pipeline DB and commit logs
- fix(pipeline): prevent historical fallback from overwriting clean success pipeline payloads
- fix(ci): resolve gocritic switch statement violations and ResultSlice method wrappers
