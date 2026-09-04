## Quick Install v6.185.0

### Windows (PowerShell 5.1+)
```powershell
irm https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.185.0/install.ps1 | iex
```

### Linux / macOS (Bash)
```bash
curl -fsSL https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.185.0/install.sh | bash
```

## Changelog v6.185.0

- Overhauled .github/scripts/e2e-cli-smoke.py with async worker group execution running parallel tests concurrently via asyncio
- Added --all (-a) CLI flag to e2e-cli-smoke.py to control verbose pass/fail printing vs concise summary
- Optimized default smoke test output to only report failed tests or a single clean pass confirmation line
- Structured stateful CLI commands (schedule and macro chains) into isolated sequential worker tasks executing in parallel with independent tests
- Bumped version to v6.185.0 across all Single Source of Truth manifests
