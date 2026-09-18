# Plan 155: Antigravity Cross-Platform Installer Suite

- **Slug**: `155-antigravity-crossplatform-installer`
- **Date**: 2026-09-14
- **Status**: completed
- **Execution Budget**: N = 150 loops
- **Loops Taken**: 2 self-loop iterations (Phase 1: Planning, Asset Analysis & Subtask Decomposition; Phase 2: Implementation, Code Generation, Guidelines Compliance & Consolidation)

---

## 1. Problem Statement & Objectives

The user requested reading, analyzing, and implementing the Google Antigravity IDE and CLI installer specification located at `D:\work\antigravity-installer-v1\01-installer` (`README.md`, `agy-install.ps1`, `agy-install.sh`) as a proper, robust, cross-platform installation engine for Ubuntu/Linux, Windows, and macOS, categorically detecting and announcing the host OS and architecture before installation.

### Key Pain Points Resolved:
1. **Categorical OS Announcement**: The installer now clearly detects and announces the target OS, distribution details, and architecture before proceeding (e.g. `[INFO] Detected Platform: Linux (Ubuntu 22.04.4 LTS) [arch: x64]`).
2. **Canonical GCS Download Matrix**: Dynamic URL derivation `{BASE_URL}/{VERSION}-{BUILD_ID}/{PLATFORM}/{ARTIFACT}` supporting Windows (`windows-x64`, `windows-arm`), macOS (`darwin-arm`, `darwin-x64`), and Linux (`linux-x64`, `linux-arm`).
3. **Resilient Retry-Backed Download**: 3 download attempts with exponential backoff and progress logging.
4. **Prerequisite Validation**: Linux glibc >= 2.28 and toolchain (`tar`, `curl`); Windows build >= 17763; macOS version >= 12.0.
5. **Cross-Platform Deployment**:
   - **Linux**: Native tarball unpacking with retry, `0755` executable permissions, SUID on `chrome-sandbox` (4755) when root, dual symlinks (`antigravity` & `agy` in `/usr/local/bin` or `~/.local/bin`), `.desktop` launcher creation with icon discovery and `update-desktop-database` cache refresh.
   - **Windows**: Silent NSIS `/S` installer execution (with fallback to interactive), dual CMD shims (`antigravity.cmd` & `agy.cmd`) and PowerShell shims (`antigravity.ps1` & `agy.ps1`) in `%LOCALAPPDATA%\agy\bin`, and user environment PATH injection.
   - **macOS**: DMG mounting via `hdiutil attach -nobrowse -quiet`, `.app` bundle copy via `ditto`, guaranteed detach in defer/finally, and quarantine attribute removal via `xattr -dr com.apple.quarantine`.
   - **Privileged to Unprivileged Fallback**: Unprivileged runs fall back cleanly to user-local paths (`~/.local/share/antigravity`, `~/Applications`, `%LOCALAPPDATA%\Programs\Antigravity`).
6. **Dual Database Registration & Universal Uninstall**:
   - Both `"antigravity"` and `"agy"` recorded in `installation.db` and synchronized to `gitmap.db`.
   - `gitmap uninstall antigravity` and `gitmap uninstall agy` cleanly sweep binaries, shims, desktop launchers, and purge both database tables.

---

## 2. Task-Specific Rule Set (Non-Negotiable Constraints)

1. **Categorical OS Announcement**: The installer MUST print the detected platform, architecture, artifact, and download URL before starting downloads.
2. **Pure-Go Primary Engine**: All core operations (detection, download, extraction, shims, desktop integration, DB recording) MUST be implemented natively in Go without external shell dependencies, returning structured `*apperror.AppError`.
3. **Function Sizing**: Every Go function MUST be <= 15 lines (target <= 8 lines). Decompose complex operations into small helper functions.
4. **Affirmative Booleans Only**: Booleans must use `is*` or `has*` only (`isFound`, `isSupported`, `isDryRun`, `isForce`, `isPrivileged`, `hasIcon`). Zero negative variables.
5. **Multi-Database & Multi-Binary Parity**: Both `antigravity` and `agy` MUST be linked into PATH, recorded in `installation.db` and `gitmap.db`, and uninstalled cleanly by either command name.

---

## 3. Subsystem Architecture

