# CI/CD Issue 43: Installer Seed Data 404s & Release Dry-Run Script Path Drift

- **Job**: Installer Dry-Run Windows & Installer Smoke Windows (release)
- **Type**: FAIL
- **Detected**: 2026-09-13
- **Status**: resolved
- **Pipeline Run**: #34744532717

## Error

```text
Section [1/2]: Release #34744532717 ➔ Job: Installer Dry-Run Windows | Step: Run install.ps1 -DryRun and capture report
pwsh -File gitmap/scripts/install.ps1 `
The argument 'gitmap/scripts/install.ps1' to the -File parameter does not exist. Provide the path to an existing '.ps1' file as an argument to the -File parameter.
Write-Error "::error::install.ps1 -DryRun exited 64"

Section [2/2]: Release #34744532717 ➔ Job: Installer Smoke Windows (release) | Step: Run installer smoke (release mode)
Install self-check FAILED: required seed file is missing.
  Expected: C:\Users\RUNNER~1\AppData\Local\Temp\...\gitmap-cli\data\downloader-config.json
  Why this matters:
    gitmap reads downloader-config.json on every run to configure parallel downloads...
Write-Error "::error::Install-SeedData failed for downloader-config.json"
```

## Root Cause

1. **Repository Layout Migration Drift in `release.yml`**:
   The GitMap codebase migrated from `gitmap/` to `cli/`. In `.github/workflows/release.yml:887`, the dryrun step still invoked `pwsh -File gitmap/scripts/install.ps1`, which no longer existed in the repo checkout.
2. **Seed Data URLs Hardcoded to `gitmap/data/`**:
   In `cli/scripts/install.ps1` (and `cli/scripts/install.sh`), `Install-SeedData` attempted to download seed files (`downloader-config.json`, `config.json`, `git-setup.json`, `seo-templates.json`) from `https://raw.githubusercontent.com/$Repo/$version/gitmap/data/$name` and fallback `https://raw.githubusercontent.com/$Repo/main/gitmap/data/$name`. Both returned HTTP 404 because seed files were moved to `cli/data/`.
3. **No Local Checkout Fallback**:
   When tests or offline environments executed the installer, neither script checked for the presence of local data files next to the installer script, leading to complete failure if remote endpoints returned 404 or experienced latency.
4. **Stale Go Constants and Helptext**:
   Embedded strings in `cli/constants/constants_selfinstall.go`, `constants_release.go`, `constants_messages.go`, and documentation helptext in `cli/helptext/` retained dead `gitmap/scripts/` URLs.

## Fix Applied

1. **Updated `.github/workflows/release.yml`**:
   - Changed dryrun invocation from `pwsh -File gitmap/scripts/install.ps1` to `pwsh -File cli/scripts/install.ps1`.
   - Updated release notes template installer URLs to `cli/scripts/install.ps1` and `cli/scripts/install.sh`.
2. **Hardened PowerShell & Bash Installers**:
   - Updated primary and fallback URLs to `cli/data/$name`.
   - Added candidate local file fallback paths (`..\data\$name`, `data\$name`, `..\..\cli\data\$name`, `cli\data\$name`) in `Install-SeedData` and `install_seed_data()`.
   - Added local checkout candidate fallback to `Assert-InstallSelfCheck` prior to throwing failure.
   - Synchronized root `install.ps1` and `install.sh` byte-for-byte with `cli/scripts/`.
3. **Synchronized Go Constants & Helptext**:
   - Updated `SelfInstallRemotePwsh`, `SelfInstallRemoteBash`, `MsgInstallHintWindows`, `MsgInstallHintUnix`, `ReleaseSnippetTemplate`, and `MsgUpdateHelpNoRepo` to point to `cli/scripts/`.
   - Updated `cli/cmdinstall/installscripts.go` default sources to `cli/scripts/`.
   - Updated command help markdown docs in `cli/helptext/`.
4. **Version Bump & Release**:
   - Bushed minor version from `v6.226.0` to `v6.227.0` via `python 03-ai-scripts/29-release-orchestrator.py --tier minor --skip-tests`.
   - Validated dry-run command locally with exit code 0.