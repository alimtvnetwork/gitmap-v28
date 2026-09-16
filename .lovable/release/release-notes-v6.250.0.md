## Quick Install v6.250.0

### Windows (PowerShell 5.1+)
```powershell
irm https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.250.0/install.ps1 | iex
```

### Linux / macOS (Bash)
```bash
curl -fsSL https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.250.0/install.sh | bash
```

## Changelog v6.250.0

- Anchor pipeline database to CLI binary data directory in dedicated pipeline/ folder with per-repo slug isolation
- Resolve pipeline DB path dynamically and eliminate hardcoded repo-root fallback
- Display canonical CLI pipeline database path in all summaries, error logs, and root identity footers
- Decompose pipeline split DB connection and schema management adhering to strict 100-line coding guidelines
