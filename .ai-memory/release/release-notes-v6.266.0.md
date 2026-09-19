## Quick Install v6.266.0

### Windows (PowerShell 5.1+)
```powershell
irm https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.266.0/install.ps1 | iex
```

### Linux / macOS (Bash)
```bash
curl -fsSL https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.266.0/install.sh | bash
```

## Antigravity Manager Installation & Update

### Windows (PowerShell)
```powershell
irm https://raw.githubusercontent.com/alimtvnetwork/Antigravity-Manager/main/install.ps1 | iex
```

### UNIX / Linux / macOS (Bash)
```bash
curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/Antigravity-Manager/main/install.sh | bash
```

## Changelog v6.266.0

- Update Antigravity Manager (`agm`, `ag-manager`, `antigravity-manager`) installer and updater to canonical one-liners: PowerShell script execution for Windows and curl | bash execution for UNIX
- Add `gitmap agm update` subcommand with aliases (`up`, `u`) and support `--dry-run`
- Route `gitmap update agm` (and aliases `ag-manager`, `antigravity-manager`) in root utility dispatcher directly to Antigravity Manager updater
- Replace legacy `lbjlaq/Antigravity-Manager` repository URLs with canonical `alimtvnetwork/Antigravity-Manager`
- Remove legacy binary installer helpers and unused executable copier functions from `cli/cmdinstall`
- Fix `03-ai-scripts/09-cli-help-auditor.py` concurrency argument passing
