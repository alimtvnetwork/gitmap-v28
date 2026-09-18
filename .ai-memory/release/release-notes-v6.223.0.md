## Quick Install v6.223.0

### Windows (PowerShell 5.1+)
```powershell
irm https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.223.0/install.ps1 | iex
```

### Linux / macOS (Bash)
```bash
curl -fsSL https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.223.0/install.sh | bash
```

## Changelog v6.223.0

- Unify slow and fast unit tests into a single 32-worker priority thread pool, eliminating thread starvation
- Bypass external GitHub API dials and subprocess probing in unit tests, speeding up cmd tests by up to 110x
- Accelerate Go Test Coverage Profile from 193s down to 2.5s while meeting 100% of coverage floor gates
