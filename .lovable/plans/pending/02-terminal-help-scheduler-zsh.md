# Parent Task: Terminal Help Alignment, Scheduler CLI, and ZSH/OS Startup Hooks

## Overview

This specification details the implementation of a comprehensive update to the GitMap CLI help aesthetics (alignment, colored flags, headers spacing), adding missing help sections (SSH, Macro, etc.), building a new cross-platform Scheduler CLI via SQLite, implementing an OS-level startup mechanism for macros, and expanding Ubuntu setups with ZSH integration.

## Architectural Goals & Coding Guidelines

- **Booleans**: Must use `is`, `has`, `can`, `should`. No negatives (`isNotReady` banned).
- **Naming**: Strict semantic naming. No `temp`, `data`, `obj`.
- **Functions**: Max 15 lines. Arguments wrapped at 100 characters.
- **Error Handling**: Wrap all errors with `AppError`.
- **Cross-Platform**: Support Windows (PowerShell/CMD/Registry) and Linux (Bash/ZSH/systemd/cron).

## Subtasks Execution Plan

### Task 1: CLI Help Re-Alignment & Colors

- **File**: `gitmap/cmd/rootusage*.go`, `gitmap/constants/constants_help*.go`
- **Features**:
  - Add line gaps between major Help headers (e.g. `\n\n` before sections).
  - Add distinct terminal colors (e.g., cyan/green/yellow) to the flags (like `--refresh`, `--config`) to make them stand out.
  - Fix alignment in the output sections (like `clone-fix-repo-pub`, `interactive`, etc.) to match perfectly.
  - Inject missing help definitions for: SSH, SSH Join, Coding Guideline, Installer, Macro.

### Task 2: Scheduler CLI Core & SQLite Engine

- **File**: `gitmap/cmd/schedule_cmd.go`, `gitmap/store/scheduler.go`
- **Features**:
  - `gitmap schedule <taskname>`: Launches interactive mode to define the command (bash/pwsh), the macro, or the script.
  - `--delay`, `--interval`: Support per second, minute, hour, day, week, month.
  - `gitmap schedule status`: Lists all active/running schedulers from SQLite (`scheduler_tasks` table).
  - Add SQLite schema for `scheduler_tasks`.

### Task 3: OS Startup Macro Hook & OS Management

- **File**: `gitmap/cmd/os_startup.go`, `gitmap/osutil/startup_*.go`
- **Features**:
  - Add OS startup capability so macros/scheduled tasks can boot automatically on login.
  - Windows: Write to `HKCU\Software\Microsoft\Windows\CurrentVersion\Run`.
  - Linux: Write `~/.config/autostart/gitmap-macro.desktop` or `systemd` user service.
  - Support `gitmap schedule restart` / `shutdown` commands using native OS calls (`shutdown /r /t`, `systemctl reboot`).

### Task 4: Ubuntu ZSH Installer & Setup Hook

- **File**: `gitmap/cmd/setup_ubuntu.go` (or `setup.go`)
- **Features**:
  - Add a dedicated routine to install ZSH, Oh-My-Zsh (or similar config), and theme switching for Ubuntu environments.
  - Ensure users can opt-in during `gitmap setup` on Linux.

### Task 5: Web UI Terminal Sync Verification

- **File**: `gitmap/cmd/hd_server.go`, Web UI components
- **Features**:
  - Confirm the Localhost endpoint and terminal commands proxy perfectly into the Web UI.
  - Expose the Help string changes through the UI.

### Task 6: Release & Architecture Map

- **Files**: `version.json`, `.lovable/memory/release-architecture-map.md`
- **Features**:
  - Bump MINOR version in `version.json`.
  - Append to changelog and perform release ceremony.

## Verification Checklist (Pre-Commit)

- [ ] Coding Guidelines & Master Consolidated File enforced.
- [ ] Boolean Examples & Fixations strictly followed.
- [ ] Anti-Garbage Naming enforced (no generic temp names).
- [ ] Semantic Tests.
- [ ] Function Size <= 15 lines.
- [ ] Error Handling uses wrappers.
- [ ] Code adheres to explicit booleans, Type-suffixed Enums.
- [ ] Formatting & Acronyms strictly PascalCase (e.g., `SwapIpWindows`).
- [ ] Temp-Scripts ignored.
