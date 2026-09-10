# Plan 93: Google Antigravity Desktop IDE Installer Fix & CLI Decoupling

**Title:** Google Antigravity Desktop IDE Application Installer Fix & CLI Tool Decoupling  
**Status:** Completed  
**Parent Task Workflow:** Parent Task N-Step Continuous Loop & Multi-Agent Orchestration (Prompt v2.1.0)  
**Budget (N):** 130 steps  
**Target Directory:** `.lovable/plans/completed/`  

---

## 1. Root Cause & Problem Statement

1. **The Conflation**:
   - Both `scripts-fixer` (Script 69 `69-install-antigravity`) and `gitmap` (`gitmap install antigravity`) were mistakenly programmed to download and install the **Antigravity CLI** (`agy` tool from `https://api.github.com/repos/google-antigravity/antigravity-cli/releases/latest` or `https://antigravity.google/cli/install.sh`).
   - When a user requests to install **Antigravity**, they expect the full **Google Antigravity Desktop IDE / Application** (`Antigravity.tar.gz` on Linux, `Antigravity-x64.exe` on Windows), which provides the standalone AI editor, GUI canvas, and desktop application.
   - The CLI (`agy`) is a complementary command-line tool, not the core Antigravity application itself.

2. **The Fix**:
   - **`antigravity` / `antigravity-ide`**: Refactor to download, install, and verify the official **Google Antigravity Desktop IDE Application**:
     - Windows: `%LOCALAPPDATA%\Programs\Antigravity\Antigravity.exe` (installed silently via `/S` from official Google Storage `Antigravity-x64.exe`).
     - Ubuntu/Linux: `/opt/antigravity/antigravity` or `~/.local/share/antigravity/antigravity` (extracted from official Google Storage `Antigravity.tar.gz`), symlinked to `/usr/local/bin/antigravity` (or `~/.local/bin/antigravity`), with desktop launcher entry `antigravity.desktop`.
   - **`antigravity-cli` / `agy`**: Decouple into a distinct tool:
     - `scripts-fixer`: Script ID `78` (`78-install-antigravity-cli`) and `install-antigravity-cli.sh`.
     - `gitmap`: `ToolAgy = "agy"`, installable via `gitmap install agy`, `gitmap install antigravity-cli`, or `gitmap agy install cli`.

---

## 2. Task-Specific Rules & Constraints

1. **Zero Data Loss & Path Integrity**:
   - Never delete existing workspace configs or user data during install.
   - Verify existing installation before initiating downloads.
2. **Repository Coding Guidelines**:
   - Function length $\le 15$ lines.
   - Zero nested `if` statements (guard-return pattern only).
   - Affirmative boolean naming (`is*`, `has*` only, zero negative booleans like `!isFound`).
   - Universal `*apperror.AppError` wrapping with contextual operational metadata.
3. **Cross-OS Parity**:
   - Full support for Linux/Ubuntu (tarball + symlink + desktop file) and Windows (NSIS silent installer).
4. **Strict Relative Git Paths**:
   - All links in plans and subtasks use repository-relative paths.

---

## 3. Subtasks Ledger

- **Subtask 93.01**: Scripts-Fixer Antigravity Desktop & CLI Split (`scripts/69-install-antigravity/run.ps1`, `scripts/os/ubuntu/install-antigravity.sh`, `scripts/78-install-antigravity-cli/run.ps1`, `scripts/os/ubuntu/install-antigravity-cli.sh`, registry and keyword mappings).
- **Subtask 93.02**: Gitmap Constants, Tool Mappings & Binary Probes (`gitmap/constants/constants_install.go`, `gitmap/cmd/install_packages.go`, `gitmap/cmd/installprobe.go`, `gitmap/cmd/installverify.go`).
- **Subtask 93.03**: Gitmap Desktop IDE Installer & Platform Engines (`gitmap/cmd/installantigravity.go`, `gitmap/cmd/installantigravity_fetch.go`, `gitmap/cmd/installantigravity_linux.go`, `gitmap/cmd/installantigravity_windows.go`).
- **Subtask 93.04**: Gitmap CLI Installer Decoupling & AGY Command Parity (`gitmap/cmd/install_agy.go`, `gitmap/cmd/install_agy_fetch.go`, `gitmap/cmd/agy_install.go`, `gitmap/cmd/install_handlers.go`).
- **Subtask 93.05**: Unit Tests, Quality Gate Verification & CI/CD Runner (`gitmap/cmd/installantigravity_test.go`, linters, and local CI runner).
