## Quick Install v6.204.3

### Windows (PowerShell 5.1+)
```powershell
irm https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.204.3/install.ps1 | iex
```

### Linux / macOS (Bash)
```bash
curl -fsSL https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.204.3/install.sh | bash
```

## Changelog v6.204.3

- Fix Ubuntu ZSH update prompt and prevent unwanted apt reinstall during update
- Prevent sudo password requests and preserve existing Oh-My-Zsh configurations
- Auto-detect existing ZSH binary and skip non-interactive setup runs
- Enforce LF line endings and automatic gofmt synchronization during release orchestration
