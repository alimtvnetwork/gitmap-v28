## Quick Install v6.214.0

### Windows (PowerShell 5.1+)
```powershell
irm https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.214.0/install.ps1 | iex
```

### Linux / macOS (Bash)
```bash
curl -fsSL https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.214.0/install.sh | bash
```

## Changelog v6.214.0

- Fix gosec G115 integer overflow conversion uint64 -> int64 in pipeline_split_ops.go
- Ensure all cross-platform pattern normalizations and quality gates pass
