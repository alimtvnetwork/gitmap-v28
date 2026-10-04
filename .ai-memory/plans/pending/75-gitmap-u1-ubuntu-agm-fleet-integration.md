# Plan 75: GitMap U1 Ubuntu Fleet Integration, AGM Migration & Cross-OS Automation

## Status: Completed
- **Plan ID:** 75
- **Spec Reference:** [02-spec/21-app/205-gitmap-u1-ubuntu-agm-fleet-integration/01-architecture-spec.md](../../02-spec/21-app/205-gitmap-u1-ubuntu-agm-fleet-integration/01-architecture-spec.md)
- **Scope:** Cross-Platform Fleet Management, Terminal UI Polish, Cross-OS Toolchain, AGM Linux Migration
- **Created At:** 2026-10-04
- **Parent Goal:** Connect to Ubuntu VM node `U1` (`/home/a/git-work/`), probe and cache system telemetry, deploy cross-OS PowerShell installer, polish GitMap scan deduplication and vibrant terminal colors, migrate AGM settings from Windows, diagnose and resolve Linux account switching, and produce comprehensive documentation and runbooks.

---

## 1. Executive Summary

This plan orchestrates the end-to-end integration of Ubuntu node `U1` with GitMap, fixes terminal UI and scanning validation issues identified in user captures, delivers a cross-OS PowerShell installer, and completes the migration and troubleshooting of Antigravity Manager (AGM) account switching on Linux.

---

## 2. Subtasks Breakdown

### Subtask 75.1: U1 SSH Node Discovery, Telemetry Probe & Local Caching
- **Spec Reference:** Spec 205 §3.1
- **File:** `cli/cmdnode/node_telemetry.go`, `cli/store/node_cache.go`
- **Actions:**
  - Query node `U1` via SSH for OS distribution, version, kernel, architecture, and installed GitMap version.
  - Confirm default workspace `/home/a/git-work`.
  - Cache telemetry in local GitMap store so subsequent commands instantly detect environment without redundant network roundtrips.

### Subtask 75.2: Cross-OS PowerShell & Toolchain Provisioning in GitMap
- **Spec Reference:** Spec 205 §3.2
- **File:** `cli/cmdinstall/install_pwsh.go`, `cli/crossplatform/shell.go`
- **Actions:**
  - Audit missing toolchains on `U1` (PowerShell is currently missing).
  - Implement/verify `gitmap install pwsh` with dual support: APT repo registration and portable tarball fallback to `~/.local/bin/pwsh`.
  - Verify cross-OS PowerShell execution via `gitmap ps` on both Windows and Ubuntu.

### Subtask 75.3: Scan Validation Warning Deduplication & Vibrant Terminal Palette
- **Spec Reference:** Spec 205 §3.3
- **File:** `cli/cmdscanner/scanner.go`, `cli/clicolors/colors.go`, `cli/cmdvisualizer/tree.go`, `cli/cmdprojectmanager/projectmanager.go`
- **Actions:**
  - Deduplicate scan validation warnings so missing clone URL notices (e.g. `.oh-my-zsh`) are outputted only once instead of repeating across CSV/JSON export loops.
  - Clean duplicate phrasing in VS Code Project Manager error message.
  - Upgrade terminal tree colors from dull ANSI 32/33 to High-Intensity Bright Bold Green (`\033[92;1m`) and Vibrant Gold Yellow (`\033[93;1m`).

### Subtask 75.4: U1 Codebase Pull, Build & Bi-Directional Communication
- **Spec Reference:** Spec 205 §1.1, §2
- **File:** `cli/cmdssh/ssh_exec.go`, remote build scripts
- **Actions:**
  - Pull updated `gitmap` repository on `U1` at `/home/a/git-work/gitmap`.
  - Build native Linux binary on `U1` and install to `/usr/local/bin/gitmap` or `~/.local/bin/gitmap`.
  - Verify terminal visualizer and command execution on `U1`.

### Subtask 75.5: Antigravity Manager (AGM) Linux Account Switching RCA & Migration
- **Spec Reference:** Spec 205 §3.4
- **File:** `scripts/agm/export-import.sh`, `scripts/agm/account-switch-helper.py`
- **Actions:**
  - Inspect AGM codebase on `U1` at `/home/a/git-work/Antigravity-Manager`.
  - Analyze root cause for Linux account switching failure (DBus Secret Service dependency, headless session keyring locks, Chrome singleton sockets).
  - Export Windows AGM profiles and settings via GitMap, transfer via SCP, and import on `U1`.
  - Validate profile switching and establish headless credential fallback.

### Subtask 75.6: Audit Trail, Automated Runbook & Upstream AI Instruction Prompt
- **Spec Reference:** Spec 205 §4
- **File:** `docs/migration-windows-to-ubuntu-agm.md`, `scripts/run-migration.sh`, `docs/agm-linux-account-switching-fix.prompt.md`
- **Actions:**
  - Compile exhaustive MD migration report with rationale, commands run, and decision trail into the repo docs/cache.
  - Generate standalone executable `run-migration.sh` script for zero-touch VM provisioning.
  - Create standalone AI prompt for upstream AGM developers to implement native Linux headless keyring support.
