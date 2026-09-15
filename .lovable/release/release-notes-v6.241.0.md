## Quick Install v6.241.0

### Windows (PowerShell 5.1+)
```powershell
irm https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.241.0/install.ps1 | iex
```

### Linux / macOS (Bash)
```bash
curl -fsSL https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.241.0/install.sh | bash
```

## Changelog v6.241.0

- Add non-Linux launcher stub for cross-platform compilation
- Remove unused functions and fix SA4023 interface nilness
- Fix SQLite DB mock reconnection in cmdssh tests
- Correct unit test assertions across install, pull, and ssh suites
