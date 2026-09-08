# 01-installer-research: Installer Package Resolution & Command Handling for Nginx, WordPress, and Laravel

**Date:** 2026-09-09  
**Status:** Completed  
**Audited Codebase:** `gitmap/cmd/install*.go`, `gitmap/constants/constants_install.go`, `gitmap/store/`  
**Target Capabilities:** Add Nginx, WordPress, and Laravel as first-class install targets across Windows (`choco` / `winget`), Linux (`apt`), and macOS (`brew`).

---

## 1. Executive Summary & Objective

GitMap provides a unified cross-platform tool installer CLI accessible via `gitmap install <tool>` (alias `gitmap in <tool>`). The subsystem manages:
- Dynamic package manager resolution (`choco`, `winget`, `apt`, `brew`, `snap`).
- Tool name validation, canonicalization, and alias expansion.
- Pre-installation existence checks via binary inspection and registry queries.
- Audited command generation and dry-run simulation.
- Post-install binary verification and version extraction.
- Recording install state in `installation.db` (`InstalledTool` and `InstallationLog`).
- Interactive grouped tool listing (`gitmap install --list`) with real-time status dots.
- Safe tool uninstallation (`gitmap uninstall <tool>`).

This research details the structural requirements, package identifiers, execution mechanics, binary verification quirks, and implementation roadmap to onboard **Nginx**, **WordPress**, and **Laravel** as first-class citizens in GitMap.

---

## 2. Deep Dive: GitMap Installer Subsystem Architecture

The installer pipeline operates across a 10-phase lifecycle:

```
[CLI Args]
    │
    ▼
1. Flag Parsing & Canonicalization (cmd/install.go)
    ├─ parseInstallFlags: --manager, --version, --verbose, --dry-run, --check, --list, --yes
    ├─ resolveToolAlias: maps aliases to canonical tool constants
    └─ validateToolName: checks InstallToolDescriptions (exits E9000 if invalid)
    │
    ▼
2. Handler Routing (cmd/install.go -> cmd/install_handlers.go)
    ├─ specialInstallHandler: custom pipelines (vscode-linux, chrome-linux, build-essential)
    └─ executeGenericInstall (cmd/install_prompt.go)
    │
    ▼
3. Pre-Installation Existence Detection (cmd/install_prompt.go & cmd/installverify.go)
    ├─ alreadyInstalled: prints MsgInstallChecking
    ├─ detectInstalledVersion: checks expectedExePath -> LookPath -> getInstalledVersion
    └─ if found: skips installation; if --check: reports status and returns
    │
    ▼
4. Package Manager Resolution (cmd/installdetect.go)
    ├─ resolvePackageManager: honors CLI flag override (--manager)
    └─ detectPackageManager:
         Windows: checks LookPath("choco") -> LookPath("winget") -> default choco
         macOS: LookPath("brew") -> default brew
         Linux: CurrentOS() (Ubuntu/Debian -> apt, Fedora/CentOS -> dnf) -> LookPath fallback
    │
    ▼
5. Package Identifier Mapping (cmd/install_packages.go & cmd/install_packages_extra.go)
    ├─ resolvePackageName(manager, tool):
         choco  -> chocoPackageMap[tool]
         winget -> wingetPackageMap[tool]
         apt    -> aptPackageMap[tool]
         brew   -> brewPackageMap[tool]
         snap   -> snapPackageMap[tool]
    └─ fallback: returns tool name unchanged
    │
    ▼
6. Command Construction & Execution (cmd/installtools.go & cmd/install_audit_runner.go)
    ├─ buildInstallCommand:
         choco:  choco install <pkg> -y --no-progress [--version <ver>]
         winget: winget install <pkg> --accept-package-agreements --accept-source-agreements --silent [--version <ver>]
         apt:    sudo apt install -y <pkg>[=<ver>] (with pre-step sudo apt-get update)
         brew:   brew install [--cask] <pkg>
    ├─ handleDryRunInstall: prints MsgInstallDryCmd and halts if --dry-run
    └─ runInstallCommand: executeCommandWithAudit -> records duration, exit code, stdout, stderr
    │
    ▼
7. Verification & Version Discovery (cmd/installverify.go)
    ├─ verifyInstallation:
         1. Checks expectedExePath(tool) on Windows
         2. Checks isGUITool(tool) (skips CLI execution to avoid blocking GUI popups)
         3. Resolves toolBinaryName(tool)
         4. Executes getInstalledVersion(binary)
    └─ runPostInstall(tool): tool-specific hooks (e.g., git lfs install, core.longpaths)
    │
    ▼
8. State Persistence & Audit Logging (cmd/installtools.go & store/)
    ├─ splitDB.SaveInstalledTool(tool, version, manager) -> SQLite table InstalledTool
    └─ splitDB.RecordLog(...) & splitDB.RecordExecution(...) -> SQLite table InstallationLog
    │
    ▼
9. Catalog & Grouped Listing (cmd/installlist.go & constants/constants_install.go)
    ├─ InstallToolDescriptions: map[string]string (tool descriptions)
    ├─ InstallToolCategories: map[string][]string (grouped by category)
    └─ printInstallListGrouped: renders status dots (● installed, ○ missing, ? unknown)
    │
    ▼
10. Uninstallation Flow (cmd/uninstall.go)
    ├─ resolveUninstallManager: reads manager from InstalledTool record
    ├─ buildUninstallCommand: choco uninstall, winget uninstall, apt remove/purge, brew uninstall
    └─ removes database record upon successful execution
```

