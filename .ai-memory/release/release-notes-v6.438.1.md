## Quick Install v6.438.1

### Windows (PowerShell)

```powershell
Invoke-WebRequest -Uri https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.438.1/install.ps1 -OutFile install.ps1; .\install.ps1 -TargetDir ".ai-memory/prompts" -Version "v6.438.1"
```

### Unix / Linux / macOS (Bash)

```bash
curl -sL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.438.1/install.sh | bash -s -- ".ai-memory/prompts" "v6.438.1"
```

---

## What's Changed in v6.438.1

### Added
- Fix nested-if in cpar, sync generated commands and register fix-credential test constants
