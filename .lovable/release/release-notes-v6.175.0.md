## Quick Install v6.175.0

### Windows (PowerShell 5.1+)
```powershell
irm https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.175.0/install.ps1 | iex
```

### Linux / macOS (Bash)
```bash
curl -fsSL https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.175.0/install.sh | bash
```

## Changelog v6.175.0

- Added gitmap agy pin-projects command suite (ls, add, rm, --json, --all) to pin, list, and unpin Antigravity projects
- Added --pinned / -p filter flag to gitmap agy ls to quickly inspect pinned projects
- Added first-class dynamic root-level execution for saved macros allowing gitmap <macro-name> directly
- Added root-level macro utility command aliases (macro-list, macro-add, macro-run, macro-record, macro-show, macro-rm)
- Added comprehensive documentation in docs/commands/agy/pin-projects.md and docs/commands/macro.md
- Formatted function signatures, parameter declarations, and invocations across Go codebase to Rule 9a/9b multi-line standards
- Standardized value-based parameter structs (*Params) and eliminated bare void functions across Go domain and service layers
- Enforced universal Result[T] envelope in Go and TypeScript with .IsSuccess(), .IsFailed(), .HasError(), .HasNoError(), and .HasValidError()
- Enforced code hygiene, Unix LF line endings across all files, UTF-8 (no BOM) encoding, and Markdown heading spacing (MD022/MD032)
- Enforced vertical newline styling rules (R13-R16) with blank lines before if, after closing braces, and before return
- Passed all 16 local CI/CD quality gates across 3 sequential batches with 100% green verification (exit code 0)
