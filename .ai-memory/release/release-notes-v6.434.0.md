## Quick Install v6.434.0

### Windows (PowerShell)

```powershell
Invoke-WebRequest -Uri https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.434.0/install.ps1 -OutFile install.ps1; .\install.ps1 -TargetDir ".ai-memory/prompts" -Version "v6.434.0"
```

### Unix / Linux / macOS (Bash)

```bash
curl -sL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.434.0/install.sh | bash -s -- ".ai-memory/prompts" "v6.434.0"
```

---

## What's Changed in v6.434.0

### Added
- RCA-55: Fix test auto-dest-repo workspace pollution in VS Code and Antigravity, correct w4 os, distinguish offline vs auth_failed in nodes clone
