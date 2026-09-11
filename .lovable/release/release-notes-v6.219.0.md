## Quick Install v6.219.0

### Windows (PowerShell 5.1+)
```powershell
irm https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.219.0/install.ps1 | iex
```

### Linux / macOS (Bash)
```bash
curl -fsSL https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.219.0/install.sh | bash
```

## Changelog v6.219.0

- Relocate CI/CD runner temporary telemetry and session runs from OS temp directory to repository-internal .lovable/cicd
- Ensure all failure banner stream paths and artifact locations are displayed as repository-relative paths
- Add .lovable/cicd/ to .gitignore to avoid untracked working directory clutter
