## Quick Install v6.242.0

### Windows (PowerShell 5.1+)
```powershell
irm https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.242.0/install.ps1 | iex
```

### Linux / macOS (Bash)
```bash
curl -fsSL https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.242.0/install.sh | bash
```

## Changelog v6.242.0

- Fix test DDL for ssh_hosts port column
- Resolve gosec G115 integer conversions in storage display and POSIX metrics
- Safeguard nil context and avoid test deadlock in cmdssh
