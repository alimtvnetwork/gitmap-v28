## Quick Install v6.265.0

### Windows (PowerShell 5.1+)
```powershell
irm https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.265.0/install.ps1 | iex
```

### Linux / macOS (Bash)
```bash
curl -fsSL https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.265.0/install.sh | bash
```

## Changelog v6.265.0

- Fix SQL double precision keywords (DOUBLE and FLOAT) in cli/cmdautomation/db_generate.go to prevent misspell false positive
- Replace deprecated strings.Title with custom formatDomainTitle rune helper in cli/cmdautomation/plan_consolidate.go
- Rename cachedErr to errCached in cli/cmdai/ai_python_detector.go to satisfy Go errname linter convention
- Remove unused option getter functions in cli/cmdautomation (version_sync_cmd.go, release_bump_cmd.go, milestones_cmd.go) and unused aiCmd in cli/cmdai/ai_cmd.go
- Resolve inverted success check (!isSuccess) in cli/cmdautomation/phase2_test.go by renaming to isGenerated
- Sanitize repository path references in RCA documents to maintain 100% relative path compliance
