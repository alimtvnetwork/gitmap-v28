## Quick Install v6.438.2

### Windows (PowerShell)

```powershell
Invoke-WebRequest -Uri https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.438.2/install.ps1 -OutFile install.ps1; .\install.ps1 -TargetDir ".ai-memory/prompts" -Version "v6.438.2"
```

### Unix / Linux / macOS (Bash)

```bash
curl -sL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.438.2/install.sh | bash -s -- ".ai-memory/prompts" "v6.438.2"
```

---

## What's Changed in v6.438.2

### Added
- Fix helptext examples formatting for fix-credential
