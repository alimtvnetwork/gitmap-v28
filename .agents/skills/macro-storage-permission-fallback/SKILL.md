---
name: macro-storage-permission-fallback
description: Autonomously diagnose, repair, and verify macro storage permissions, multi-tier fallback directory resolution, sudo user ownership auto-healing, and cross-directory macro discovery across GitMap.
---

# Macro Storage Permission & Fallback Suite

## Core Capabilities

1. **Multi-Tier Writable Directory Resolution**:
   - Primary: `~/.gitmap/macros` (standard user macro store).
   - Secondary: `~/.config/gitmap/macros` (standard XDG configuration directory, user-owned).
   - Tertiary: `~/.local/share/gitmap/macros` (standard XDG data directory).
   - Quaternary: `./.gitmap/macros` (repository-local workspace fallback).
   - Quinquenary: `%TEMP%/.gitmap-<user>/macros` or `/tmp/.gitmap-<user>/macros` (isolated ephemeral fallback).

2. **Proactive Permission Healing & Auto-Chmod**:
   - Detect `permission denied` on `os.MkdirAll` and attempt permission healing via `os.Chmod` with `0777` or `0755` when allowed.
   - When running with elevated privileges under `sudo`, automatically detect `SUDO_USER` / `SUDO_UID` / `SUDO_GID` and ensure any created files or directories in the target user's home directory are chowned back to that user.

3. **Unified Cross-Directory Macro Discovery & Operations**:
   - `LoadMacro(name)`: Sequentially search all candidate directories until the macro definition is found.
   - `ListMacros()`: Scan all candidate directories and return deduplicated macros by name.
   - `MacroExists(name)`: Check for macro presence across all candidate paths.
   - `DeleteMacro(name)`: Remove macro files across all candidate paths where present.
   - `SaveMacro(&m)`: Write atomically to the first verified writable candidate directory.

4. **Directory Writability Verification**:
   - Verify write permissions dynamically by creating and immediately removing an ephemeral probe token (`.perm_test_<timestamp>`) before committing writes.
