## Quick Install v6.212.0

### Windows (PowerShell 5.1+)
```powershell
irm https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.212.0/install.ps1 | iex
```

### Linux / macOS (Bash)
```bash
curl -fsSL https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.212.0/install.sh | bash
```

## Changelog v6.212.0

- Prevent infinite dynamic timeline polling loop in tests and CI mode
- Skip live GitHub API polling in unit test suite under CI
- Verify all 20+ Go packages pass with zero failures
