## Quick Install v6.183.0

### Windows (PowerShell 5.1+)
```powershell
irm https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.183.0/install.ps1 | iex
```

### Linux / macOS (Bash)
```bash
curl -fsSL https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.183.0/install.sh | bash
```

## Changelog v6.183.0

- Overhauled GitMap CLI footer metadata output with dedicated Short and Long variants
- Implemented short footer displaying Version and Git Commit SHA, integrated across all help invocations (`--help`, `-h`, `gitmap help <topic>`)
- Implemented long footer displaying complete binary identity block: Name, Git URL (`https://github.com/alimtvnetwork/gitmap-v28`), Version, Commit SHA, Database path, and Installed binary path
- Added dynamic binary and repository metadata resolution with cascading fallbacks (ldflags -> git config -> origin URL -> `version.json`), guaranteeing zero empty fields
- Flattened nested conditionals across `04-code/` streamwriter reference implementation, achieving 0 nested `if` violations across 3,089 repository files
- Bumped version to `v6.183.0` across all Single Source of Truth manifests (`version.json`, `package.json`, `gitmap/constants/constants.go`, `changelog.md`)
- Verified 100% green pass across all 16 CI/CD quality gates locally and synchronized release assets
