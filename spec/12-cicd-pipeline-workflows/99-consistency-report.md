# Consistency Report — CI/CD Pipeline Workflows

**Version:** 3.2.0  
**Updated:** 2026-04-24

---

## File Inventory

### Root Files

| # | File | Version | Status |
|---|------|---------|--------|
| 1 | `00-overview.md` | 5.0.0 | ✅ Present |
| 2 | `01-shared-conventions.md` | — | ✅ Present |
| 3 | `02-github-release-standard.md` | — | ✅ Present |
| 4 | `03-vulnerability-scanning.md` | — | ✅ Present |
| 5 | `04-install-script-generation.md` | — | ✅ Present |
| 6 | `05-code-signing.md` | — | ✅ Present |
| 7 | `06-self-update-mechanism.md` | 2.0.0 | ✅ Present |
| 8 | `07-release-body-and-changelog.md` | 2.0.0 | ✅ Present |
| 9 | `08-installation-flow.md` | 1.0.0 | ✅ Present |
| 10 | `09-changelog-integration.md` | 1.0.0 | ✅ Present |
| 11 | `10-version-and-help.md` | 1.0.0 | ✅ Present |
| 12 | `11-environment-variable-setup.md` | 1.0.0 | ✅ Present |
| 13 | `12-terminal-output-standards.md` | 1.0.0 | ✅ Present |
| 14 | `13-binary-icon-branding.md` | 1.0.0 | ✅ Present |

### Subfolder: 01-browser-extension-deploy

| # | File | Status |
|---|------|--------|
| 15 | `01-browser-extension-deploy/00-overview.md` | ✅ Present |
| 16 | `01-browser-extension-deploy/01-ci-pipeline.md` | ✅ Present |
| 17 | `01-browser-extension-deploy/02-release-pipeline.md` | ✅ Present |
| 18 | `01-browser-extension-deploy/99-consistency-report.md` | ✅ Present |

### Subfolder: 02-go-binary-deploy

| # | File | Status |
|---|------|--------|
| 19 | `02-go-binary-deploy/00-overview.md` | ✅ Present |
| 20 | `02-go-binary-deploy/01-ci-pipeline.md` | ✅ Present |
| 21 | `02-go-binary-deploy/02-release-pipeline.md` | ✅ Present |
| 22 | `02-go-binary-deploy/03-complete-workflow-reference.md` | ✅ Present |
| 23 | `02-go-binary-deploy/99-consistency-report.md` | ✅ Present |

**Total:** 24 files (excluding this report)

---

## Cross-Reference Integrity

- [x] All overview files list their child documents
- [x] All subfolder overviews link to shared convention files
- [x] Consolidated summary exists at `../17-consolidated-guidelines/15-cicd-pipeline-workflows.md`
- [x] Bidirectional cross-refs with `../14-update/`
- [x] Root files (01–13) all have cross-reference sections
- [x] Subfolder consistency reports present in both archetypes
- [x] All cross-references verified via automated scan (2026-04-13)

---

## v8.0.0 Changes

- Added `02-go-binary-deploy/03-complete-workflow-reference.md` to inventory (was missing from v7.0.0 report)
- Fixed broken anchor in `03-complete-workflow-reference.md` (`#multiple-binaries` → `#multiple-binaries-multi-module-build`)
- Updated total file count from 23 to 24
- Verified all cross-references with `14-update/` including new files 13–15

---

## Summary

- **Errors:** 0
- **Health Score:** 100/100 (A+)

---

## Validation History

| Date | Version | Action |
|------|---------|--------|
| 2026-04-09 | 1.0.0 | Initial folder creation with overview placeholder |
| 2026-04-09 | 2.0.0 | Added shared conventions, two subfolders with CI and release specs |
| 2026-04-09 | 3.0.0 | Added root files 04–07, expanded cross-references |
| 2026-04-10 | 5.0.0 | Updated for 06-self-update-mechanism v2.0.0, added subfolder consistency reports to inventory |
| 2026-04-10 | 6.0.0 | Removed 08-ci-failure-logs (commit-back approach rejected) |
| 2026-04-11 | 7.0.0 | Added 08–13 from gitmap-v2 pipeline specs (installation, changelog, version, env, terminal, icon) |
| 2026-04-13 | 8.0.0 | Added 03-complete-workflow-reference to go-binary-deploy inventory; fixed broken anchor; verified all cross-refs |

---

*Consistency Report — updated: 2026-04-13*
