---
name: fix-antigravity-and-universal-uninstall
description: Autonomously diagnose, repair, and verify Antigravity IDE and CLI installations using official Google Cloud Storage artifacts and universal uninstall engine across Linux and Windows with full stack traces.
---

# Fix Antigravity Installation & Universal Uninstall Engine

## Overview
This skill guides the autonomous remediation and verification of Google Antigravity Desktop IDE, Antigravity CLI (`agy`), and universal uninstallation across all tracked tools and platforms.

## Core Rules & Architecture
1. **Official Google Cloud Storage Sources**:
   - Linux x64: `https://storage.googleapis.com/antigravity-public/antigravity-hub/2.13.0-6362815968182272/linux-x64/Antigravity.tar.gz`
   - Linux ARM: `https://storage.googleapis.com/antigravity-public/antigravity-hub/2.13.0-6362815968182272/linux-arm/Antigravity.tar.gz`
   - Windows x64: `https://storage.googleapis.com/antigravity-public/antigravity-hub/2.13.0-6362815968182272/windows-x64/Antigravity-x64.exe`
   - macOS ARM: `https://storage.googleapis.com/antigravity-public/antigravity-hub/2.13.0-6362815968182272/darwin-arm/Antigravity.dmg`
   - macOS x64: `https://storage.googleapis.com/antigravity-public/antigravity-hub/2.13.0-6362815968182272/darwin-x64/Antigravity.dmg`
2. **True Installation Verification**:
   - Never report "installed" unless the binary exists on disk, has non-zero size, is executable (`chmod +x` on Unix), and produces valid version output.
   - Do not leave orphan `.desktop` files or broken symlinks if download or extraction fails.
3. **Universal Uninstallation**:
   - Every supported installable tool (CLI, Desktop IDE, custom scripts, package manager tools) MUST have a fully functional `uninstall` flow.
   - Support both flags before and after positional tool arguments using `reorderFlagsBeforeArgs`.
   - Dual-database purging (`gitmap.db` and `installation.db`).
   - Clean disk file sweeping for standalone tools.
4. **AppError with Mandatory Stack Trace**:
   - Always wrap execution failures in domain `*apperror.AppError` preserving caller stack traces.
