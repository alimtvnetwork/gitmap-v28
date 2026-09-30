## Quick Install v6.415.0

### Windows (PowerShell)

```powershell
Invoke-WebRequest -Uri https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.415.0/install.ps1 -OutFile install.ps1; .\install.ps1 -TargetDir ".ai-memory/prompts" -Version "v6.415.0"
```

### Unix / Linux / macOS (Bash)

```bash
curl -sL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.415.0/install.sh | bash -s -- ".ai-memory/prompts" "v6.415.0"
```

---

## What's Changed in v6.415.0

### Added
- JSON envelope v2 schema with enriched metadata attributes, working directory telemetry, and top-level variables dictionary
- Multi-JSON merge command (`gitmap merge-json` / `gitmap json merge`) with item deduplication and 1-based indexing
- Full phrasing variants for fleet key deployment: `gitmap deploy all-keys`, `gitmap deploy all keys`, `gitmap deploy keys all`, `gitmap deploy key all`, and `gitmap deploy keys`
- Screen buffer clearing integrated into `gitmap terminal clear` and `gitmap clear terminal`

### Changed
- Bare `gitmap import` without arguments or `--confirm` now displays an interactive guidance catalog instead of throwing E9000 stack traces
- Typos and unknown root commands now print concise error notifications and recommendations instead of dumping the multi-page help catalog
- Removed all underscores from SSH node JSON representations (camelCase) and updated repo-secrets manifests with 1-based machine indexing
