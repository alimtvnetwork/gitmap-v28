# Plan 149: Installer URLs, Release Dry-Run Path & v6.227.0 Release

- Slug: 149-installer-urls-dryrun-path-and-v6-227-0-release
- Date: 2026-09-13
- Status: completed

## 1. Executive Summary & Objective

Autonomously audit, remediate, verify, and release GitMap version `v6.227.0` to resolve two post-publish pipeline failures in GitHub Actions Release run `#34744532717`:
1. `Installer Dry-Run Windows`: Failed with exit code 64 because `.github/workflows/release.yml` invoked `pwsh -File gitmap/scripts/install.ps1`, which was missing due to repository migration to `cli/scripts/install.ps1`.
2. `Installer Smoke Windows (release)`: Failed because `Install-SeedData` in `cli/scripts/install.ps1` downloaded from dead `https://raw.githubusercontent.com/$Repo/$version/gitmap/data/$name` URLs instead of `cli/data/$name`, causing `Assert-InstallSelfCheck` to fail on missing `downloader-config.json`.

Per user directives ("fix and bump the minor version and release again", "don't need to run tests please"), the pipeline paths, seed data URLs, local fallbacks, Go constants, and release body templates were completely synchronized, verified, and released under `v6.227.0`.

## 2. Changes Summary

1. **GitHub Actions Release Workflow (`.github/workflows/release.yml`)**:
   - Updated dryrun command at line 887 to `pwsh -File cli/scripts/install.ps1`.
   - Updated release notes body installer URLs to `cli/scripts/install.ps1` and `cli/scripts/install.sh`.
   - Updated job comments to reference `cli/scripts/`.

2. **PowerShell & Bash Installers (`install.ps1`, `cli/scripts/install.ps1`, `install.sh`, `cli/scripts/install.sh`)**:
   - Updated seed data download URLs to `cli/data/` across PowerShell and Bash scripts.
   - Added robust local repository candidate fallback checks for seed files (`downloader-config.json`, `config.json`, `git-setup.json`, `seo-templates.json`).
   - Added local checkout candidate fallback in `Assert-InstallSelfCheck` before throwing failure.
   - Synchronized root `install.ps1` and `install.sh` to match `cli/scripts/` versions byte-for-byte.

3. **Go Constants & Command Helptext**:
   - `cli/constants/constants_selfinstall.go`: Updated remote fallback URLs to `cli/scripts/`.
   - `cli/constants/constants_release.go`: Updated install hint snippets to `cli/scripts/`.
   - `cli/constants/constants_messages.go`: Updated re-install one-liner to `cli/scripts/install.ps1`.
   - `cli/cmdinstall/installscripts.go`: Updated `defaultScriptSources` to use `cli/scripts/`.
   - `cli/downloaderconfig/types.go` & `cli/constants/deploy_manifest.go`: Updated comments to `cli/`.
   - Documentation helptext (`cli/helptext/release.md`, `self-install.md`, `self-uninstall.md`, `setup.md`): Updated commands to `cli/scripts/`.

4. **Release Orchestration (`03-ai-scripts/29-release-orchestrator.py`)**:
   - Executed `--tier minor --scope "fix installer seed data URLs and release dryrun path" --skip-tests`.
   - Bumped canonical version to `6.227.0` across `version.json`, `package.json`, `cli/constants/constants.go`, `readme.md`, `changelog.md`.
   - Pushed `release/v6.227.0` branch and `v6.227.0` tag, published GitHub release notes.

## 3. Verification

- Ran `pwsh -File cli/scripts/install.ps1 -DryRun -Version v6.226.0 -NoDiscovery -Arch amd64` locally: exit code 0, verified report matches contract.
- Verified Go package compilation: clean build across all CLI packages.
- GitHub Actions Release pipeline `#34746417661` triggered by `v6.227.0` tag.
