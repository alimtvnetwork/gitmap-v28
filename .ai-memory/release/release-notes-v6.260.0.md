## Quick Install v6.260.0

### Windows (PowerShell 5.1+)
```powershell
irm https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.260.0/install.ps1 | iex
```

### Linux / macOS (Bash)
```bash
curl -fsSL https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.260.0/install.sh | bash
```

## Changelog v6.260.0

- feat(storage): clean storage ls table UI with root folder header and relative database paths
- feat(table): utilize reusable termtable auto-alignment and spacing framework
- fix(ci): complete verification across all 46 CI/CD local quality gates and 5 remote workflows
- refactor(storage): decompose database listing to storage_ls.go to maintain <100 line canonical file sizes
