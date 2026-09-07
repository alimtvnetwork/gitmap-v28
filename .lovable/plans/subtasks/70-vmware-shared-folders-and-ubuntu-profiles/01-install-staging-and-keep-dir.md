# Subtask 01: Two-Tier Download Staging & Persistent Keep Architecture

## Status
Completed

## Context & Objectives
Ensure that on Linux/Ubuntu, gitmap never pollutes user directories or leaves orphaned download files.
1. **Tier 1 (In-Flight Staging)**: All active downloads and temporary extractions must stage strictly in `/tmp` (or `mktemp -d` / `os.TempDir()`). On cancellation or failure, transient files are automatically cleaned up.
2. **Tier 2 (Persistent Keep)**: If any downloaded packages, tools, or archives must be retained, store them inside `~/.gitmap-installation` (or `/root/.gitmap-installation` if operating as root).
   - `~/.gitmap-installation/downloads/`: Verified installers, tarballs, .deb packages.
   - `~/.gitmap-installation/scripts/`: Generated helper scripts.
   - `~/.gitmap-installation/logs/`: Installation logs and receipts.
3. **Cross-Device EXDEV Link Fallback**: Moving files from `/tmp` (often a `tmpfs` RAM disk) to `${HOME}` will fail with `EXDEV` on standard `os.Rename`. Implement `PromoteStagedFile(src, dst)` that attempts rename, and if `EXDEV`, falls back to buffered `io.Copy`, `Sync()`, and `os.Remove(src)`.

## Files to Create / Modify
- [NEW] `gitmap/downloaderconfig/staging.go` (<= 200 lines): Staging directory resolver, persistent keep directory resolver, and `PromoteStagedFile` with `EXDEV` fallback.
- [NEW] `gitmap/downloaderconfig/staging_test.go` (<= 200 lines): Unit tests validating path resolution, directory creation, and promotion across simulated partitions.
- [MODIFY] `gitmap/constants/constants_downloader.go`: Declare `DirInFlightPrefix = "gitmap-stage-"`, `DirPersistentKeep = ".gitmap-installation"`, and subfolders `downloads`, `scripts`, `logs`.

## Verification Steps
- `go test -v ./gitmap/downloaderconfig/... -run "TestStaging"` passes cleanly.
- Verify functions adhere to $\le 15$ lines, blank line before returns, affirmative booleans.
