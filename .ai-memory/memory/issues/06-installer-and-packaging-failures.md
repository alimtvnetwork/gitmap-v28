# Issue Domain 06: Installer and Packaging Failures

- **Domain:** NSIS Auto-Detection, Shell Installers, and Binary Extraction
- **Status:** Consolidated Problem & Resolution Matrix

## 1. Missing Binary in Smoke Installer Tests
- **Symptoms:** Installer smoke test failed with "Binary not found at $DEST/gitmap".
- **Root Cause:** Archive extraction path flattened the destination directory unexpectedly.
- **Resolution:** Enforced deterministic destination paths and added post-extraction file existence verification.

## 2. Path Delimiter Mismatches in Cross-Platform Installers
- **Symptoms:** `run.ps1` failed when executed on Windows with Unix-style path separators.
- **Root Cause:** Raw string concatenations mixed `/` and `\` delimiters.
- **Resolution:** Standardized on `filepath.Join` across Go and `Join-Path` across PowerShell scripts.
