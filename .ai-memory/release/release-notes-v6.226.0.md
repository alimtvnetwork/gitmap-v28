## Quick Install v6.226.0

### Windows (PowerShell 5.1+)
```powershell
irm https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.226.0/install.ps1 | iex
```

### Linux / macOS (Bash)
```bash
curl -fsSL https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.226.0/install.sh | bash
```

## Changelog v6.226.0

- Fixed unbalanced if statement syntax error in E2E GitHub Actions workflow
- Guarded TypeScript AST imports in check-nested-ifs against missing node_modules
- Auto-migrated and ensured ssh_hosts and ssh_history tables exist across cmdssh subcommands
- Resolved relative path and nested-if lint violations
