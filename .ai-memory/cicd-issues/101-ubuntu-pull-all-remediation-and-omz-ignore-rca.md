# RCA-101: Ubuntu Pull-All Remediation Engine, Oh-My-Zsh Scanner Exclusions, and Cross-OS Path Healing

**Date:** 2026-10-05
**Status:** ✅ Resolved
**Severity:** High (Missing Repo Remediation Silent Failure, Cross-OS Separator Corruption & Scanner Noise)
**Affected Components:** `cli/cmdpull`, `cli/scanner`, `cli/constants`, `cli/fsutil`, `repo-secrets/05-scripts`

---

## 1. Symptom

On the remote Ubuntu fleet node (`u1`), running `gitmap pa` (`gitmap pull-all`) reported 15 repository failures:
- 14 repositories failed with:
  ```text
  • <repo> failed
      ↳ Reason: Repository directory does not exist on disk
      ↳ Next Step: gitmap clone <repo>
      ↳ Option 1 (Inspect Status): gitmap status <repo>
      ↳ Option 2 (Re-pull Repo): gitmap pull <repo>
  ```
- 1 repository (`letsmarknow-v2`) failed due to a transient merge conflict.
- When selecting interactive remediation option `1` (*Resolve all failed repositories at once*):
  - Every missing repository failed with `✗ [<repo>] Remediation failed.` without stack traces, error outputs, or diagnostic hints.
  - Shell framework directory `.oh-my-zsh` was tracked and polled as a repository.
  - The failure summary rendered flat `↳` entries with illogical options (`gitmap status` and `gitmap pull` on a repository that does not exist on disk).

---

## 2. Root Cause

1. **Cross-OS Windows Backslash Ingestion:**
   - 11 of the 14 repositories actually existed on disk under `/home/a/git-work/`, but their `AbsolutePath` in `gitmap.db` contained Windows backslashes (`\`) (e.g., `/home/a/git-work/03-aukgo\core`).
   - On Linux, backslashes are treated as literal characters rather than directory separators, causing `os.Stat` and `git -C` to fail with `No such file or directory`.
2. **Casing and Suffix Mismatches:**
   - `Antigravity-Manager` existed on disk in lowercase (`antigravity-manager`).
   - `movie-cli-v8` existed on disk without the version suffix (`movie-cli`).
3. **Missing Scanner Default Exclusions:**
   - `DefaultConfig().ExcludeDirs` and `cli/scanner/scanner.go` lacked built-in default directory exclusions. Running `gitmap scan ~` crawled `~/.oh-my-zsh` and registered it in `gitmap.db`.
4. **Remediation Engine Blind Execution:**
   - `remediateSingleRepo` in `cli/cmdpull/pull_remediation.go` unconditionally executed `git -C <repoPath> pull` for all failures. When `repoPath` did not exist, `git -C` failed immediately. It never attempted to clone from known remote URLs or offer database-only pruning.
5. **Tree Formatting & Options Disconnect:**
   - `resolveStateOrAuthDualHints` omitted a check for `isMissingRepoFailure`, falling back to nonsensical status/pull commands.
   - The tree renderer used flat indentation instead of hierarchical box-drawing connectors (`├──`, `└──`, `│   `).
6. **Error Logging & RFC3339 Driver Incompatibility:**
   - `pull_errors` in `gitmap-pull.db` was populated with blank `NodeID` and `StackTrace`. Modernc SQLite returned RFC3339 datetimes which failed `time.Parse("2006-01-02 15:04:05")`. No file logger existed at `.gitmap/logs/pull-errors.log`.

---

## 3. Resolution

1. **Centralized Default Scanner Exclusions:**
   - Created `DefaultScanExcludeDirs` in `cli/constants/constants_scan.go` (.oh-my-zsh, oh-my-zsh, ohmyzsh, .ohmyzsh, omz, .omz, omizssh, node_modules, vendor, .cache, .venv, etc.).
   - Integrated into `cli/scanner/scanner.go:buildExcludeSet`, `DefaultConfig` in `cli/model/record.go`, `cli/fsutil/child_repos.go`, and `cli/fsutil/recursive_top_level.go`.
   - Added `--force-include` to `gitmap scan` and `--force` to `gitmap add`.
2. **Cross-OS Path Normalization & Pre-Pull Auto-Healing:**
   - Implemented `HealRepoPathSeparators`, `AutoHealRepoPathInDB`, and `TryHealMissingRecordPath` in `cli/cmdpull/pull_path_heal.go`.
   - Automatically normalizes `\` to `/` on Unix and updates `gitmap.db` before pull workers execute.
3. **Missing Repo Remediation Engine:**
   - Enhanced `remediateSingleRepo` in `cli/cmdpull/pull_missing_remediate.go` to branch on missing directories.
   - Automatically executes `git clone <remote_url> <path>` using recorded `HttpsUrl` / `SshUrl`.
   - Provides safe database-only removal (`DeleteRepoRecordByPath`).
4. **Hierarchical Failure Subtree:**
   - Refactored `renderSingleFailedItem` and `renderStructuredOptions` in `cli/cmdpull/pull_efficient_render.go` using box-drawing connectors (`├──`, `└──`, `│   `).
   - Added `resolveMissingRepoDualHints` returning Option 1 (`gitmap clone <repo>`) and Option 2 (`gitmap rm --db-only <repo>`).
5. **Structured Error Logging & Diagnostic Hints:**
   - Implemented `pull_error_logger.go` appending JSONL entries to `.gitmap/logs/pull-errors.log`.
   - Populated `NodeID`, `NodeVersion`, and `StackTrace` in `gitmap-pull.db`.
   - Added multi-layout datetime parsing (`time.RFC3339`, `"2006-01-02 15:04:05"`) in `cli/store/pull_split_db_errors.go`.
   - Output terminal hint: `└── To inspect stack trace: gitmap pull-error <repo> (or: gitmap pe)`.
6. **Remote Ubuntu Fleet Diagnostic & Healing Script:**
   - Created `repo-secrets/05-scripts/heal-u1-pull-errors.py` (and `.sh`).
   - Executed on `u1` over SSH: normalized all 11 backslash paths, resolved casing/suffixes, and pruned stale `.oh-my-zsh`.

---

## 4. Verification

- `python linter-scripts/check-nested-ifs.py --all`: ✅ 0 violations across 4,186 files.
- `python linter-scripts/check-boolean-guidelines.py --all`: ✅ 0 violations across 4,186 files.
- `python linter-scripts/check-enum-and-boolean.py`: ✅ 0 violations across 3,195 files.
- `python linter-scripts/check-error-management.py --all`: ✅ 0 violations across 4,246 files.
- `python linter-scripts/check-relative-paths.py`: ✅ 0 violations across 8,678 files.
- `go vet ./...` in `cli/`: ✅ Clean (exit code 0).
- `golangci-lint run --issues-exit-code=1 ./...` in `cli/`: ✅ Clean (exit code 0).
- **Remote Node Audit (`u1`):** 26 repositories verified on disk (0 missing).
