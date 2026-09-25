## Quick Install v6.343.0

### Windows (PowerShell 5.1+)
```powershell
irm https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.343.0/install.ps1 | iex
```

### Linux / macOS (Bash)
```bash
curl -fsSL https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.343.0/install.sh | bash
```

## Changelog v6.343.0

- Git Pull Efficient (PAE): Evaluates repository activity strictly on actual git log / commit trace changes within 24h window (skipping quiescent repositories), saving complete commit traces into SQLite gitmap-pull.db (PullRepoRun.Notes).
- Sub-Millisecond Commit Trace Search: Added SearchPullTraces with indexes on LastCommitSha, RepoPath, HasChanges, and CreatedAt for instant querying across all historical repository commit logs without invoking git subprocesses.
- Antigravity IDE Injection Hardening: Enhanced gitmap pe agy fix with automatic fallback to default gitmap repo root and conversation resolution via resolveConversationForDispatch when outside git repos.
- CI/CD & Exhaustive Switch Fixes: Resolved exhaustive linter checks on PullStepType in cmdpull and ScriptTemplateType in cmdai, keeping 100% CI compliance.
