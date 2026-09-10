## Quick Install v6.208.0

### Windows (PowerShell 5.1+)
```powershell
irm https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.208.0/install.ps1 | iex
```

### Linux / macOS (Bash)
```bash
curl -fsSL https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.208.0/install.sh | bash
```

## Changelog v6.208.0

- Connect repodb to dbengine ORM with typed repositories and row scanners
- Upgrade 30-db-struct-enum-generator.py to support db tags, entity models, and typed mutations
- Normalize SQLite database primary keys to PascalCase <Entity>Id across repodb and store
- Eradicate swallowed errors across pipelinedb, dbengine, and repodb with universal AppError wrapping
- Implement in-memory batch timestamp caching, 8KB binary sniffer, and directory pruning in indexer
- Fix Linux corrupted install directories with safe recovery and cross-OS Darwin/Windows support
- Decouple Google Antigravity Desktop IDE installer from agy CLI utility
