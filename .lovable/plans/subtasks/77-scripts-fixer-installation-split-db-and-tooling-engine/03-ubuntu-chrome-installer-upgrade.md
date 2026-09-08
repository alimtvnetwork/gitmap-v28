# Subtask 03: Ubuntu Chrome Installer Upgrade

**Parent Plan:** [77-scripts-fixer-installation-split-db-and-tooling-engine.md](../../pending/77-scripts-fixer-installation-split-db-and-tooling-engine.md)
**Status:** Completed
**Target Files:**
- `gitmap/cmd/install_chrome_deb.go`
- `gitmap/cmd/install_chrome_deb_test.go`

## Objectives
1. Refactor `runInstallChromeLinux` in `gitmap/cmd/install_chrome_deb.go` to adopt `scripts-fixer` 6-step pipeline:
   - Step 1: Run `sudo apt-get update`
   - Step 2: Ensure fetch utilities: `sudo apt-get install -y wget curl`
   - Step 3: Direct download of deb package to `/tmp/google-chrome-stable_current_amd64.deb`
   - Step 4: Install via `sudo apt-get install -y /tmp/google-chrome-stable_current_amd64.deb`
   - Step 5: Clean up scratch deb file from `/tmp`
   - Step 6: Verify binary via `google-chrome --version`
2. Record granular execution metrics for each phase into `InstallationLog` in `installation.db`.
3. Add unit test `gitmap/cmd/install_chrome_deb_test.go`.
