## Quick Install v6.298.0

### Windows (PowerShell 5.1+)
```powershell
irm https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.298.0/install.ps1 | iex
```

### Linux / macOS (Bash)
```bash
curl -fsSL https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.298.0/install.sh | bash
```

## Changelog v6.298.0

- MultiClone Subsystem: Added `gitmap multiclone` (aliases `mc`, `mutliclone`) to batch clone repositories from markdown codeblocks, raw lists, or stdin with shorthand expansion, description stripping, and deduplication.
- AI Split DB Command Tracking: Integrated global `--ai` execution history logging into `~/.gitmap/ai/instructions.db` with duration and exit code metrics.
- Frequent AI Commands & Clipboard Export: Implemented `gitmap ai ls` and `gitmap ai history` to inspect top frequent commands with `--copy` to OS clipboard.
- Native Automation Search Benchmark: Verified Go native search is 46.1x faster (1.5ms vs 70.5ms) than cached Python grep with 0 disk bloat.
- Rich Terminal Help Menus: Implemented two-column styled help screens for `gitmap sync` and `gitmap github` (`gd`) with formatted error diagnostics.
- Anti-Gravity Settings Protection: Protected all Antigravity IDE configuration files and settings across Windows and Unix OS clean operations.
- WinUtil & LinUtil Integration Architecture: Authored comprehensive plan for native Go OS auto-login, Windows tweaks, and Linux system cleanup.
