---
name: install-tar-linux
description: Autonomously install .tar, .tar.gz, .tgz, .gz, and .zip archives on Linux with intelligent format detection, installation strategy resolution, PATH symlinking, desktop integration, and dual-database uninstallation tracking.
---

# Linux Archive Installer (`gitmap install tar`)

## Overview
This skill guides the installation of archive packages (`.tar`, `.tar.gz`, `.tgz`, `.gz`, `.zip`, `.tar.bz2`, `.tar.xz`) on Linux/Ubuntu via `gitmap install tar <archive-path>`.

## Core Requirements & Execution Steps
1. **Linux-Only Enforcement**:
   - Verify `runtime.GOOS == "linux"`. If run on non-Linux, emit structured `*apperror.AppError` explaining Linux support only.
2. **Intelligent Archive Detection**:
   - Discriminate between `.tar.gz`/`.tgz`, `.tar.xz`/`.txz`, `.tar.bz2`, `.tar`, single compressed `.gz`, and `.zip` files by file signature/magic bytes and extension.
3. **Extraction & Strategy Detection**:
   - Extract cleanly into repository-scoped temporary directory (`tempdir.RepoTempDir("archive-install")`).
   - Analyze extracted content for:
     - Pre-compiled binaries (`bin/<app>`, root binary, ELF format).
     - Standard install scripts (`install.sh`, `setup.sh`).
     - Build configurations (`Makefile`, `./configure`, `CMakeLists.txt`).
     - Single binary unpacked from `.gz`.
     - `.desktop` files and application icons (`.png`, `.svg`).
4. **Target Destination & Symlinking**:
   - Extract/install to `~/.local/share/<appname>/` (or `/opt/<appname>/` when run with root/sudo).
   - Create symlink in `~/.local/bin/<appname>` (or `/usr/local/bin/<appname>`), ensuring `0755` permissions.
   - Install `.desktop` to `~/.local/share/applications/` and run `update-desktop-database`.
5. **Database Tracking & Uninstallation Parity**:
   - Record installation in `gitmap.db` and `installation.db` as a custom archive tool.
   - Support `gitmap uninstall <appname>` cleanly removing the directory, symlink, desktop entry, and DB records.
6. **Coding Guidelines**:
   - Functions <= 15 lines.
   - Single return types or structured Result wrappers.
   - Affirmative booleans (`is*`, `has*`), zero explicit `== true`.
   - Mandatory `*apperror.AppError` with stack traces.
