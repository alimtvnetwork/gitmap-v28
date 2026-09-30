## Quick Install v6.431.0

### Windows (PowerShell)

```powershell
Invoke-WebRequest -Uri https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.431.0/install.ps1 -OutFile install.ps1; .\install.ps1 -TargetDir ".ai-memory/prompts" -Version "v6.431.0"
```

### Unix / Linux / macOS (Bash)

```bash
curl -sL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.431.0/install.sh | bash -s -- ".ai-memory/prompts" "v6.431.0"
```

---

## What's Changed in v6.431.0

### Added
- Fix cross-platform tests and boolean linter nested if policy
