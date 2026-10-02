## Quick Install v6.438.3

### Windows (PowerShell)

```powershell
Invoke-WebRequest -Uri https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.438.3/install.ps1 -OutFile install.ps1; .\install.ps1 -TargetDir ".ai-memory/prompts" -Version "v6.438.3"
```

### Unix / Linux / macOS (Bash)

```bash
curl -sL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.438.3/install.sh | bash -s -- ".ai-memory/prompts" "v6.438.3"
```

---

## What's Changed in v6.438.3

### Added
- Flatten nested ifs, switch applyOption, and handle cache db scan errors
