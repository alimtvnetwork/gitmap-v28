# Architecture Specification: Ubuntu Pull-All Remediation, Oh-My-Zsh Exclusions & Cross-OS Path Healing

- **Task Slug:** `219-ubuntu-pull-all-remediation-omz-ignore-and-error-logging`
- **Application Version:** `6.475.0`
- **Scope:** Cross-Platform Scanner Exclusions, Oh-My-Zsh Discovery Isolation, Cross-OS Backslash Path Normalization, SQLite Split-DB Auto-Healing, and Enhanced Missing-Repo Pull Remediation with Clone Fallback.

---

## 1. Executive Summary & Problem Analysis

In multi-node development fleets and heterogenous environments combining Windows workstations and Linux/Ubuntu servers, GitMap encountered three critical failure modes during repository discovery and batch pull operations (`gitmap pa` / `gitmap pull --all`):

### 1.1 Oh-My-Zsh & Tool Subtree Discovery Pollution
On Linux/Ubuntu nodes, developers typically install Oh-My-Zsh (`~/.oh-my-zsh`), package manager caches, and Python virtual environments directly within user home directories. 
- Running `gitmap scan ~` recursively descends into `~/.oh-my-zsh/` and its subdirectories (`plugins/`, `custom/themes/`, `tools/`).
- Because OMZ and its plugins are themselves git checkouts (often detached HEADs, read-only submodules, or locked trees), GitMap registers dozens of noise repositories in `gitmap.db`.
- Subsequent `gitmap pa` commands attempt to pull all these internal framework repos, triggering authentication warnings, detached HEAD merge failures, and long timeouts.

