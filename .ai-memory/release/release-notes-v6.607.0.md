# GitMap v6.607.0

## What's Changed in v6.607.0

- **Incoming Pulled Changes Diffstat**: Automatically renders the detailed `git diff --stat` of all incoming commits integrated by `gitmap pull`, listing modified files and line deltas.
- **Git Status on Dirty Working Tree**: Renders standard Git short status (`git status --short`) with ANSI color-coding for dirty working trees, displaying staged, modified, deleted, and untracked files with full fidelity.
- **Repository Path First Tree Header**: Formats the initial command execution header with the exact filesystem repository path (`-> <repoPath>`), eliminating redundant command echoing.

### Quick Install One-Liners

**Windows (PowerShell):**
```powershell
irm https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.607.0/install.ps1 | iex
```

**Linux / macOS (Bash):**
```bash
curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.607.0/install.sh | sh
```
