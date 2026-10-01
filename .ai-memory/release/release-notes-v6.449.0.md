## Quick Install v6.449.0

### Windows (PowerShell 5.1+)

```powershell
irm https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.449.0/install.ps1 | iex
```

### Linux / macOS (Bash)

```bash
curl -fsSL https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.449.0/install.sh | bash
```

### Script / Prompt Sync (PowerShell)

```powershell
Invoke-WebRequest -Uri https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.449.0/install.ps1 -OutFile install.ps1; .\install.ps1 -TargetDir ".ai-memory/prompts" -Version "v6.449.0"
```

### Script / Prompt Sync (Bash)

```bash
curl -sL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.449.0/install.sh | bash -s -- ".ai-memory/prompts" "v6.449.0"
```

---

## What's Changed in v6.449.0

### Fixed
- Resolve nested ifs, swallowed db errors, and runner path attributes (RCA-095)