### 1.2 Cross-OS Path Sensitivity & Backslash Breakage
When `gitmap.db` is copied, synchronized (via SSH cluster sync, dotfiles, or backup restores), or shared across OS boundaries:
- Repositories discovered on Windows contain backslash-separated absolute paths (e.g., `D:\projects\app` or `C:\Users\admin\dev\repo`).
- On Unix/Ubuntu, backslashes (`\`) are treated as literal characters within filename strings rather than path component delimiters.
- Calling `os.Stat(rec.AbsolutePath)` or invoking `git -C <path> pull` fails with:
  `fatal: cannot change to 'D:\projects\app': No such file or directory`.
- GitMap lacked an automated, pre-flight auto-healing mechanism in SQLite to normalize Windows backslashes to Unix forward slashes.

### 1.3 Missing Repository Pull Breakdowns & Remediation Gap
When a tracked repository directory is missing on disk (`cloner.IsMissingRepo` returns true):
- The existing pull worker marks the repo as `PullStepTypeSkipped` with error `"missing repository directory"`.
- If remediation is invoked (`gitmap pull-fix` or interactive remediation menu in `cli/cmdpull/pull_remediation.go`), the engine blindly runs:
  `git -C <missing_path> pull --progress --no-rebase --no-edit --autostash`.
- This inevitably fails because the target directory does not exist.
- GitMap lacked an intelligent remediation fallback that:
  1. Attempts path healing and cross-drive relocation.
  2. Offers automated cloning from known remote URLs (`rec.HTTPSUrl`, `rec.SSHUrl`, `rec.DiscoveredURL`).
  3. Offers safe database-only removal (`DELETE FROM Repo WHERE RepoId = ?`) so phantom records stop poisoning batch pull workflows.

---

## 2. Architectural Blueprint & Data Flow

```
+-----------------------------------------------------------------------------------------------+
|                                      DISCOVERY & SCAN ENGINE                                  |
+-----------------------------------------------------------------------------------------------+
                                                |
               +--------------------------------+--------------------------------+
               |                                                                 |
     [ gitmap scan <path> ]                                           [ gitmap add <path> ]
               |                                                                 |
    Filter with Default Excludes                                      Check Manual Exclude Gate
    (OMZ, node_modules, vendor, .cache)                                          |
               |                                                       Is Target Dir Excluded?
    Has --force-include flag?                                                    |
         /           \                                                  /         \
       Yes            No                                              Yes          No
       /               \                                              /             \
  Bypass Match     Skip Traversal                              Has --force flag?   Add Record
                          |                                          /         \
                   Prune Directory                                  Yes         No
                                                                    /            \
                                                            Bypass Filter     Reject with Hint

+-----------------------------------------------------------------------------------------------+
|                                     CROSS-OS PRE-PULL AUTO-HEAL                               |
+-----------------------------------------------------------------------------------------------+
                                                |
                                     [ gitmap pa / pull ]
                                                |
                                   HealCrossOSPaths(db)
                                                |
                    +---------------------------+---------------------------+
                    |                                                       |
            Running on Unix?                                      Running on Windows?
             (GOOS != "windows")                                    (GOOS == "windows")
                    |                                                       |
        Detect Windows Backslashes (\)                             Detect Mixed Slashes
        Convert to Unix Slashes (/)                                Normalize to filepath.Clean
                    |                                                       |
        Verify Path on Disk                                        Verify Path on Disk
                    |                                                       |
        Update Repo.AbsolutePath in SQLite                         Update Repo.AbsolutePath in SQLite

+-----------------------------------------------------------------------------------------------+
|                               ENHANCED PULL REMEDIATION ENGINE                                |
+-----------------------------------------------------------------------------------------------+
                                                |
                                    [ ExecuteTrackedPull ]
                                                |
                                    Is Repository Missing?
                                       (IsMissingRepo == true)
                                                |
                     +--------------------------+--------------------------+
                     |                                                     |
             Attempt Path Heal                                     Directory Not Found
                     |                                                     |
              Does Path Exist?                                     Check Remote URLs
                /          \                                       (HTTPSUrl / SSHUrl)
              Yes           No                                             |
              /              \                              +--------------+--------------+
      Update Path         Proceed to                        |                             |
      Execute Pull        Clone Fallback              Remote Exists                No Remote URL
                                                            |                             |
                                                     [ Option 1: Clone ]           [ Option 2: Remove ]
                                                     git clone <url> <path>        DELETE FROM Repo
                                                            |                      WHERE RepoId = ?
                                                     Re-link & Pull Up-to-date
```

---

## 3. Core Architectural Pillars

### 3.1 Scanner Exclusions & Oh-My-Zsh Isolation
1. **Centralized Exclusion Constants (`cli/constants/constants_scan.go`):**
   - Define canonical exclude slice `DefaultScanExcludeDirs`:
     ```go
     var DefaultScanExcludeDirs = []string{
         ".oh-my-zsh",
         "oh-my-zsh",
         "ohmyzsh",
         ".ohmyzsh",
         "omz",
         ".omz",
         "omizssh",
         "node_modules",
         "vendor",
         ".cache",
         ".venv",
         "venv",
         ".git",
         ".terraform",
         ".next",
         ".turbo",
         "dist",
         "build",
         "bin",
         "obj",
         "target",
     }
     ```
2. **Scanner Baseline Merge (`cli/scanner/scanner.go:buildExcludeSet`):**
   - `buildExcludeSet` merges `DefaultScanExcludeDirs` with user-specified directories loaded from `gitmap.config.json` or CLI flags.
   - When a folder matches an exclusion, `filepath.SkipDir` or worker skip is triggered immediately without examining children.
3. **Application Configuration Alignment (`cli/model/record.go`):**
   - Update `DefaultConfig()` to ensure `ExcludeDirs` contains the baseline set by default.
4. **Filesystem Utility Synchronization:**
   - `cli/fsutil/child_repos.go:DiscoverChildGitRepos`: Check each directory entry against the exclusion set before running `os.Stat(filepath.Join(subPath, ".git"))`.
   - `cli/fsutil/recursive_top_level.go:DiscoverTopLevelGitRepos`: Replace hardcoded string checks (`.git`, `node_modules`, `.cache`, `vendor`) with unified lookups against `constants.DefaultScanExcludeDirs`.
5. **Explicit Force Override Mechanisms:**
   - `--force-include <pattern|all>` on `gitmap scan`:
     Allows developers who legitimately manage Oh-My-Zsh plugins or vendor packages to include them on demand.
     Example: `gitmap scan ~/.oh-my-zsh --force-include .oh-my-zsh,omz`.
   - `--force` on `gitmap add`:
     Allows registering a specific repository that would otherwise match default exclusion rules.
     Example: `gitmap add ~/.oh-my-zsh/custom/plugins/my-plugin --force`.

---

### 3.2 Cross-OS Path Normalization & Pre-Pull Auto-Healing

1. **Path Symptom on Linux/Ubuntu:**
   - A row in `gitmap.db` contains `AbsolutePath: "D:\\work\\project"` or `"/home/ubuntu\\sub\\repo"`.
   - On Linux, backslash `\` is not a path separator. `filepath.Dir` or `os.Stat` on `"/home/ubuntu\\sub\\repo"` looks for a single file named `sub\repo` inside `/home/ubuntu`, which fails.
2. **Path Auto-Healing Algorithm (`cli/store/repo_sanitize.go`):**
   - Implement `HealCrossOSPaths(db *sql.DB, isQuiet bool) (int, error)`:
     * Query all repositories: `SELECT RepoId, AbsolutePath FROM Repo`.
     * Inspect if `runtime.GOOS != "windows"` and `strings.Contains(path, "\\")`.
     * Convert: `normalized := strings.ReplaceAll(path, "\\", "/")`.
     * Windows drive adaptation:
       - If path matches `^[A-Za-z]:/`: check if `/mnt/<drive_lower>/...` exists (WSL) or if the path tail relative to project roots exists.
     * Disk verification: If `cloner.IsGitRepo(normalized)`:
       - Execute `UPDATE Repo SET AbsolutePath = ? WHERE RepoId = ?`.
       - Increment healed counter.
3. **Integration Point in Pull Workflow (`cli/cmdpull/pull_dedup.go`):**
   - In `runPullAll` and `runPull`, before pulling worker channels, execute `HealCrossOSPaths(db, opts.quiet)` in tandem with `OptimizeRedundantRepos(db, opts.quiet)`.
   - Output clear visual telemetry:
     `✓ SQLite: Auto-healed 4 cross-OS repository path(s) to Unix slashes before pull.`

---

### 3.3 Enhanced Pull Remediation & Automatic Clone Fallback

1. **Missing Repo Detection in Worker (`cli/cmdpull/pull_worker.go`):**
   - When `cloner.IsMissingRepo(rec.AbsolutePath)` evaluates to true:
     * **Step 1 (Path Self-Heal):** Attempt real-time in-memory backslash-to-slash replacement. If the adjusted path exists on disk, update `rec.AbsolutePath` and continue normal pull stream.
     * **Step 2 (Missing State Tagging):** If directory is physically absent, record `rec.HTTPSUrl` / `rec.SSHUrl` in the `PullRepoState` failure metadata.
2. **Missing Repo Remediation Engine (`cli/cmdpull/pull_remediation.go`):**
   - In `remediateSingleRepo(f PullFailureSummary)`:
     * Do NOT execute `git -C <missing_path> pull`.
     * Check if `cloner.IsMissingRepo(f.RepoPath)` is true.
     * If missing, invoke `remediateMissingRepo(f)`.
3. **Remediation Choices for Missing Repositories:**
   - **Action A: Clone Fallback (Automatic or User-Confirmed):**
     * Query repository remote from `gitmap.db` using `f.RepoPath` or `f.RepoName`.
     * Verify remote URL is valid (`https://...` or `git@...`).
     * Execute `git clone <remoteURL> <f.RepoPath>`.
     * If clone succeeds:
       `✓ [repo-name] Missing directory recreated via clone fallback.`
   - **Action B: Database-Only Removal:**
     * If no remote URL is available or user selects prune:
       `DELETE FROM Repo WHERE AbsolutePath = ? OR Slug = ?`.
     * Report:
       `✓ [repo-name] Stale entry removed from database (no disk files modified).`
   - **Action C: Skip:**
     * Preserve record and bypass without raising panic.
4. **Interactive & Batch Remediation Menus:**
   - Update `printRemediationPromptMenu` and `promptAndRemediateSingle` to distinguish between missing directory errors vs. dirty/merge conflict errors.
   - For missing repos, offer:
     `[c]lone from remote / [d]elete record from db / [s]kip`.
   - For batch all mode (`RunBatchRemediateAll` / `gitmap pull-fix --all`):
     * If remote exists: attempt clone fallback.
     * If no remote: skip and warn with clear remediation hint.

---

## 4. Detailed Component Changes & File Touches

| Component / File | Purpose & Architectural Responsibility |
| :--- | :--- |
| `cli/constants/constants_scan.go` | Declare `DefaultScanExcludeDirs`, `--force-include`, and `--force` CLI constants. |
| `cli/scanner/scanner.go` | Update `buildExcludeSet` to merge default scan exclusions; add `ForceIncludeDirs` to `ScanOptions`. |
| `cli/model/record.go` | Update `DefaultConfig()` to ensure baseline exclusions are populated. |
| `cli/fsutil/child_repos.go` | Filter entries against `constants.DefaultScanExcludeDirs` before checking `.git`. |
| `cli/fsutil/recursive_top_level.go` | Replace hardcoded exclude list with `constants.DefaultScanExcludeDirs`. |
| `cli/cmdscan/flags.go` | Register `--force-include` string flag on scan flag set. |
| `cli/cmdscan/scan.go` | Propagate `--force-include` tokens to scanner exclude filter. |
| `cli/cmd/rootadd.go` | Add `--force` flag handling to allow registering excluded repos explicitly. |
| `cli/store/repo_sanitize.go` | Implement `HealCrossOSPaths(db *sql.DB, isQuiet bool)` for backslash conversion. |
| `cli/cmdpull/pull_dedup.go` | Hook `HealCrossOSPaths` into pre-pull startup sequence. |
| `cli/cmdpull/pull_worker.go` | In-memory path normalization and missing repo remote tagging. |
| `cli/cmdpull/pull_remediation.go` | Enhanced missing repo handler, clone fallback execution, and DB-only prune. |
| `cli/cmdpull/pull_not_found.go` | Updated help hints and missing repo diagnostics. |

---

## 5. CLI UX, Commands & Flag Specifications

### 5.1 Scanning with Exclusions & Force Override
```bash
# Default scan: automatically ignores .oh-my-zsh, node_modules, vendor, .cache, .venv
gitmap scan ~

# Explicitly force-include Oh-My-Zsh directories
gitmap scan ~ --force-include .oh-my-zsh,omz

# Explicitly force-include all directories (disable baseline exclusions)
gitmap scan ~ --force-include all
```

### 5.2 Manual Repository Addition with Force Flag
```bash
# Attempting to add an excluded repository gives a protective warning
gitmap add ~/.oh-my-zsh
# Output:
#   error: '/home/ubuntu/.oh-my-zsh' matches default exclusion '.oh-my-zsh'.
#   Run with --force to add anyway: gitmap add ~/.oh-my-zsh --force

# Forcing addition of an excluded directory
gitmap add ~/.oh-my-zsh --force
# Output:
#   ✓ Tracked 1 repository in /home/ubuntu/.oh-my-zsh (forced)
```

### 5.3 Batch Pull with Auto-Heal & Missing Repo Remediation
```bash
# Run batch pull on Ubuntu after syncing database from Windows
gitmap pa

# Console Output:
#   ✓ SQLite: Auto-healed 4 cross-OS repository path(s) to Unix slashes before pull.
#   ✓ SQLite: Optimized 0 redundant repository record(s).
#   Pulling 42 repositories across 8 workers...
#   ...
#   [!] 2 repositories failed or missing.

# Run batch pull remediation
gitmap pull --remediate
# Interactive Menu:
#   [?] Repository 'legacy-tool' directory is missing: /home/ubuntu/dev/legacy-tool
#       Known Remote: https://github.com/myorg/legacy-tool.git
#       [1/c] Clone from remote URL
#       [2/d] Delete record from gitmap database
#       [3/s] Skip
#   Choice: 1
#   → Cloning https://github.com/myorg/legacy-tool.git...
#   ✓ Successfully cloned and up-to-date.
```

---

## 6. Acceptance Criteria & Quality Gates

1. **Scanner Exclusions Verified:**
   - `gitmap scan <dir>` does not recurse into `.oh-my-zsh`, `ohmyzsh`, `omz`, or `node_modules`.
   - `--force-include` allows scanning excluded folders when explicitly requested.
2. **Cross-OS Auto-Heal Verified:**
   - Paths stored with Windows `\` backslashes are automatically converted to `/` when running on Unix before pull begins.
   - Database record `AbsolutePath` in `Repo` table is updated on disk verification.
3. **Missing Repo Remediation Verified:**
   - Missing repository paths never invoke `git -C <missing> pull`.
   - Remediation engine provides automatic clone fallback using remote URLs and DB-only removal for dead repos.
4. **Code Quality & Relative Paths:**
   - Zero hardcoded magic strings; all constants declared in `cli/constants/`.
   - Zero absolute filesystem paths in Go or Python scripts.
   - Clean pass on `python 03-ai-scripts/05-guideline-autofixer.py` and `python linter-scripts/check-relative-paths.py`.
