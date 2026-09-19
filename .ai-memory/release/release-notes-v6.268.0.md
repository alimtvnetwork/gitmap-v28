## Quick Install v6.268.0

### Windows (PowerShell 5.1+)
```powershell
irm https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.268.0/install.ps1 | iex
```

### Linux / macOS (Bash)
```bash
curl -fsSL https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.268.0/install.sh | bash
```

## Changelog v6.268.0

- Add `gitmap ssh nodes` command to list registered SSH nodes directly instead of failing with alias not found
- Add `gitmap ssh ls` command to list registered SSH nodes
- Improve SSH nodes and ls terminal view with rich ANSI colors (Cyan headers, Yellow roles for control-plane/master, Green bullet status, White aliases and host:port)
- Add top and bottom newline padding and precise column alignment (100-char table layout) for terminal readability
- Enhance alias-not-found suggestion view with the same rich colored registered hosts table and recall guidance
- Add `gitmap ssh-join nodes` and `gitmap sj nodes` aliases for nodes listing
