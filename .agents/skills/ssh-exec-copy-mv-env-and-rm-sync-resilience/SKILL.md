---
name: ssh-exec-copy-mv-env-and-rm-sync-resilience
description: Autonomously implement and verify SSH parallel execution with UI padding, command injection banners, offline node detection, --except filtering, SSH copy and move with path macro expansion (~, %win%, %temp%), cross-platform environment variable management (env add/rm/ls), pipeline errors clear alias, and resilient rm/sync operations across GitMap.
---

# SSH Execution, Remote Copy/Move, Environment Management & CLI Resilience Suite

## Overview
This skill provides comprehensive instructions for enhancing GitMap's SSH remote orchestration suite (`gitmap ssh exec`, `gitmap ssh copy`, `gitmap ssh mv`), cross-platform environment variable management (`gitmap env add/rm/ls`), pipeline error clearing, resilient untracking/removal (`gitmap rm`), and sync commands.

## Key Capabilities & Requirements

### 1. SSH Execution & UI Polish (`gitmap ssh exec`)
- **Parallel Multi-Node Execution**: Run commands across all nodes in parallel with clean visual padding (top, bottom, left).
- **Command Injection Banner**: Display clear summary of commands being injected and target machines/IPs upfront without credential leakage (suppress noisy "No password or key configured" messages).
- **Offline Node Pre-Detection & Summary**:
  - Detect offline/unreachable machines at start and announce them upfront.
  - Summarize offline machines at completion.
- **Command List & Filter Syntax**:
  - Support comma-separated or space-separated command sequences (`cmd1,cmd2,cmd3`).
  - Support `--except <machine-name|alias|id>` to skip specific nodes.
- **Progress & Result Reporting**: Display dynamic processing state and print structured results per machine upon completion.

### 2. SSH Remote Copy & Move (`gitmap ssh copy`, `gitmap ssh mv`)
- **Copy Engine (`gitmap ssh copy <from> <to>`)**:
  - Copy from host machine to remote default working directory across all nodes.
  - Support `--except <alias|name|id>` filter.
  - Detailed help documentation and interactive examples (`gitmap ssh copy help`).
- **Move Engine (`gitmap ssh mv <from> <to>`)**:
  - Move/transfer files from host machine to target remote locations across nodes with `--except` support.
- **Path Macro Expansion**:
  - Expand `~` (user home directory).
  - Expand `%win%` (Windows system folder: `C:\Windows` or equivalent).
  - Expand `%win-drive%` (Windows system drive, e.g. `C:`).
  - Expand `%temp%` / `%appdata%` / `$ENV_VAR` / `%VAR%` across Linux and Windows.

### 3. Cross-Platform Environment Management (`gitmap env add/rm/ls`)
- Export and persist environment variables cross-platform (Ubuntu Linux and Windows).
- Subcommands: `gitmap env add <key> <val>`, `gitmap env rm <key>`, `gitmap env ls`, `gitmap env help`.

### 4. CLI Resilience & Error Handling
- **Pipeline Errors Clear**: Ensure `gitmap pipeline errors clear` works as a top-level alias alongside `storage reset-errors`.
- **Resilient Package Removal (`gitmap rm`)**:
  - If a file/folder is locked by another Windows process (`unlinkat ... The process cannot access the file`), retry or report clean error without dumping internal raw AppError stack traces.
  - If a package/path is missing on disk, cleanly untrack it from GitMap database rather than aborting fatally.
- **Sync & Package Database Clear**:
  - Ensure `gitmap sync` commands are fully documented and functional.
  - Provide clear guidance to reset or clear package databases.

### 5. Architectural & Coding Guidelines
- Functions $\le$ 8–15 lines.
- Affirmative booleans (`is*`, `has*`).
- Single return types (zero error tuples).
- Universal `*apperror.AppError` envelopes.
- Strict Unix LF line endings.
