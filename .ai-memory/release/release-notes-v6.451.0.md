# GitMap v6.451.0

## What's Changed in v6.451.0

- **GitIgnore Split-DB Cache Enforcement across Pull Commands:**
  - Integrated `isCWDIgnoreCheckRecent` in `handleCWDIgnoreChecks` to query the SQLite Split-DB cache (`.gitmap/data/gitignore/cache/sql.db`) prior to any ignore audits in current working directory pulls.
  - Enforced recording of `status = "skipped"` for every repository when skipping interactive ignore resolution or running in non-interactive/CI environments, preventing re-prompt loops across repeated `gitmap pa` and `gitmap pull` commands within the 24-hour TTL.
  - Added `normalizeRepoPathForCache` across all Split-DB methods to reliably canonicalize drive letter casing, relative paths, and path separators (`/`).
  - Extended `IsCheckRecent` to treat `"clean"`, `"skipped"`, and `"remediated"` records as valid cached states.
  - Ensured clean repository statuses are recorded into the Split-DB during `fix-ignore-all` batch audits.
- **Pending Task Completion Warning Suppression:**
  - Expanded `isIgnorablePendingTaskError` in `cli/cmd/pendingtaskhelper.go` to cleanly suppress non-fatal `no rows` and `not found` errors during `CompleteTask` invocations, preventing false-positive warnings when idempotent or already-completed tasks are processed.
- **Automated Manifest Synchronization:**
  - Upgraded `03-ai-scripts/37-bump-version.py` to keep `what-to-read.md`, `cli/constants/constants.go`, and release notes synchronized across all release ceremonies.

### Quick Install One-Liners

**Windows (PowerShell):**
```powershell
irm https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.451.0/install.ps1 | iex
```

**Linux / macOS (Bash):**
```bash
curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.451.0/install.sh | sh
```
