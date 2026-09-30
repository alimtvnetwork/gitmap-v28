## Quick Install v6.419.0

### Windows (PowerShell)

```powershell
Invoke-WebRequest -Uri https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.419.0/install.ps1 -OutFile install.ps1; .\install.ps1 -TargetDir ".ai-memory/prompts" -Version "v6.419.0"
```

### Unix / Linux / macOS (Bash)

```bash
curl -sL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.419.0/install.sh | bash -s -- ".ai-memory/prompts" "v6.419.0"
```

---

## What's Changed in v6.419.0

### Added
- JSON envelope multi-pass variable interpolation (`${keyPath}`, `${adminUser}`, `${workDir}`) and structured `workDirectory.variables` support
- Single-argument current-user OS password change (`gitmap os change-password <password>`) with confirmation and dual-argument admin mode (`gitmap os change-password <user> <password>`)
- Interactive salted-RSA SSH join vault password prompt samples and documentation
- Full repo-secrets sequence, lowercase naming, and lowerCamelCase envelope harmonization across `01-gitmap` through `07-final-network-machine`
- PowerShell completion and predictive suggestions for `nodes`, `node`, `ping`, `os`, and enriched `ssh` subcommands with automatic installer session activation
