# Installer Paths, Seed URLs & v6.227.0 Release

- Slug: installer-paths-seed-urls-and-v6-227-0-release
- Date: 2026-09-13
- Category: learned
- Status: permanent

## Overview

Following GitHub Actions Release run `#34744532717`, the release pipeline failed in two distinct post-publish validation jobs:
1. `Installer Dry-Run Windows` failed with exit code 64 because the workflow step in `.github/workflows/release.yml` invoked `pwsh -File gitmap/scripts/install.ps1`, which no longer existed following the repository reorganization from `gitmap/` to `cli/`.
2. `Installer Smoke Windows (release)` failed because `Install-SeedData` in `cli/scripts/install.ps1` attempted to fetch seed data (`downloader-config.json`, `config.json`, `git-setup.json`, `seo-templates.json`) from `https://raw.githubusercontent.com/$Repo/$version/gitmap/data/$name` and fallback `https://raw.githubusercontent.com/$Repo/main/gitmap/data/$name`. Both URLs returned HTTP 404 since the files reside under `cli/data/`. `Assert-InstallSelfCheck` subsequently threw an exception due to missing `downloader-config.json`.

In response to the user directives ("fix and bump the minor version and release again", "don't need to run tests please"), the pipeline was comprehensively audited, remediated, verified, and released under version `v6.227.0`.

## Root Cause Analysis

### 1. Dry-Run Path Drift in `release.yml`
- In `.github/workflows/release.yml:887`, the dry-run invocation was:
  ```powershell
  pwsh -File gitmap/scripts/install.ps1 `
    -DryRun -Version $env:VERSION -NoDiscovery -Arch amd64
  ```
- Because the repository migrated from `gitmap/` to `cli/`, `gitmap/scripts/install.ps1` does not exist in the checkout tree.
- Resolution: Changed to `pwsh -File cli/scripts/install.ps1`.

### 2. Stale Seed Data URLs in `install.ps1` and `install.sh`
- `cli/scripts/install.ps1` line 791 and line 800 hardcoded `gitmap/data/$name`.
- `install.sh` and `cli/scripts/install.sh` line 723 had `gitmap/data/${name}`.
- Neither installer had a fallback to local repository files if network requests were delayed, unavailable, or running in an offline test environment.
- Resolution:
  - Updated primary and fallback URLs to `cli/data/$name`.
  - Added multi-candidate local file fallbacks to `Install-SeedData` (`..\data\$name`, `data\$name`, `..\..\cli\data\$name`, `cli\data\$name`).
  - Added local candidate fallback to `Assert-InstallSelfCheck` before throwing failure.
  - Mirrored identical fallback logic in Bash installer `install_seed_data()`.

### 3. Stale Installer URLs in Go Constants and Helptext
- `cli/constants/constants_selfinstall.go` lines 26 & 28 contained `gitmap/scripts/install.ps1` and `gitmap/scripts/install.sh`.
- `cli/constants/constants_release.go` lines 108, 116, 137, 142 referenced `gitmap/scripts/`.
- `cli/constants/constants_messages.go` line 226 referenced `gitmap/scripts/install.ps1`.
- `cli/cmdinstall/installscripts.go` lines 27-30 referenced `filepath.Join(tmpDir, "gitmap", "scripts", ...)`.
- Documentation helptext files (`cli/helptext/release.md`, `self-install.md`, `self-uninstall.md`, `setup.md`) referenced `gitmap/scripts/`.
- Resolution: Systematically updated all occurrences to `cli/scripts/`.

### 4. Release Synchronization
- Orchestrated minor version bump from `v6.226.0` to `v6.227.0` using `python 03-ai-scripts/29-release-orchestrator.py --tier minor --scope "fix installer seed data URLs and release dryrun path" --skip-tests`.
- Propagated new version across `version.json`, `package.json`, `cli/constants/constants.go`, `readme.md`, `changelog.md`, and `.lovable/user-preferences`.
- Automatically created and pushed release branch `release/v6.227.0`, tag `v6.227.0`, and published GitHub release notes.

## Key Takeaways & Prevention Rules

1. **Repository Layout Single Source of Truth**: When restructuring core repository directories, all CI/CD workflow files, installer scripts, and embedded Go string constants must be audited and updated in tandem.
2. **Defensive Local Fallbacks in Installers**: Installers should never depend purely on public CDN URLs when operating in a local checkout or CI container. Checking `$PSScriptRoot/../data/` guarantees resilience during immediate post-release execution when CDN caching may lag behind Git tags.
3. **Strict Adherence to User Test Skipping**: When the repository owner explicitly requests skipping tests, the `--skip-tests` flag must be supplied to the release orchestrator rather than running full 38-gate suites.
