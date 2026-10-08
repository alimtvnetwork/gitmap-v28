# GitMap v6.513.0

## What's Changed in v6.513.0

### Highlights & Key Improvements
- **Diffstat & Commit Counting Isolation**: Resolved issue where multiple repositories displayed identical diffstat metrics `+167988/-14766 (1095)` in pull/status summaries. Added strict empty-path guards and directory checks across `cli/gitutil/last_sha.go` and `cli/cmdpull/pull_worker.go`, while ensuring `cmd.Dir = cleanDir` in `cli/cmdpull/pull.go`.
- **Remediation Sub-Node Tree UX Restructuring**: Restructured batch pull failure remediation in `cli/cmdpull/pull_remediation.go` and `cli/cmdpull/pull_missing_remediate.go` to render each repository as a parent node (`• <repo>`) with indented tree sub-nodes (`├── Attempting auto-merge...`, `└── ✓ Successfully updated.`, `└── ✗ Merge conflict detected...`).
- **`gitmap download` Subcommand with aria2c Acceleration**: Implemented unified `gitmap download <url>` (alias `gitmap dl`) in `cli/cmddownload/` with automatic `aria2c` detection (`-s 16 -x 16 -k 1M`), 3-tier fallback to `curl` and Go `net/http`, and centered symmetric terminal progress bars.
- **Smart `gitmap update` Engine**:
  - Pinned version targeting via argument `gitmap update [version]` and `--version <v>`.
  - Short-circuit fast-path when already on the latest/target version (`Already updated`).
  - Automatic fallback across up to 5 previous tags if release executables are absent or update fails.
  - Split SQLite `ReleaseCache` table in `cli/store/release_cache.go` (`ReleaseCacheId`, `Tag`, `Version`, `HasExecutables`, `AssetUrl`, `CheckedAt`, `ExpiresAt`).
  - Clean console output showing progress and percentage, plus full structured `--json` payload support.
  - Synchronized runner scripts (`install.ps1`, `install.sh`, `gitmap.ps1`, `run.ps1`).

### Specifications & Documentation
- Architectural Specification: `02-spec/21-app/242-aria2c-download-smart-update-stats-subnode-remediation/01-architecture-spec.md`
- Component Specification: `02-spec/21-app/242-aria2c-download-smart-update-stats-subnode-remediation/02-download-and-update-spec.md`
- Completed Plan: `.ai-memory/plans/completed/242-aria2c-download-smart-update-stats-subnode-remediation.md`

### Quick Install One-Liners

**Windows (PowerShell):**
```powershell
irm https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.513.0/install.ps1 | iex
```

**Linux / macOS (Bash):**
```bash
curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.513.0/install.sh | sh
```
