## Quick Install v6.222.0

### Windows (PowerShell 5.1+)
```powershell
irm https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.222.0/install.ps1 | iex
```

### Linux / macOS (Bash)
```bash
curl -fsSL https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.222.0/install.sh | bash
```

## Changelog v6.222.0

- Route 'gitmap agy install' to Antigravity Desktop IDE and add 'gitmap agm install' for Antigravity Tools/Manager
- Add dynamic CPU freeness and memory detection, scaling test worker pools to 32 threads for full CPU saturation
- Resolve multi-module go.mod paths and pass 100% of all 46 CI/CD quality gates
