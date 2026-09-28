## Quick Install v6.380.0

### Windows (PowerShell)

```powershell
Invoke-WebRequest -Uri https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.380.0/install.ps1 -OutFile install.ps1; .\install.ps1 -TargetDir ".ai-memory/prompts" -Version "v6.380.0"
```

### Unix / Linux / macOS (Bash)

```bash
curl -sL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.380.0/install.sh | bash -s -- ".ai-memory/prompts" "v6.380.0"
```

---

## What's Changed in v6.380.0

### Added
- feat(special-repos): introduce special default repositories `repo-secrets` (`rs`) and `repo-cache` (`rc` / `repo-storage`) with SQLite split-db persistence (`gitmap-special-repos.db`), first-scan one-time discovery prompt, `gitmap cd rs|rc` shortcuts, auto-sequencing (`XX-<repo>/01-<slug>.ext`), and automatic git add/commit/push
- feat(settings,ui): expose `special_repos.secrets_name` and `special_repos.cache_name` in `gitmap settings` CLI and interactive Web UI with rich help documentation
- docs(prompts,guidelines): integrate Special Repositories V2 prompt (`01-prompts/special-repos-secrets-and-cache.md`) and enforce R19 guideline in `coding-guidelines`
