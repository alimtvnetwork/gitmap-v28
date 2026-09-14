## Quick Install v6.236.0

### Windows (PowerShell 5.1+)
```powershell
irm https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.236.0/install.ps1 | iex
```

### Linux / macOS (Bash)
```bash
curl -fsSL https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.236.0/install.sh | bash
```

## Changelog v6.236.0

- First-class user@ip SSH join syntax (gitmap ssh-join user@ip [alias]), default host-<ip> alias generation, and automatic IP login user resolution
- Subnet discovery scanner (gitmap sj scan) probing port 22 and cross-referencing registered hosts
- Machine health and latency ping (gitmap sj status) with online/offline diagnostics
- Redesigned Git pull UI with 80ms active background ticker, animated Braille spinner, accurate Windows TTY detection, 4-step milestones, and worker slot concurrency
- Scripts-fixer alignment: terminal profile utilities (jq, yq, zellij), OS-aware dev/small-dev partitioning, git-compact integration, and AI profile trees
- Comprehensive leaf command help, mistake recovery guidance, and documentation parity across terminal UI, catalogs, and markdown docs
