## Quick Install v6.176.0

### Windows (PowerShell 5.1+)
```powershell
irm https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.176.0/install.ps1 | iex
```

### Linux / macOS (Bash)
```bash
curl -fsSL https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.176.0/install.sh | bash
```

## Changelog v6.176.0

- Fixed CI/CD compatibility: formatted cmd/agy_pin_projects.go to strict gofmt specifications
- Purged orphaned submodule gitlinks ensuring clean actions/checkout across GitHub Actions workflows
- Standardized US English spelling across all documentation in spec/ and de-literalized test lookup tables
- Gracefully handled missing VS Code user-data root in headless CI runners for gitmap vscode ls
- Verified 115/115 E2E installer smoke tests with 100% green verification on release pipeline
- Added interactive command builder and piped input support for gitmap macro add <name> without stack traces on validation errors
- Added gitmap agy pin-projects command suite (ls, add, rm, --json, --all) to pin, list, and unpin Antigravity projects
- Added --pinned / -p filter flag to gitmap agy ls to quickly inspect pinned projects
- Added first-class dynamic root-level execution for saved macros allowing gitmap <macro-name> directly
- Added root-level macro utility command aliases (macro-list, macro-add, macro-run, macro-record, macro-show, macro-rm)
- Enforced vertical newline styling rules (R13-R16) with blank lines before if, after closing braces, and before return
- Passed all 16 local CI/CD quality gates across 3 sequential batches with 100% green verification (exit code 0)
