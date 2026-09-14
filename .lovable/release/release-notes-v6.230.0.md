## Quick Install v6.230.0

### Windows (PowerShell 5.1+)
```powershell
irm https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.230.0/install.ps1 | iex
```

### Linux / macOS (Bash)
```bash
curl -fsSL https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.230.0/install.sh | bash
```

## Changelog v6.230.0

- Gracefully handle gitmap cd and lookup commands when targets are not found without raw stack traces or duplicate stderr output
- Add fuzzy repository and workdir suggestion engine (Did you mean: <repo>?) on cd lookup misses
- Standardize ErrorTypeNotFound across apperror, cliexit, and cmd packages
