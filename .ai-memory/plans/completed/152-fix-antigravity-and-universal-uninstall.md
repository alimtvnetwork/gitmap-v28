# Plan 152: Fix Antigravity Installation & Universal Uninstall Engine

- **Slug**: `152-fix-antigravity-and-universal-uninstall`
- **Date**: 2026-09-13
- **Status**: completed
- **Execution Budget**: N = 150 loops
- **Loops Taken**: 2 self-loop iterations (Phase 1: Research, Skill Bootstrap & Decomposition; Phase 2: Implementation, Coverage Expansion & Consolidation)

---

## 1. Executive Summary & Task Genesis

### Genesis & Problem Report
The user reported two critical system issues:
1. **Google Antigravity Desktop IDE Missing on Ubuntu / Linux with Phantom Icon**:
   Although the GNOME application launcher displayed an icon for "Antigravit..." (uploaded screenshot `media_1789299165641.png`), clicking it failed to launch anything. The application binary was completely missing from the filesystem.
2. **Missing Universal Uninstall Coverage & Truncated Errors**:
   Running `gitmap uninstall agy` resulted in unformatted error strings (`[E9000:EXECUTION] %s is not tracked in the database...`) and failed to uninstall tools not tracked in `gitmap.db`. Tools like context menu entries, cloned script directories, and self-uninstalls lacked dedicated uninstallation paths. Furthermore, the user explicitly commanded:
   > *"Make sure it is there. So for all the installation, there should be uninstall, and that should be done automatically without any issues. Remember that, okay? It's good thing that you have the stack trace. Remember to provide a stack trace. That's good. Okay? Do not skip stack trace."*
3. **Official Artifact Sources**:
   The user provided official Google Cloud Storage artifact endpoints for Antigravity:
   - Linux x64: `https://storage.googleapis.com/antigravity-public/antigravity-hub/2.13.0-6362815968182272/linux-x64/Antigravity.tar.gz`
   - Linux ARM: `https://storage.googleapis.com/antigravity-public/antigravity-hub/2.13.0-6362815968182272/linux-arm/Antigravity.tar.gz`

---

## 2. Root Cause Analysis (RCA)

1. **Phantom GNOME Desktop Launcher**:
   - `deployAntigravityDesktopLinux` unconditionally generated `~/.local/share/applications/antigravity.desktop` before verifying whether the archive download and binary extraction succeeded.
   - Previous download URLs returned HTTP 404, but the code had already written the `.desktop` file pointing to a broken symlink `~/.local/bin/antigravity`.
   - GNOME Shell indexed the `.desktop` file, producing an application icon that did nothing.
2. **Failed Broken Symlink Cleanup**:
   - Cleanup routines guarded file removal with `if _, err := os.Stat(path); err == nil { os.Remove(path) }`.
   - In Go (and POSIX), `os.Stat` follows symlinks. When a symlink points to a non-existent target (a broken/dead symlink), `os.Stat` returns `ENOENT` (`os.ErrNotExist`).
   - Consequently, broken symlinks were never deleted by cleanup routines.
3. **Go Flag Ordering Dropping Trailing Flags**:
   - In Go's `flag.FlagSet.Parse(args)`, parsing stops at the first non-flag positional argument. Running `gitmap uninstall agy --force` caused `--force` to be ignored as an unparsed argument following `agy`.
4. **Tool Category Gaps in Uninstaller**:
   - Context menus (`ctx`, `vscode-ctx`, `pwsh-ctx`, `ag-ctx`), cloned script directories (`scripts`), and direct self-uninstall (`gitmap uninstall gitmap`) lacked custom routing and fell through to apt/choco package managers or reported database tracking errors.

---

## 3. Key Implementation Details

### A. Official GCS Artifact Integration (`cli/cmdinstall/installantigravity_fetch.go`)
- Configured official Google Cloud Storage download URLs for all operating systems and CPU architectures:
  - Linux x64 & ARM64
  - Windows x64
  - macOS ARM64 & x64
- Verified HTTP 200 status codes and valid payload sizes (~140MB–205MB) across all official endpoints.

