## Quick Install v6.363.0

### Windows (PowerShell)

```powershell
Invoke-WebRequest -Uri https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.363.0/install.ps1 -OutFile install.ps1; .\install.ps1 -TargetDir ".ai-memory/prompts" -Version "v6.363.0"
```

### Unix / Linux / macOS (Bash)

```bash
curl -sL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.363.0/install.sh | bash -s -- ".ai-memory/prompts" "v6.363.0"
```

---

## What's Changed in v6.363.0

### Added
- Fix telemetry array extraction in fleet update, verified running prompts restore and deploy
