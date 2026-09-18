## Quick Install v6.252.0

### Windows (PowerShell 5.1+)
```powershell
irm https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.252.0/install.ps1 | iex
```

### Linux / macOS (Bash)
```bash
curl -fsSL https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.252.0/install.sh | bash
```

## Changelog v6.252.0

- Capture and display Go stack trace on pipeline command failures
- Diagnose workflow initialization and syntax errors when GitHub Actions log archive is missing
- Eliminate blank job labels and misleading synthetic step errors on failed gh runs