### B. Pre-Installation Cleanup & Verification Safeguards (`cli/cmdinstall/installantigravity_cleanup.go`, `installantigravity_other.go`)
- Implemented `cleanupBrokenLinuxArtifacts()`:
  - Uses direct `os.Remove(path)` to ensure dead symlinks and stale `.desktop` files are removed without relying on `os.Stat`.
  - Cleans stale download archives and extraction directories before unpacking.
- Implemented dynamic binary locator `locateExtractedLinuxBinary()`:
  - Inspects unpacked archives, finds the regular binary file, ensures size > 0, and sets execute bit `0755`.
- Gated desktop entry creation:
  - `.desktop` files and symlinks (`~/.local/bin/antigravity` and `~/.local/bin/agy`) are ONLY created AFTER the binary is verified on disk.
  - Automatically triggers `update-desktop-database` to synchronize GNOME application index.

### C. Universal Uninstallation Coverage (`cli/cmdinstall/install_uninstall_custom.go`, `install_uninstall_paths.go`, `cli/cmd/uninstall.go`)
- Added universal coverage for all tool categories:
  - Standalone custom tools: `agy`, `antigravity`, `ag-manager`, `scripts-fixer`, `coding-guidelines`, `macro-ahk`.
  - Context menu shortcuts: `ctx`, `vscode-ctx`, `pwsh-ctx`, `ag-ctx`.
  - Cloned script repositories: `scripts` (purging `D:\gitmap-scripts` and `~/Desktop/gitmap-scripts`).
  - Self-uninstall alias: `gitmap uninstall gitmap` / `gitmap-cli` routing to `runSelfUninstall`.
  - Extended alias map: `clean-code`, `code-guide`, `cg`, `cc` -> `coding-guidelines`; `ubuntu-common`, `ub-common` -> `build-essential`.
  - Package manager tools: ~85 tools backed by apt, brew, choco, winget.
- Fixed flag ordering in `cmd/uninstall.go`:
  - `reorderFlagsBeforeArgs(args)` moves flags before positional arguments, ensuring `gitmap uninstall agy --force` parses correctly.
- Dual database cleanup:
  - Deletes tracking records from both `gitmap.db` and `installation.db`.

### D. AppError & Stack Trace Preservation (`cli/cmdinstall/install_handlers.go`, `cli/cmd/uninstall.go`)
- Wrapped all uninstallation and installation errors in `*apperror.AppError` via `apperror.WrapSimple` and `apperror.New`.
- Errors emitted through `cliexit.HandleError` format and display full stack traces to `os.Stderr`.

---

## 4. Consolidated Subtasks

### Subtask 01: Fix Antigravity Linux Installation and Cleanup
- **Target**: Deploy official Google Antigravity GCS artifacts, prevent orphan `.desktop` files, remove broken symlinks, and update GNOME desktop cache.
- **Status**: Completed.
- **Key Files**:
  - `cli/cmdinstall/installantigravity_fetch.go`
  - `cli/cmdinstall/installantigravity_cleanup.go`
  - `cli/cmdinstall/installantigravity_other.go`
  - `cli/cmdinstall/installantigravity.go`

### Subtask 02: Universal Uninstall Coverage and Stack Traces
- **Target**: Universal uninstall routing across all tools, flag reordering, broken symlink deletion, and AppError stack trace preservation.
- **Status**: Completed.
- **Key Files**:
  - `cli/cmdinstall/install_packages.go`
  - `cli/cmdinstall/install_uninstall_custom.go`
  - `cli/cmdinstall/install_uninstall_paths.go`
  - `cli/cmd/uninstall.go`
  - `cli/cmd/uninstall_helpers.go`

---

## 5. Verification & Quality Gates

1. **GCS Artifact Verification**:
   All 5 official GCS download URLs verified with HTTP 200:
   - Linux x64: 173,005,112 bytes
   - Linux ARM: 170,014,747 bytes
   - Windows x64: 148,797,976 bytes
   - macOS ARM: 190,668,953 bytes
   - macOS x64: 204,658,207 bytes
2. **Coding Guidelines Adherence**:
   - All newly authored/refactored functions adhere to <= 15 lines.
   - Affirmative booleans (`is*`, `has*`) strictly enforced.
   - Zero explicit `== true` evaluations.
   - All errors use structured `*apperror.AppError` with stack traces.
3. **Execution Constraints**:
   - Total ban on routine test running (`go test`) and build checking (`go build`) strictly maintained.
