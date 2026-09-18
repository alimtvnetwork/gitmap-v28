## Quick Install v6.225.0

### Windows (PowerShell 5.1+)
```powershell
irm https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.225.0/install.ps1 | iex
```

### Linux / macOS (Bash)
```bash
curl -fsSL https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.225.0/install.sh | bash
```

## Changelog v6.225.0

- Integrated robust_rmtree with Windows read-only file attribute unlinking across all build and test cleanup hooks
- Guaranteed zero-byte temporary storage footprint in OS temp gitmap after test runs
- Added --allow-serial-runners across all golangci-lint invocations to prevent parallel lock contention