### 3.1 Platform Detection & Prerequisite Validation
- `cli/cmdinstall/installantigravity_types.go`: Platform models (`AntigravityPlatformInfo`), version constants, and error codes.
- `cli/cmdinstall/installantigravity_detect.go`: OS detection, architecture normalization, Linux distro identification from `/etc/os-release`, and categorical announcement printing.
- `cli/cmdinstall/installantigravity_prereqs.go`: Platform prerequisite checks (glibc >= 2.28, Windows build >= 17763, macOS >= 12.0).

### 3.2 Canonical URL Resolution & Resilient Downloader
- `cli/cmdinstall/installantigravity_fetch.go`: Dynamic URL builder `{BASE_URL}/{VERSION}-{BUILD_ID}/{PLATFORM}/{ARTIFACT}` and 3-attempt resilient downloader with backoff into `tempdir.RepoTempDir("antigravity-install")`.

### 3.3 Multi-Platform Deployment
- `cli/cmdinstall/installantigravity_deploy_linux.go`: Native Linux tarball unpacking, dual symlinks, desktop launcher, and desktop database refresh.
- `cli/cmdinstall/installantigravity_deploy_windows.go`: Windows silent installer execution, verification polling, CMD/PS1 shims in `%LOCALAPPDATA%\agy\bin`, and PATH persistence.
- `cli/cmdinstall/installantigravity_deploy_darwin.go`: macOS DMG mount, bundle copy, quarantine clearance, and dual symlinks.
- `cli/cmdinstall/installantigravity_other.go`: Safe stub fallback for non-supported operating systems.

### 3.4 Database Tracking & Universal Uninstall
- `cli/cmdinstall/installantigravity_db.go`: Dual database registration in `installation.db` and `gitmap.db`.
- `cli/cmdinstall/install_uninstall_paths.go` & `cli/cmd/uninstall.go`: Uninstallation sweep for custom shims, desktop files, and dual database purge.

### 3.5 CLI Integration & Documentation
- `cli/cmdinstall/installantigravity.go`: Refactored `runInstallAntigravityWithOpts` orchestrating detection, prerequisites, idempotency, dry-run, deployment, and DB registration.
- `cli/cmdinstall/install_agy.go` & `cli/cmdinstall/agy_install.go`: Routing `gitmap install agy` directly to the unified Antigravity engine.
- `cli/helptext/antigravity.md`: Updated help text documentation.
- `cli/cmdinstall/installantigravity_test.go`: Unit tests for URL resolution, announcement formatting, prerequisite checks, and path mappings.

---

## 4. Subtasks Completed & Consolidated

- **Subtask 155-01: Core Types, Detection, Announcement & Resilient Download**
  - Created `installantigravity_types.go`, `installantigravity_detect.go`, `installantigravity_prereqs.go`.
  - Refactored `installantigravity_fetch.go` and `install.go` (added `--force` and `--prefix`).
- **Subtask 155-02: Multi-Platform Deployment (Linux, Windows, macOS)**
  - Created `installantigravity_deploy_linux.go`, `installantigravity_deploy_windows.go`, `installantigravity_deploy_darwin.go`.
  - Updated `installantigravity_windows.go` and `installantigravity_other.go`.
- **Subtask 155-03: Dual-Database Tracking & Universal Uninstall Integration**
  - Created `installantigravity_db.go`.
  - Updated `install_uninstall_paths.go` and `cli/cmd/uninstall.go`.
- **Subtask 155-04: CLI Integration, Helptext & Parity Testing**
  - Refactored `installantigravity.go`, `install_agy.go`, `agy_install.go`.
  - Updated `helptext/antigravity.md` and added unit test suite `installantigravity_test.go`.

---

## 5. Verification & Quality Gates

All checks and linters executed cleanly:
- `python linter-scripts/check-nested-ifs.py`: PASS (0 violations across 2869 files).
- `python linter-scripts/check-enum-and-boolean.py`: PASS (0 violations across 2157 source files).
- `python linter-scripts/check-boolean-guidelines.py`: PASS (0 violations across 2869 files).
- `python linter-scripts/check-error-management.py`: PASS (0 violations across 2906 source files).
- `python linter-scripts/check-function-formatting.py`: PASS (Rule 9a/9b clean across all new files).
- `python linter-scripts/check-function-signatures.py`: PASS (0 violations).
