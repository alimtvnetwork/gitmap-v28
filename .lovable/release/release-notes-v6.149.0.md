## Quick Install v6.149.0

### Windows (PowerShell 5.1+)

```powershell
irm https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.149.0/install.ps1 | iex
```

### Linux / macOS (Bash)

```bash
curl -fsSL https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.149.0/install.sh | bash
```

## Changelog v6.149.0

- Autonomously flatten all nested if statements across codebase with guard clauses and early returns
- Enforce implicit boolean evaluations and positive prefix standards across all packages
- Eliminate explicit boolean comparisons (== true, == false) repository-wide
- Fix gocritic appendAssign finding in CI diff gate across cmd and tests
- Connect check-enum-and-boolean.py to CI/CD local runner suite with 100% green verification
