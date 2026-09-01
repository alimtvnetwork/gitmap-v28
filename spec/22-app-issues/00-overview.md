# GitMap Application Issues & Root Cause Analysis (RCA)

**Version:** 3.2.0  
**Updated:** 2026-09-01  
**AI Confidence:** Production-Ready  
**Ambiguity:** None

---

## Overview

The `spec/22-app-issues/` directory is the canonical repository home for all **GitMap Issue Reports, Root Cause Analysis (RCA), and Fix Specifications**. 

This folder tracks production-grade post-mortems, reproduction steps, underlying causes, architectural fixes, and regression prevention rules for the GitMap CLI ecosystem.

---

## Issue Catalog

| ID | Specification | Topic |
|---|---|---|
| 01 | [`01-update-file-lock.md`](./01-update-file-lock.md) | Windows binary file lock during self-update |
| 02 | [`02-update-flow-spec-alignment.md`](./02-update-flow-spec-alignment.md) | Self-update flow and CLI flag alignment |
| 03 | [`03-update-sync-lock-loop.md`](./03-update-sync-lock-loop.md) | Update sync lock retry loop |
| 04 | [`04-database-path-resolution.md`](./04-database-path-resolution.md) | SQLite database path resolution across OS environments |
| 05 | [`05-list-empty-db-path.md`](./05-list-empty-db-path.md) | Graceful handling when listing uninitialized databases |
| 06 | [`06-release-orphaned-meta.md`](./06-release-orphaned-meta.md) | Orphaned release metadata cleanup |
| 07 | [`07-zip-group-release-silent-failure.md`](./07-zip-group-release-silent-failure.md) | Zip group asset bundling resilience |
| 08 | [`08-autocommit-push-rejection.md`](./08-autocommit-push-rejection.md) | Auto-commit push rejections & branch synchronization |
| 09 | [`09-list-releases-repo-source.md`](./09-list-releases-repo-source.md) | Repository source attribution for releases |
| 10 | [`10-legacy-uuid-detection.md`](./10-legacy-uuid-detection.md) | Legacy UUID migration |
| 11 | [`11-auto-legacy-dir-migration.md`](./11-auto-legacy-dir-migration.md) | Automatic directory migrations |
| 12 | [`12-legacy-id-migration.md`](./12-legacy-id-migration.md) | Legacy repository ID normalization |
| 13 | [`13-release-pipeline-dist-directory.md`](./13-release-pipeline-dist-directory.md) | Distribution directory structure in release pipeline |
| 14 | [`14-security-hardening-gosec-fixes.md`](./14-security-hardening-gosec-fixes.md) | Gosec security audit hardening |
| 15 | [`15-installer-progress-bar-and-binary-detection.md`](./15-installer-progress-bar-and-binary-detection.md) | Installer progress and binary detection |
| 16 | [`16-ci-passthrough-gate-pattern.md`](./16-ci-passthrough-gate-pattern.md) | CI passthrough gate patterns |
| 17 | [`17-go-flag-ordering-silent-drop.md`](./17-go-flag-ordering-silent-drop.md) | Flag order normalization |
| 18 | [`18-ci-release-branch-cancellation-protection.md`](./18-ci-release-branch-cancellation-protection.md) | CI branch protection |
| 19 | [`19-missing-macos-binaries-and-lint-regression.md`](./19-missing-macos-binaries-and-lint-regression.md) | Cross-platform macOS binary packaging |
| 20 | [`20-path-not-available-in-other-shells.md`](./20-path-not-available-in-other-shells.md) | Cross-shell PATH environment activation |
| 21 | [`21-pending-task-durability.md`](./21-pending-task-durability.md) | Pending task serialization durability |
| 22 | [`22-installer-path-not-active-after-install.md`](./22-installer-path-not-active-after-install.md) | Post-install shell environment refresh |
| 23 | [`23-go-build-copyfile-redeclared.md`](./23-go-build-copyfile-redeclared.md) | Build script namespace collisions |
| 24 | [`24-cd-command-does-not-change-shell-directory.md`](./24-cd-command-does-not-change-shell-directory.md) | Shell wrapper directory navigation |
| 25 | [`25-powershell-cd-wrapper-not-loaded.md`](./25-powershell-cd-wrapper-not-loaded.md) | PowerShell profile module loader |
| 26 | [`26-docs-site-not-bundled-and-swallowed-errors.md`](./26-docs-site-not-bundled-and-swallowed-errors.md) | Docs site bundling error handling |
| 27 | [`27-error-management-file-path-and-missing-file-code-red-rule.md`](./27-error-management-file-path-and-missing-file-code-red-rule.md) | Error code allocation & CODE RED enforcement |
| 28 | [`28-unused-cd-profile-path-lint-failure.md`](./28-unused-cd-profile-path-lint-failure.md) | Dead code removal & lint hygiene |
| 29 | [`29-macos-pwsh-shell-not-activated-after-install.md`](./29-macos-pwsh-shell-not-activated-after-install.md) | macOS PowerShell profile activation |
| 31 | [`31-update-cleanup-phase3-observability-gap.md`](./31-update-cleanup-phase3-observability-gap.md) | Update cleanup observability |
| 32 | [`32-docs-ui-vscode-grading-missed-request.md`](./32-docs-ui-vscode-grading-missed-request.md) | Docs UI grading |
| 34 | [`34-hd-hosted-docs-fallback.md`](./34-hd-hosted-docs-fallback.md) | Hosted docs fallback |

---

## Placement Rule

All bug analysis reports, root cause investigations, and post-mortems belong in `22-app-issues/`.

---

## Cross-References

| Reference | Location |
|-----------|----------|
| App Specifications | [../21-app/00-overview.md](../21-app/00-overview.md) |
| Coding Guidelines | [../02-coding-guidelines/00-overview.md](../02-coding-guidelines/00-overview.md) |
| Error Management | [../03-error-manage/00-overview.md](../03-error-manage/00-overview.md) |
