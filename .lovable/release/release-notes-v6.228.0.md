## Quick Install v6.228.0

### Windows (PowerShell 5.1+)
```powershell
irm https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.228.0/install.ps1 | iex
```

### Linux / macOS (Bash)
```bash
curl -fsSL https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.228.0/install.sh | bash
```

## Changelog v6.228.0

- Add Linux archive package installation (tar, gz, zip) with intelligent strategy detection
- Support automated strategy detection for ELF binaries, install scripts, and source builds
- Integrate universal uninstaller for archive deployments with desktop database refresh
- Add cross-platform compilation stubs and enforce strict boolean, enum, and lint guidelines
