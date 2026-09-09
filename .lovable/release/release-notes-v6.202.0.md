## Quick Install v6.202.0

### Windows (PowerShell 5.1+)
```powershell
irm https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.202.0/install.ps1 | iex
```

### Linux / macOS (Bash)
```bash
curl -fsSL https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.202.0/install.sh | bash
```

## Changelog v6.202.0

- Remove dead runInstallAgManager and runInstallAntigravity functions to pass strict unused and lint-baseline-diff CI checks
- Upgrade .github/scripts/full-suite-lint.sh to native high-performance Bash runner with live tee streaming and SIMD grep -cE issue counting
- Modernize .github/scripts/full-suite-lint.py with line-buffered subprocess.Popen real-time stdout streaming and O(1) memory tracking
