## Quick Install v6.400.0

### Windows (PowerShell)

```powershell
Invoke-WebRequest -Uri https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.400.0/install.ps1 -OutFile install.ps1; .\install.ps1 -TargetDir ".ai-memory/prompts" -Version "v6.400.0"
```

### Unix / Linux / macOS (Bash)

```bash
curl -sL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.400.0/install.sh | bash -s -- ".ai-memory/prompts" "v6.400.0"
```

---

## What's Changed in v6.400.0

### Added
- restructure markdown help into modern box display format and universal catalog fallback
