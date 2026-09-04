## Quick Install v6.179.0

### Windows (PowerShell 5.1+)
```powershell
irm https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.179.0/install.ps1 | iex
```

### Linux / macOS (Bash)
```bash
curl -fsSL https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.179.0/install.sh | bash
```

## Changelog v6.179.0

- Added self-healing Git diagnostic command engine: `gitmap fix-git` (aliases: `gitmap fg`, `gitmap --fix-git`)
- Added automatic permission repair and Windows NTFS ACL Full Control assignment on `.git` directories (`icacls` / `chmod -R u+rwX`)
- Added elevated PowerShell self-healing script `scripts/fix-all-permissions.ps1` for system-wide ACL and ownership recovery
- Added stale lockfile detection and removal (`.git/index.lock`, `HEAD.lock`, `config.lock`)
- Added index corruption detection and auto-recovery from `HEAD` via `git reset` with timestamped backup
- Added detection and guided resolution for unmerged files and unresolved merge conflicts
- Added conflict-safe pull protection by backing up colliding untracked files to `.git/gitmap-backup/`
- Added `gitmap workdir default [path]` to inspect or configure the active workspace default directory
- Added `gitmap workdir path` to print raw absolute default workdir path for shell navigation scripts
- Resolved `gitmap cd work` and `gitmap cd default` to seamlessly navigate to the default work directory
- Suppressed internal Go runtime stack traces and line numbers on validation errors for clean user-facing error reporting
- Added interactive step builder fallback when command steps are omitted in `gitmap macro add <name>`
- Authored formal issue specification and 5-part RCA in `spec/22-app-issues/35-reconcile-prompt-nested-if-ci-failure.md`
- Flattened nested conditionals in `cmd/reconcile_prompt.go` to maintain 0 nested if violations across the codebase
- Verified 100% green pass across all 25 quality gates in `.lovable/ai-fix-scripts/06-cicd-local-runner.py`
