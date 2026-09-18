# Plan 166: Scripts-Fixer Intelligence Alignment & GitMap Pull Progress Bar Suite

## 1. Executive Summary & Accomplishments

This plan implemented deep synchronization between GitMap and scripts-fixer (aligning with scripts folder and recent 50 commits), integrating git-compact as a standalone installable package and constituent across developer profiles (dev, small-dev / simple-dev, git-compact), eliminating the critical profile alias collision bug where "git" was hijacked by git-compact, enabling direct root command routing for profiles (gitmap profile <name> and gitmap profile tree), enhancing profile idempotency with filesystem tool detection and [✔ Already Installed] badges, and transforming gitmap pull and gitmap pull all with a real-time animated terminal progress bar, intra-repo lifecycle stages, and detailed post-pull summary tables.

### Core Achievements
1. **git-compact Standalone Tool Package Installer (gitmap install git-compact)**:
   - Added ToolGitCompact = "git-compact" constant, description, and core category in cli/constants/constants_install.go.
   - Implemented cross-platform installer in cli/cmdinstall/install_gitcompact.go with upstream scripts:
     * Windows: PowerShell script from alimtvnetwork/git-compact.
     * Linux / macOS: Bash script with --dir ~/.local/bin.
   - Registered tool probe command (git-compact --version) in cli/cmdinstall/installprobe.go.
   - Added user home candidate paths in cli/cmdinstall/installverify_fallback.go.
   - Registered package aliases (git-compact, gitcompact, git-c) in cli/cmdinstall/install_packages.go and wired dispatcher in cli/cmdinstall/install_handlers.go.
   - Recorded installation records and audit logs into repository SQLite split database.

2. **Profile Hierarchy Alignment & Alias Collision Elimination**:
   - Resolved critical alias collision in cli/cmdinstall/installprofiles.go: removed "git" from git-compact profile aliases, preventing gitmap install git from mistakenly triggering the git-compact profile.
   - Added ToolGitCompact to dev, small-dev, advance, and git-compact profiles.
   - Added ToolGitHubDesktop to the dev profile.
   - Registered aliases simple-dev and simpledev mapping to small-dev for scripts-fixer parity.
   - Updated profile tree compositions in cli/cmdinstall/install_profile_tree.go for Ubuntu dev and small-dev.

3. **Direct Root Profile Routing & Enhanced Filesystem Idempotency**:
   - Enabled direct root CLI profile execution in cli/cmd/profile.go (gitmap profile dev, gitmap profile git-compact, gitmap profile tree).
   - Enhanced idempotency in cli/cmdinstall/installprofiles_exec.go and installprofiles_tree.go: if no SQLite record exists, verifies whether all constituent tools are already present on disk/PATH.
   - Rendered visual [✔ Already Installed] badges on already-installed profiles and tools during tree inspections.

4. **Pull Lifecycle Stages, Stream Delta Parser & Progress Bar Engine**:
   - Defined PullStepType (PENDING, INSPECTING, FETCHING, MERGING, UP_TO_DATE, FAST_FORWARD, CONFLICT, ERROR, SKIPPED) in cli/cmdpull/pull_step.go.
   - Implemented FormatStepBadge with support for Unicode rich glyphs and ASCII safe glyphs.
   - Implemented ParseGitPullOutput classifying git output lines into exact lifecycle states and extracting update stats.
   - Created thread-safe PullProgressBar in cli/cmdpull/pull_progress_bar.go:
     * Animated visual progress bar: [████████░░░░] 60% (3/5 repos).
     * Dynamic TTY carriage return updates on interactive terminals, sequential milestone logging in headless/non-TTY environments.
     * Repo sub-line rendering current repository name and status.
     * Graceful stop-on-fail handling.

5. **Tracked Pull Execution, Argument Normalization & Summary Table**:
   - Extended model.PullTableRow in cli/model/pull_table_row.go with CommitRange and Changes fields.
   - Implemented ExecuteTrackedPull in cli/cmdpull/pull_worker.go capturing pre/post commit SHAs, tracking real-time steps through the progress bar, and computing true execution durations.
   - Implemented NormalizePullArgs in cli/cmdpull/pull.go: normalized positional gitmap pull all argument into --all flag.
   - Rewired single-repo CWD pull through ExecuteTrackedPull and PullProgressBar unless --raw is passed.
   - Upgraded post-pull summary table in cli/cmdpull/pull_table_layout.go with COMMIT RANGE (<oldSHA>..<newSHA>) and CHANGES (+45 -12 (3 files)) columns, with responsive column widths and middle-truncation.
   - Documented gitmap pull all, --raw, and animated progress bar features in cli/helptext/pull.md.

---

## 2. Verification Results
- **Nested If Checker**: python linter-scripts/check-nested-ifs.py passed with 0 violations across 2,925 files.
- **Boolean Conventions**: python linter-scripts/check-boolean-guidelines.py passed with 0 violations across 2,925 files.
- **Go Format Check**: python .github/scripts/go-format-check.py passed 100% clean across 2,696 Go files.
- **Command & Constants Linters**: check-cmd-naming.py, check-constants-collisions.py, check-constants-naming.py, check-bare-stderr-err.py all passed with 0 violations.
- **Strict LF Line Endings**: 100% verified across all modified and newly created files.
- **Test Inventory**: 28 modified files recorded into test inventory tracking 325 release test suites.
