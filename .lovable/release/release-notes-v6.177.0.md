## Quick Install v6.177.0

### Windows (PowerShell 5.1+)
```powershell
irm https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.177.0/install.ps1 | iex
```

### Linux / macOS (Bash)
```bash
curl -fsSL https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.177.0/install.sh | bash
```

## Changelog v6.177.0

- Fixed pull and status table line wrapping and column misalignment in terminal output
- Implemented smart branch prefix omission for common branch prefixes (feature/, feat/, release/, bugfix/, hotfix/, fix/, dependabot/)
- Added middle-truncation engine with ellipsis (...) preserving starting characters and ending 5 characters
- Enforced fixed column caps across repository name, branch, latest branch, PR status, pull status, and SHA
- Corrected ANSI escape sequence padding calculation for styled Lipgloss elements to ensure exact column alignment
- Added comprehensive unit tests in cmd/pull_table_test.go covering prefix omission, middle-truncation, and ANSI padding
