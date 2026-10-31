# GitMap v6.635.0

## What's Changed in v6.635.0

- **Terminal Help Parity (`thelp` / `th`):** Aliased `thelp` and `th` to the core help dispatcher so users can query terminal help intuitively (`gitmap thelp <command>`).
- **Merge & Transfer Suite Menu (`thelp merge` / `gitmap merge`):** Added structured 2-column rich help menu displaying one-way file merges (`merge-right`, `merge-left`), two-way sync (`merge-both`), chronological commit replay (`merge-in` / `commit-in`), and pull requests (`merge write`).
- **Merge Subcommand Routing:** Supported `gitmap merge right`, `gitmap merge left`, `gitmap merge both`, and `gitmap merge in` natively with automatic help menu on bare invocation.
- **First-Class `merge-in` Alias:** Added `merge-in` as a first-class alias for `commit-in` (`cin`) for multi-source history stitching.
- **AI Automation Terminal Documentation:** Added complete documentation and executable examples for `login set ai-auto <agent>`, `login set global ai-auto`, `settings set ai_auto.*`, and `--ai-agent <name>` flags.
- **Curated Help Examples Registry:** Registered canonical examples for `merge`, `merge-right`, `merge-left`, `merge-in`, `commit-in`, `commit-both`, `commit-left`, `commit-right`, `thelp`, `specific`, `login`, and `settings`.
- **Query Specific Guidance (`thelp specific`):** Added dedicated help guide for discovering command-specific terminal help screens.

### Quick Install One-Liners

**Windows (PowerShell):**
```powershell
irm https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.635.0/install.ps1 | iex
```

**Linux / macOS (Bash):**
```bash
curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.635.0/install.sh | sh
```
