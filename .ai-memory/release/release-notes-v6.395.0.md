## Quick Install v6.395.0

### Windows (PowerShell)

```powershell
Invoke-WebRequest -Uri https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.395.0/install.ps1 -OutFile install.ps1; .\install.ps1 -TargetDir ".ai-memory/prompts" -Version "v6.395.0"
```

### Unix / Linux / macOS (Bash)

```bash
curl -sL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.395.0/install.sh | bash -s -- ".ai-memory/prompts" "v6.395.0"
```

---

## What's Changed in v6.395.0

### Added
- resolve ee external repo targets, fix cfr bare repo suggestions, and add cache cleanup guard
