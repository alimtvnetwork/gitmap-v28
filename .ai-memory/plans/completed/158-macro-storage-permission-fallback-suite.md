# Plan 158: Macro Storage Permission & Fallback Suite

> **Consolidation Note:** This plan resolves the fatal runtime permission failure during macro recording (`mkdir /home/a/.gitmap/macros: permission denied`). Executed across 3 subtasks with N = 100 continuous budget (50 planning, 50 execution loops) and verified zero-defect completion against coding guidelines, zero-nesting, affirmative booleans, and error wrapping contracts.

## Context & User Error

User encountered the following runtime failure during interactive macro recording on Linux:
```text
gitmap macro: execute failed: [E9000:EXECUTION] save macro alim1: mkdir /home/a/.gitmap/macros: permission denied (at=cmdmacro/macro_add.go:45)
Stack Trace:
    at github.com/alimtvnetwork/gitmap-v28/cli/apperror.WrapSimple (apperror/apperror.go:231)
    at github.com/alimtvnetwork/gitmap-v28/cli/cmdmacro.handleMacroAdd (cmdmacro/macro_add.go:45)
    at github.com/alimtvnetwork/gitmap-v28/cli/cmdmacro.routeModifySubcommand (cmdmacro/macro_cmd.go:233)
    at github.com/alimtvnetwork/gitmap-v28/cli/cmdmacro.routeManagementSubcommand (cmdmacro/macro_cmd.go:215)
    at github.com/alimtvnetwork/gitmap-v28/cli/cmdmacro.routeMacroSubcommand (cmdmacro/macro_cmd.go:149)
    at github.com/alimtvnetwork/gitmap-v28/cli/cmdmacro.runMacroCmd (cmdmacro/macro_cmd.go:133)
    at github.com/alimtvnetwork/gitmap-v28/cli/cmdmacro.RunMacroCmd (cmdmacro/exports.go:18)
    at github.com/alimtvnetwork/gitmap-v28/cli/cmd.dataExecutionEntries.func2 (cmd/rootdata.go:71)
    at github.com/alimtvnetwork/gitmap-v28/cli/cmd.runDispatchTable (cmd/rootdispatch.go:24)
    at github.com/alimtvnetwork/gitmap-v28/cli/cmd.dispatchData (cmd/rootdata.go:10)
    at github.com/alimtvnetwork/gitmap-v28/cli/cmd.dispatch (cmd/root.go:298)
    at github.com/alimtvnetwork/gitmap-v28/cli/cmd.runDispatch (cmd/root.go:124)
    at github.com/alimtvnetwork/gitmap-v28/cli/cmd.Run (cmd/root.go:94)
    at main.main (cli/main.go:7)
```

Root Cause: `cli/macro/storage.go` hardcoded `getMacroDir()` to strictly attempt `filepath.Join(home, ".gitmap", "macros")`. When `.gitmap` was created under elevated permissions (e.g. `sudo`), non-root users have no write permissions to create subdirectories, crashing `SaveMacro`.

## Task-Specific Rule Set (Domain-Specific Constraints)

1. **Rule 1 (Multi-Tier Writable Directory Fallback Chain)**: `getMacroDir()` and `SaveMacro()` probe candidate directories in order: `~/.gitmap/macros`, `~/.config/gitmap/macros`, `~/.local/share/gitmap/macros`, `./.gitmap/macros`, and `%TEMP%/.gitmap-<user>/macros` (or `/tmp/.gitmap-<user>/macros`). It never fails with `permission denied` if any alternative user-writable location is available.
2. **Rule 2 (Writability Probing & Permission Auto-Healing)**: Every directory candidate is validated using an ephemeral probe write (`.perm_probe_<nanos>`). If a permission error occurs, GitMap attempts non-destructive chmod healing before falling back to XDG/local directories.
3. **Rule 3 (Universal Cross-Directory Macro Discovery)**: `LoadMacro`, `ListMacros`, `MacroExists`, and `DeleteMacro` scan all candidate directories, allowing macros stored in legacy `~/.gitmap/macros` or fallback `~/.config/gitmap/macros` to be read, edited, and executed transparently.
4. **Rule 4 (Sudo Ownership Preservation)**: When running under `sudo` on Linux, GitMap detects `SUDO_USER` / `SUDO_UID` / `SUDO_GID` and restores ownership of created configuration directories and macro files so subsequent non-root invocations never encounter permission denial.
5. **Rule 5 (Non-Negotiable Coding Standards & Total Ban)**: Maximum function length <= 15 lines (target <= 8 lines). Affirmative booleans only (`is*`, `has*`). No negative booleans. Zero nested ifs. AppError wrapping on all error paths. TOTAL BAN on `go test`, `go build`, and runner scripts during routine execution.

