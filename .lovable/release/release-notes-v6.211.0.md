## Quick Install v6.211.0

### Windows (PowerShell 5.1+)
```powershell
irm https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.211.0/install.ps1 | iex
```

### Linux / macOS (Bash)
```bash
curl -fsSL https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.211.0/install.sh | bash
```

## Changelog v6.211.0

- Add targeted test runner by file path, file name, or Go package name
- Migrate CI/CD runner artifacts to cross-platform OS temp folder
- Add isolated per-test failure log files and hashed session folders
- Flatten nested-if conditionals across mkdir.go and store.go
