## Quick Install v6.186.0

### Windows (PowerShell 5.1+)
```powershell
irm https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.186.0/install.ps1 | iex
```

### Linux / macOS (Bash)
```bash
curl -fsSL https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.186.0/install.sh | bash
```

## Changelog v6.186.0

- Added smart Chrome profile import supporting directory scanning, globbing (*.json), email lookup, and auto-detecting current directory
- Implemented safe non-destructive Chrome profile import: matches existing profiles by email and automatically creates new profile directories without breaking existing profiles
- Added step-by-step progress logging across all import phases (inspect, resolve, restore, stage extensions, and Local State registration)
- Added gitmap chrome profile import-check (inspect) command to preview snapshot metadata and planned import actions before execution
- Enhanced gitmap chrome profile ls to display account emails alongside discovered snapshot files in the active directory
- Added --except / --exclude flag to skip profile IDs, slugs, names, emails, or prefix patterns during import
- Bumped version to v6.186.0 across all Single Source of Truth manifests
