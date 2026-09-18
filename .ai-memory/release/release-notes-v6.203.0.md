## Quick Install v6.203.0

### Windows (PowerShell 5.1+)
```powershell
irm https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.203.0/install.ps1 | iex
```

### Linux / macOS (Bash)
```bash
curl -fsSL https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.203.0/install.sh | bash
```

## Changelog v6.203.0

- Add qBittorrent and uTorrent cross-platform installers across Windows (choco, winget), Ubuntu/Debian (apt), and macOS (brew)
- Implement portable JSON configuration export/import engine (export-config, import-config, improt-config) for VS Code, qBittorrent, and uTorrent
- Support batch export and import across folders with automatic OS path and file format translation (.ini <-> .conf)
- Add dedicated help documentation (export-config.md, import-config.md) and commands.ts UI metadata
