# Plan 199: SSH Fleet Package Installation Delegation & Function Reduction Audit

## Task Lifecycle & Completion Summary
- **Initiated**: User requested comprehensive audit and completion of all missing SSH elements, including remote package installation (`gitmap ssh install <pkg> [target]`), fleet updates, AppError diagnostics, and function size compliance.
- **Status**: COMPLETED (100% Verified)

---

## 1. Overview & Problem Statement
Following the terminal SSH execution improvements, the user requested an audit of all missing elements:
1. Fleet package installation and updating (`gitmap ssh install` / `gitmap ssh update`):
   - Ability to install/update GitMap or any supported package (e.g. `agy`, `devbox`) across remote machines.
   - Fault tolerance: offline machines are skipped with clear notice (`[alias|ip] OFFLINE (skipped: reason)`), while online machines execute concurrently.
2. Coding Guidelines & Canonical Sizing:
   - All Go files in `cli/cmdssh/` must remain strictly $\le 100$ lines.
   - All Go functions must remain strictly $\le 15$ lines.
   - Zero nested ifs (nesting depth $\le 1$), zero boolean inversions, zero anti-ok variables.
   - Elimination of `go vet` warnings (e.g. redundant newlines in `fmt.Println`).
3. Unit test coverage for install arguments parsing, package resolution, and remote shell detection.

---

## 2. Changes Made & Consolidated Subtasks

### Subtask 01: Fleet Remote Package Installation & Updates
- **`cli/cmdssh/ssh_install_remote.go`** (93 lines):
  - `runSSHInstallCLI`: entry point for fleet-wide or targeted installs.
  - `executeFleetInstall`: iterates over target machines and executes installs.
  - `parseInstallTargetAndPackage` & `parseSingleArgInstall`: flexibly parses 0, 1, or 2 positional arguments.
  - `installOrUpdateNodeWithPackage`: conducts pre-flight TCP reachability check before executing actions.
  - `dispatchNodePackageAction`: handles gitmap presence verification, updates existing installations, or installs specified packages.
- **`cli/cmdssh/ssh_install_actions.go`** (59 lines):
  - `installFreshGitmap`: generates official cross-platform bootstrap script and executes remotely.
  - `updateExistingGitmap`: executes `gitmap update` remotely with fallback reinstall.
  - `installRemotePackage`: invokes `gitmap install <package>` remotely on target machines.
  - `isNodeAvailable`: pre-flight TCP reachability probe.
  - `resolveRemoteShell`: selects `ps` for Windows and `bash` for Linux/macOS.
  - `reportRemoteExecution`: standardizes success and error outputs.
- **`cli/cmdssh/ssh_update_remote.go`** (58 lines):
  - `runSSHUpdateCLI`: entry point for fleet-wide updates.
  - `executeFleetUpdate`: iterates over filtered connections.
  - `updateSingleSSHNode`: checks node reachability before triggering update.
  - `executeRemoteUpdate`: executes `gitmap update` or `gitmap update <package>`.

### Subtask 02: Antigravity & VS Code Remote Delegation Modularization
- **`cli/cmdssh/ssh_target_nodes.go`** (33 lines):
  - Centralized `loadTargetNodes` and `checkRemoteNodeOnline` helpers shared between `agy` and `code` subcommands.
- **`cli/cmdssh/ssh_agy_cmd.go`** (91 lines):
  - Decomposed into `runSSHAgyCLI`, `executeAgyOnFleet`, `showSSHAgyUsage`, `parseRemoteTargetAndArgs`, `runAgyOnNode`, `establishAgyClient`, and `executeAgyRemoteCommand`.
  - All functions reduced to $\le 14$ lines.
- **`cli/cmdssh/ssh_code_cmd.go`** (73 lines):
  - Decomposed into `runSSHCodeCLI`, `executeCodeOnFleet`, `showSSHCodeUsage`, `runCodeOnNode`, `resolveTargetCodePath`, and `tryLaunchLocalCodeRemote`.
  - All functions reduced to $\le 14$ lines.
- **`cli/cmdssh/ssh_code_remote.go`** (37 lines):
  - Extracted `runRemoteCodeBinary` and `executeRemoteCodeSession` for fallback SSH execution.
- **`cli/cmdssh/ssh_compare_table.go`** (83 lines):
  - Fixed redundant newline in `fmt.Println` to achieve 0 `go vet` warnings.

### Subtask 03: Unit Test Coverage
- **`cli/cmdssh/ssh_install_test.go`** (54 lines):
  - `TestParseInstallTargetAndPackageZeroArgs`: validates default `gitmap` and `all`.
  - `TestParseInstallTargetAndPackageSingleArgGitmap`: validates explicit `gitmap`.
  - `TestParseInstallTargetAndPackageSingleArgTarget`: validates known target detection.
  - `TestParseInstallTargetAndPackageSingleArgPackage`: validates package target defaulting to `all`.
  - `TestParseInstallTargetAndPackageTwoArgs`: validates package + target syntax.
  - `TestResolveRemoteShell`: validates Windows PowerShell vs Linux Bash resolution.

---

## 3. Verification & Compliance Matrix

| Quality Gate | Tool / Command | Result |
|---|---|---|
| **Go File Sizing** | Canonical limit: $\le 100$ lines for every Go file in `cli/cmdssh/` | **100% PASS** (all files $\le 93$ lines) |
| **Function Length** | Canonical limit: $\le 15$ lines per function | **100% PASS** (all functions $\le 14$ lines) |
| **Go Vet** | `go vet ./cmdssh` | **100% PASS** (0 warnings) |
| **Unit Tests** | `go test -v ./cmdssh` | **100% PASS** (all unit tests passing) |
| **Control Flow (Nested Ifs)** | `python linter-scripts/check-nested-ifs.py` | **100% PASS** (0 violations across 3,098 files) |
| **Booleans & Enums** | `python linter-scripts/check-enum-and-boolean.py` | **100% PASS** (0 violations across 2,334 files) |
| **Go Compilation Gate** | `python 03-ai-scripts/06-cicd-local-runner.py --filter "Compile" --no-cache` | **100% PASS** (compiled in 3.25s) |
