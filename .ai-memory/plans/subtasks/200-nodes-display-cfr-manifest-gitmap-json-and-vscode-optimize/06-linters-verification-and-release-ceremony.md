# Subtask 06: Linters Verification and Release Ceremony

> **Parent Plan:** `200-nodes-display-cfr-manifest-gitmap-json-and-vscode-optimize.md`  
> **Status:** `PENDING`  
> **Target Files:**
> - `version.json`
> - `package.json`
> - `cli/constants/constants.go`
> - `changelog.md`
> - `02-spec/98-changelog.md`
> - `02-spec/99-release-notes/release-v6.425.0.md`

---

## Technical Specification

1. **Local Python Linter Verification**:
   - `python linter-scripts/check-nested-ifs.py` (0 errors)
   - `python linter-scripts/check-enum-and-boolean.py` (0 errors)
   - `python linter-scripts/check-boolean-guidelines.py` (0 errors)
   - `python linter-scripts/check-relative-paths.py` (0 errors)
   - `python linter-scripts/check-error-management.py` (0 errors)

2. **Version Bump**:
   - `v6.425.0` in `version.json`, `package.json`, and `cli/constants/constants.go`.
   - Update `changelog.md` and `02-spec/98-changelog.md`.
   - Author release notes.

3. **Atomic Commit & Push**:
   - Commit all changes in a single grouped commit.
   - Create git tag `v6.425.0`.
   - Push `main`, `release/v6.425.0`, and tag `v6.425.0`.
   - Monitor CI/CD with `gitmap pipeline-ai status --json` until green.
