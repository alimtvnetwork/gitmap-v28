# RCA-101: Ubuntu Pull-All Remediation Engine, Oh-My-Zsh Scanner Exclusions, and Cross-OS Path Healing

**Date:** 2026-10-05
**Status:** ✅ Resolved
**Severity:** High (Missing Repo Remediation Silent Failure, Cross-OS Separator Corruption & Scanner Noise)
**Affected Components:** `cli/cmdpull`, `cli/scanner`, `cli/constants`, `cli/fsutil`, `cli/store`, `repo-secrets/05-scripts`

---

## Part 1: Symptom & Environmental Manifestation

### Environmental Context
- **Nodes Affected:** Remote Ubuntu Linux fleet server node (`u1`) running GitMap CLI (`v6.474.0` / `v6.475.0`) in a distributed environment synchronized with Windows workstations.
- **Filesystem Context:** Linux ext4 filesystem (strictly case-sensitive, forward slash `/` separators) paired with a SQLite repository catalog (`gitmap.db`) originally cataloged or populated from Windows workstations (NTFS case-insensitive, backslash `\` separators).

### Observable Failure Symptoms
During automated fleet synchronization (`gitmap pa` / `gitmap pull-all`), 15 repository pull failures occurred:
1. **Directory Resolution Failures (14 Repositories):**
   ```text
   • <repo> failed
       ↳ Reason: Repository directory does not exist on disk
       ↳ Next Step: gitmap clone <repo>
       ↳ Option 1 (Inspect Status): gitmap status <repo>
       ↳ Option 2 (Re-pull Repo): gitmap pull <repo>
   ```
2. **Transient Conflict (1 Repository):**
   - `letsmarknow-v2` stalled due to a transient merge conflict.
3. **Interactive Remediation Silent Failure:**
   - When users selected interactive batch fix Option 1 (*Resolve all failed repositories at once*), all missing repositories failed immediately with `✗ [<repo>] Remediation failed.` without stack traces, diagnostic context, or actionable repair instructions.
4. **Shell Dotfile Pollution:**
   - The shell framework directory `.oh-my-zsh` was discovered, indexed as a repository in `gitmap.db`, and polled during batch pull operations, causing detached HEAD errors.
5. **Zero-Time Diagnostic Timestamps:**
   - In `gitmap-pull.db`, error diagnostic records in `pull_errors` rendered timestamps as `0001-01-01 00:00:00 UTC` because datetime strings could not be parsed by single-format parsers.

---

## Part 2: Root Causes

### 1. Windows Backslash Cross-Pollination
- **Mechanism:** 11 of the 14 reported missing repositories were physically present under `/home/a/git-work/`, but their `AbsolutePath` records in `gitmap.db` contained Windows backslash separators (`\`) (e.g., `/home/a/git-work/03-aukgo\core`).
- **Impact:** On Linux ext4 filesystems, the backslash `\` is a valid literal filename character rather than a directory separator. Calls to `os.Stat` and `git -C <path>` treated the path as a non-existent single file, throwing `No such file or directory`.

### 2. Linux ext4 Filesystem Case-Sensitivity and Version Suffix Divergence
- **Mechanism:**
  - Case Discrepancy: `Antigravity-Manager` existed on disk in lowercase (`antigravity-manager`). On Windows NTFS, case is preserved but insensitive; on Linux ext4, case mismatches result in `ENOENT`.
  - Version Suffix Discrepancy: The repository `movie-cli-v8` existed physically as `movie-cli`. Without automatic suffix reconciliation, GitMap flagged the path as non-existent.

### 3. Scanner Omission of Framework Dotfiles (Oh-My-Zsh)
- **Mechanism:** `DefaultConfig().ExcludeDirs` and `cli/scanner/scanner.go` lacked default exclusions for common shell dotfiles and developer tool trees.
- **Impact:** Running `gitmap scan ~` recursively walked `~/.oh-my-zsh`, discovering multiple internal plugin and theme git repositories and inserting them into `gitmap.db`.

### 4. Remediation Engine Blind Execution & UI Disconnect
- **Mechanism:** `remediateSingleRepo` in `cli/cmdpull/pull_remediation.go` unconditionally ran `git -C <path> pull` for all failures. When the path did not exist on disk, `git -C` failed immediately without attempting a clone fallback or offering database pruning.
- **Tree Formatting Disconnect:** Failure cards displayed flat `↳` entries recommending `gitmap status` and `gitmap pull` on repositories that did not exist on disk.

### 5. Rigid SQLite Datetime Parsing
- **Mechanism:** The SQLite error storage layer in `cli/store/pull_split_db_errors.go` parsed datetime strings using a single hardcoded layout (`2006-01-02 15:04:05`).
- **Impact:** When modern SQLite drivers or API clients wrote ISO8601 or RFC3339 formatted timestamps (e.g. `2026-10-05T06:14:00.123456789Z`), parsing failed silently, falling back to Go's zero-value `time.Time{}` (`0001-01-01 00:00:00 UTC`).

---

## Part 3: Integrated Multi-Tier Solutions

### 1. Centralized Scanner Exclusions & Force Flag (`cli/constants`, `cli/scanner`)
- Codified `DefaultScanExcludeDirs` in `cli/constants/constants_scan.go` including `.oh-my-zsh`, `oh-my-zsh`, `ohmyzsh`, `.ohmyzsh`, `omz`, `.omz`, `node_modules`, `vendor`, `.cache`, `.venv`, and `.git`.
- Integrated exclusion sets across `cli/scanner/scanner.go`, `cli/model/record.go`, `cli/fsutil/child_repos.go`, and `cli/fsutil/recursive_top_level.go`.
- Added `--force-include` flag to `gitmap scan` and `--force` to `gitmap add` for explicit user overrides.

### 2. Pre-Pull Path Auto-Healing (`cli/cmdpull/pull_path_heal.go`)
- Implemented `HealRepoPathSeparators` to normalize `\` to `/` on Unix and translate Windows drive letters (`D:/` -> `/mnt/d/` or stripping prefix).
- Implemented `AutoHealRepoPathInDB` and `TryHealMissingRecordPath` to update `gitmap.db` before pull workers dispatch.

### 3. Missing Repository Remediation Engine (`cli/cmdpull/pull_missing_remediate.go`)
- Upgraded `remediateSingleRepo` to differentiate between dirty working trees, network timeouts, and missing directories.
- For missing directories, automatically clones from recorded remote URLs (`HttpsUrl` / `SshUrl`) using `gitutil.Clone()`.
- Provides non-destructive database-only removal option via `DeleteRepoRecordByPath`.

### 4. Hierarchical Diagnostic Tree (`cli/cmdpull/pull_efficient_render.go`)
- Formatted failure cards using standard box-drawing connectors (`├──`, `└──`, `│   `).
- Rendered tailored dual options for missing repositories:
  - Option 1 (Clone): `gitmap clone <repo>`
  - Option 2 (Prune DB): `gitmap rm --db-only <repo>`
- Output standardized diagnostic terminal hints:
  ```text
  └── To inspect stack trace: gitmap pull-error <repo> (or: gitmap pe)
  ```

### 5. Structured Dual-Persistence Error Logging & Multi-Layout Datetime Engine
- **JSONL Append Logger (`cli/cmdpull/pull_error_logger.go`):**
  - Appends structured records to `.gitmap/logs/pull-errors.log` synchronized with `pullErrorLogMu`.
  - Captures `ErrorID`, `RepoSlug`, `RepoPath`, `NodeID`, `NodeVersion`, `ErrorType`, `ErrorText`, `StackTrace`, and `CreatedAt`.
- **Split-DB Pipeline (`cli/cmdpull/pull_db_sync.go`):**
  - Captures `NodeID`, `NodeVersion`, `CreatedAt`, and `StackTrace` into `gitmap-pull.db` (`pull_errors` table).
- **Multi-Layout Parser (`cli/store/pull_split_db_errors.go`):**
  - Implemented `ParseFlexibleDBTimestamp` supporting `RFC3339Nano`, `RFC3339`, SQLite standard (`2006-01-02 15:04:05.999999999`), ISO8601 (`2006-01-02T15:04:05`), and date-only layouts.
  - Guarantees non-zero UTC return values on any input string or driver value representation.

### 6. Automated Remote Node Healing Tooling (`repo-secrets/05-scripts/heal-u1-pull-errors.py`, `.sh`)
- Automated database backup (`gitmap.db.bak.YYYYMMDD_HHMMSS`) prior to any mutation.
- Path normalization: converts backslashes to POSIX forward slashes and strips Windows drive prefixes.
- Casing & suffix reconciliation: resolves `Antigravity-Manager` -> `antigravity-manager` and `movie-cli-v8` -> `movie-cli`.
- Deletes `.oh-my-zsh` entries from SQLite table `Repo` and appends ignore patterns to `.gitmapignore`.
- Scans candidate work directories and prunes stale `.git/index.lock` files older than 600 seconds.
- Performs physical filesystem audit reporting verified vs missing repositories.

---

## Part 4: Verification & Preventive Governance

### Automated Linter Verification
- `python linter-scripts/check-nested-ifs.py --all`: ✅ 0 violations across 4,186 files.
- `python linter-scripts/check-boolean-guidelines.py --all`: ✅ 0 violations across 4,186 files.
- `python linter-scripts/check-enum-and-boolean.py`: ✅ 0 violations across 3,195 files.
- `python linter-scripts/check-error-management.py --all`: ✅ 0 violations across 4,246 files.
- `python linter-scripts/check-relative-paths.py`: ✅ 0 violations across 8,678 files.

### Script Syntax & Compilation
- `python -m py_compile repo-secrets/05-scripts/heal-u1-pull-errors.py`: ✅ Clean compilation (exit code 0).
- `repo-secrets/05-scripts/heal-u1-pull-errors.sh`: ✅ POSIX shell wrapper with argument forwarding and executable permissions.

### Remote Fleet Node (`u1`) Validation
1. **Pre-Healing Inspection:**
   - 11 backslash paths detected.
   - 1 stale `.oh-my-zsh` record detected.
   - 14 repositories reported missing.
2. **Post-Healing Execution:**
   - Backup created: `gitmap.db.bak.<timestamp>`
   - 11 paths normalized to POSIX forward slashes.
   - 2 casing/suffix paths reconciled (`antigravity-manager`, `movie-cli`).
   - `.oh-my-zsh` deleted from `Repo` table and added to `.gitmapignore`.
   - Physical disk audit: 26 repositories verified on disk (0 missing).
3. **End-to-End Pull Verification:**
   - Remote execution `gitmap ssh exec u1 "gitmap pa"`: exited with returncode `0` and 0 failures.
   - Remote error query `gitmap ssh exec u1 "gitmap pe all --json"`: returned empty list `[]`.

### Preventive Governance Protocols
1. **Pre-Pull Path Assertion:** All pull executors automatically normalize path separators before invoking git subprocesses.
2. **Scanner Filter Guard:** Built-in exclude lists prevent registration of shell and tool directories.
3. **Dual-Persistence Observability:** All pull errors are written to both JSONL files and indexed SQLite tables with complete stack traces, enabling instant diagnostic retrieval via `gitmap pe`.
