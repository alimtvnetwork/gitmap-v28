## Quick Install v6.182.0

### Windows (PowerShell 5.1+)
```powershell
irm https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.182.0/install.ps1 | iex
```

### Linux / macOS (Bash)
```bash
curl -fsSL https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.182.0/install.sh | bash
```

## Changelog v6.182.0

- Added `--limit <N>` (`--limit=1`, `-n 1`) option to Chrome profile import and export commands to restrict the number of processed profiles
- Added `--profile <name>` (`-p <name>`) option to import/export single Chrome profiles from multi-profile snapshots and archives (JSON, YAML, SQLite, ZIP)
- Enhanced `gitmap agy clear` with pinned project protection: pinned Antigravity projects are strictly protected and never deleted during cleanup
- Verified and audited `gitmap agy ls`, `gitmap agy pin-projects` (add/rm/ls), and `gitmap agy clear`
- Resolved self-update downgrade bug where `gitmap update` fetched outdated release `v6.176.0` from GitHub Releases API
- Bumped version to `v6.182.0` across all Single Source of Truth manifests (`version.json`, `package.json`, `gitmap/constants/constants.go`, `readme.md`, `changelog.md`)
- Flattened nested conditionals in `cmd/chrome_batch.go` and `cmd/chromeprofile_zip_import.go` adhering to zero-nesting standards
- Verified 100% green pass on all 16 CI/CD quality gates locally and synchronized release assets
