## Quick Install v6.201.0

### Windows (PowerShell 5.1+)
```powershell
irm https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.201.0/install.ps1 | iex
```

### Linux / macOS (Bash)
```bash
curl -fsSL https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.201.0/install.sh | bash
```

## Changelog v6.201.0

- Add automatic Git remote tag release discovery (https://github.com/lbjlaq/Antigravity-Manager.git) for Antigravity Manager GUI with semver sorting and multi-OS asset matching
- Add specific version targeting --version <ver> to gitmap install ag-manager and gitmap agy install manager
- Enable full CLI flag parity (--dry-run, --yes/-y, --verbose/-v, --version) in gitmap agy install [manager|cli|all]
- Fix gitmap in alias routing and display grouped tool catalog, descriptions, and profiles when run without arguments
- Add and document antigravity, ag-manager, ag-ctx, build-essential, and installation profiles (dev, ubuntu, ubuntu-dev-ai, ai, backend, fullstack, minimal)
