# Root Cause Analysis: Release Skew Incident

## Incident Overview

During the `v6.87.1` deployment cycle, the CI/CD pipeline correctly triggered but abruptly failed during the `Installer Smoke Test`, leaving a broken release published on GitHub without its requisite assets or patch notes.

## Technical Post-Mortem

1. **Tag Push out of Band**: An operator/agent executed `git push origin v6.87.1` to force a release bump, but the source code (specifically `constants.go` and `changelog.md`) was untouched and still set to `v6.87.0`.
2. **Missing Patch Notes**: The `release.yml` GitHub Action triggered. Its `awk` script parsed `changelog.md` for `[v6.87.1]`. Finding none, it output `"No changelog entry found"` directly into the GitHub Release body.
3. **Smoke Test Failure**: The GitHub Action compiled the Go binary (`go build -ldflags...`). However, because `constants.Version` is a `const` in Go, `-ldflags` was incapable of overriding it at build time. The binary was natively compiled as `v6.87.0`.
4. **Fatal Crash**: The workflow proceeded to the Installer Smoke Test phase. It executed the newly built binary expecting the output `gitmap v6.87.1`. The binary instead output `v6.87.0`. The script forcefully aborted (`Version mismatch. expected: v6.87.1, actual: v6.87.0`).

## Resolution & Prevention

To resolve this, an orchestrating AI agent explicitly executed the **PowerShell Version Sweep** across all canonical files, pushed a `main` commit, and then issued a `v6.88.0` tag in series. 

**Prevention Rule**: All future automated releases must adhere strictly to the procedures documented in `01-ai-release-synchronization.md`. Source control dictates the tag; tags never dictate the source control.