---

## 3. Analysis of Target 1: Nginx

### 3.1 Characteristics
- **Type:** High-performance HTTP server, reverse proxy, and edge gateway daemon.
- **Category:** `ToolCategoryDevOps` ("DevOps & Containers")
- **Tool Constant:** `constants.ToolNginx = "nginx"`
- **Tool Description:** `"Nginx high-performance HTTP server and reverse proxy"`
- **Aliases:** `"ngx"`, `"engine-x"`

### 3.2 Package IDs Across Package Managers
| Platform / Manager | Package ID Constant | Upstream Package ID | Command Syntax |
| :--- | :--- | :--- | :--- |
| **Windows (Chocolatey)** | `constants.ChocoPkgNginx` | `nginx` | `choco install nginx -y --no-progress` |
| **Windows (Winget)** | `constants.WingetPkgNginx` | `nginxinc.nginx` | `winget install nginxinc.nginx --accept-package-agreements --accept-source-agreements --silent` |
| **Linux (APT)** | `constants.AptPkgNginx` | `nginx` | `sudo apt install -y nginx` |
| **macOS (Homebrew)** | `constants.BrewPkgNginx` | `nginx` | `brew install nginx` (Formula) |

### 3.3 Binary & Verification Quirks
1. **Binary Name:** `nginx` (`nginx.exe` on Windows).
2. **Standard Output vs Standard Error Trap:**
   - Standard `nginx -v` writes version string exclusively to `stderr` (e.g., `nginx version: nginx/1.24.0`).
   - GitMap's `getInstalledVersion` must use `CombinedOutput()` and pass `-v` for `nginx`.
3. **Windows Installation Path:**
   - Chocolatey installs Nginx to `C:\tools\nginx\nginx.exe`.

---

## 4. Analysis of Target 2: WordPress

### 4.1 Characteristics
- **Type:** Content Management System & Developer CLI (`wp-cli`).
- **Category:** `ToolCategoryCore`.
- **Tool Constant:** `constants.ToolWordPress = "wordpress"`
- **Tool Description:** `"WordPress web publishing platform and CMS / WP-CLI"`
- **Aliases:** `"wp"`, `"wp-cli"`, `"wpcli"`

### 4.2 Package IDs Across Package Managers
| Platform / Manager | Package ID Constant | Upstream Package ID | Command Syntax |
| :--- | :--- | :--- | :--- |
| **Windows (Chocolatey)** | `constants.ChocoPkgWordPress` | `wordpress` | `choco install wordpress -y --no-progress` |
| **Windows (Winget)** | `constants.WingetPkgWordPress` | `Automattic.Wordpress` | `winget install Automattic.Wordpress --accept-package-agreements --accept-source-agreements --silent` |
| **Linux (APT)** | `constants.AptPkgWordPress` | `wordpress` | `sudo apt install -y wordpress` |
| **macOS (Homebrew)** | `constants.BrewPkgWordPress` | `wp-cli` | `brew install wp-cli` |

### 4.3 Binary & Verification Mechanics
- Binary name: `wp` (with GUI fallback detection).

---

## 5. Analysis of Target 3: Laravel

### 5.1 Characteristics
- **Type:** Modern PHP Web Application Framework and CLI Installer.
- **Category:** `ToolCategoryLanguages`.
- **Tool Constant:** `constants.ToolLaravel = "laravel"`
- **Tool Description:** `"Laravel PHP web application framework and installer CLI"`
- **Aliases:** `"artisan"`, `"laravel-installer"`

### 5.2 Package IDs Across Package Managers
| Platform / Manager | Package ID Constant | Upstream Identifier | Command Syntax |
| :--- | :--- | :--- | :--- |
| **Windows (Chocolatey)** | `constants.ChocoPkgLaravel` | `laravel` | `choco install laravel -y --no-progress` |
| **Windows (Winget)** | `constants.WingetPkgLaravel` | `Laravel.Laravel` | `winget install Laravel.Laravel --accept-package-agreements --accept-source-agreements --silent` |
| **Linux (APT)** | `constants.AptPkgLaravel` | `laravel` | `sudo apt install -y laravel` |
| **macOS (Homebrew)** | `constants.BrewPkgLaravel` | `laravel` | `brew install laravel` |
