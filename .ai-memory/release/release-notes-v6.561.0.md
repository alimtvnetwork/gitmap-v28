# GitMap v6.561.0

## What's Changed in v6.561.0

- **Histogram Diff Automerge Optimization:** Passes `-X diff-algorithm=histogram` during automerge attempts, eliminating spurious merge conflicts.
- **Git Rerere Auto-Memory Engine:** Auto-activates Git `rerere` to record and re-apply past conflict resolutions.
- **Structured Conflicted File Extractors:** Parses conflict stdout and git status porcelain to cleanly extract all unmerged files.
- **Multi-IDE Launcher & Guidance:** Detects VS Code, Cursor, Windsurf, or Codium and provides step-by-step 3-Way Merge Editor walkthroughs.
- **GitMap Conflict Resolution CLI:** First-class `gitmap resolve <repo> [--theirs | --ours | --ide | --abort]` and `gitmap conflicts [repo]`.
- **Extended `gitmap fix <repo>` Actions:** Added `theirs`, `ours`, `resolve`, `ide`, `code`, `abort`, and `rerere` actions.
- **Post-Pull Merge Conflict UI Card:** Styled Lipgloss conflict card during pull failures with one-key resolution actions.
- **CLI & Web UI Documentation:** Added `helpdoc/resolve.md`, `helpdoc/conflicts.md`, and dashboard troubleshooting registry.

### Quick Install One-Liners

**Windows (PowerShell):**
```powershell
irm https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.561.0/install.ps1 | iex
```

**Linux / macOS (Bash):**
```bash
curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.561.0/install.sh | sh
```