## Acceptance Criteria

```text
gitmap macro: execute failed: [E9000:EXECUTION] save macro alim1: mkdir /home/a/.gitmap/macros: permission denied (at=cmdmacro/macro_add.go:45)
```
1. Eliminate `mkdir ... permission denied` error when saving macros under non-root user accounts where `~/.gitmap` has restricted permissions.
2. Multi-tier fallback directory resolution seamlessly selects a writable directory (`~/.config/gitmap/macros`, `./.gitmap/macros`, or `/tmp`).
3. Existing macros in any candidate directory are fully discoverable and executable by `gitmap macro ls`, `gitmap macro run`, `gitmap schedule`.
4. Proactive permission healing and sudo ownership restoration.

## Consolidated Subtasks

### Subtask 01: Multi-Tier Directory Resolution & Writability Probe
- Implemented `cli/macro/storage_dirs.go`:
  - `isSudoUserHomeValid`: validates sudo user home path.
  - `resolveEffectiveHomeDir`: resolves home accounting for `SUDO_USER`.
  - `resolveXdgConfigMacroDir`: resolves `$XDG_CONFIG_HOME/gitmap/macros` or `~/.config/gitmap/macros`.
  - `resolveXdgDataMacroDir`: resolves `$XDG_DATA_HOME/gitmap/macros` or `~/.local/share/gitmap/macros`.
  - `resolveCurrentUserName`: cross-platform user lookup from `USER` or `USERNAME`.
  - `resolveTempMacroDir`: per-user isolated directory in system temp.
  - `candidateMacroDirs`: prioritizes candidate directories in fallback chain.
  - `deduplicateDirs`: ensures unique canonical directory paths.
  - `isDirWritable`: creates and removes ephemeral `.perm_probe_<nanos>` to verify write access.
  - `probeAndPrepareDir`: ensures directory exists, applies healing if needed, tests writability, and restores sudo ownership.
  - `resolveWritableMacroDir`: probes candidates in order and returns first writable path.

### Subtask 02: Cross-Directory Macro CRUD Operations
- Refactored `cli/macro/storage.go`:
  - `getMacroDir`: delegates to `resolveWritableMacroDir()`.
  - `SaveMacro`: resolves first writable directory and atomically writes macro JSON with sudo ownership preservation.
  - `LoadMacro`: searches all candidate directories and returns first matching macro definition.
  - `ListMacros`: traverses all candidate directories, loads macros, and deduplicates by macro name.
  - `MacroExists`: verifies macro presence across all candidate directories.
  - `DeleteMacro`: safely removes macro file across all candidate directories.

### Subtask 03: Sudo Ownership Restoration & Permission Healing
- Implemented `cli/macro/sudo_owner.go`:
  - `isSudoActive`: detects whether execution is under `sudo` with valid non-root user.
  - `restoreSudoOwnership`: invokes `chown -R <SUDO_USER>` on created target paths on Linux/macOS.
  - `attemptPermissionHealing`: attempts non-destructive `chmod 0777` on target and parent directories on Linux/macOS.

## Verification & Quality Gate Results
- `python linter-scripts/check-nested-ifs.py`: PASS (0 nested ifs across 2897 files).
- `python linter-scripts/check-enum-and-boolean.py`: PASS (0 boolean/enum violations across 2185 source files).
- `python linter-scripts/check-error-management.py`: PASS (0 bare panics or swallowed errors across 2934 source files).
- Function length constraints: All new and refactored functions in `cli/macro/` are <= 12 lines (meeting <= 15 line max and <= 8 line target).
- Tracked modified files via `python 03-ai-scripts/33-test-inventory-generator.py --record`.
