---
name: antigravity-crossplatform-installer
description: Autonomously implement and verify cross-platform Google Antigravity IDE and CLI installation across Ubuntu/Linux, Windows, and macOS with categorical OS detection, prerequisite validation, retry downloads, AppError envelopes, and uninstallation tracking.
---

# Antigravity Cross-Platform Installer Engine

## Core Instructions

1. **Categorical OS & Architecture Detection**:
   - Detect host operating system (`linux`, `windows`, `darwin`) and CPU architecture (`amd64`/`x64`, `arm64`).
   - Clearly announce target OS and architecture before beginning installation (e.g., `[INFO] Detected Platform: Linux (Ubuntu) [arch: x64]`).

2. **Canonical Google Cloud Storage Artifact Matrix**:
   - Base URL: `https://storage.googleapis.com/antigravity-public/antigravity-hub/{VERSION}-{BUILD_ID}/{PLATFORM}/{ARTIFACT}`
   - Default `VERSION=2.13.0`, `BUILD_ID=6362815968182272`.
   - Windows x64 / arm64: `windows-x64/Antigravity-x64.exe` or `windows-arm/Antigravity-arm64.exe`.
   - macOS Silicon / Intel: `darwin-arm/Antigravity.dmg` or `darwin-x64/Antigravity.dmg`.
   - Linux x64 / arm64: `linux-x64/Antigravity.tar.gz` or `linux-arm/Antigravity.tar.gz`.

3. **Step-by-Step Progress Display & Error Handling**:
   - Step 1: Detect & announce platform and architecture.
   - Step 2: Validate OS prerequisites (Windows build, macOS version, Linux glibc).
   - Step 3: Check idempotency (skip if already installed unless `--force`).
   - Step 4: Download with 3 retries and progress logging into repo temp dir.
   - Step 5: Install (Windows silent `/S`, macOS DMG mount & ditto, Linux tarball extraction & chmod).
   - Step 6: Link binaries (`antigravity` and `agy` into PATH/symlinks) & Linux `.desktop` launcher.
   - Step 7: Post-install validation and dual-database registration (`installation.db` and `gitmap.db`).

4. **DRY & AppError Architecture**:
   - Use AppError envelopes on all error returns (`*apperror.AppError`).
   - Functions <= 15 lines (target <= 8 lines).
   - Affirmative booleans only (`is*`, `has*`).
