## Quick Install v6.436.0

### Windows (PowerShell)

```powershell
Invoke-WebRequest -Uri https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.436.0/install.ps1 -OutFile install.ps1; .\install.ps1 -TargetDir ".ai-memory/prompts" -Version "v6.436.0"
```

### Unix / Linux / macOS (Bash)

```bash
curl -sL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.436.0/install.sh | bash -s -- ".ai-memory/prompts" "v6.436.0"
```

---

## What's Changed in v6.436.0

### Added
- RCA-57: Align Fleet Clone Table, Suppress Stderr Escape Leaks, and Harden W3 Liveness Resilience
