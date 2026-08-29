# 29-lfs-zip-drift-changelog-sync-and-jq-diff-argjson

## Error Summary

The CI/CD pipeline failed across 3 jobs:
1. `Generate Drift Check`:
   ```text
   Encountered 3 files that should have been pointers, but weren't:
   data/npp-settings/npp-settings.zip
   settings/01 - notepad++/02. Notepad++ settings.zip
   settings/03 - obs/25 - Personal 15 Jul 2025 Malaysia.zip
   Drifted files:
   data/npp-settings/npp-settings.zip
   settings/01 - notepad++/02. Notepad++ settings.zip
   settings/03 - obs/25 - Personal 15 Jul 2025 Malaysia.zip
   Error: Generated files are out of sync with constants.
   Error: Run 'cd gitmap && go generate ./...' locally and commit the result.
   ```
2. `Installer Smoke (source-build)`:
   ```text
   → constants.Version = 6.143.0
   ✗ CHANGELOG drift: constants.Version is 6.143.0 but no matching '## v6.143.0' heading exists in changelog.md.
   ```
3. `Lint Baseline Guard (unused, gosec G115, misspell, gocritic, exhaustive)`:
   ```text
   GOCRITIC DIFF (baseline-diff, full-path only)
   jq: invalid JSON text passed to --argjson
   Process completed with exit code 2.
   ```

---

## 4-Part Root Cause Analysis

### 1. Why it happened

- **Generate Drift Check**: `.gitattributes` defined `*.zip filter=lfs diff=lfs merge=lfs -text`. However, 3 settings zip archives in `data/` and `settings/` had been committed into Git history as raw binary files instead of LFS pointers. In `actions/checkout@v6`, Git flagged them as non-pointer files. The `Generate Drift Check` job ran `git diff --exit-code` without directory scoping (`.`), which caused Git to evaluate the repository root and mistake these LFS filter transformations for code generation drift.
- **Changelog Version Sync**: `version.json` and `readme.md` had already been bumped to `6.144.0`, and `changelog.md` had a heading `## [6.144.0] - 2026-08-29`, but `gitmap/constants/constants.go` remained at `var Version = "6.143.0"`. In addition, `check-changelog-version-sync.sh` required a literal `v` prefix in the markdown heading regex, failing to recognize `## [6.144.0]`.
- **Lint Diff jq Crash**: In `check-single-linter-diff.sh` and `check-misspell-diff.sh`, `IFS='|' read -r FILE LINE TEXT` processed linter outputs. When `LINE` contained non-numeric or multiline text, `jq --argjson l "$LINE"` failed with exit code 2 because `jq` expects valid JSON tokens for `--argjson`. Furthermore, upgrading to Go 1.25 invalidated the older cached Go 1.24 lint baseline.

### 2. How it happened

- `actions/checkout@v6` on Ubuntu checked out the repository without downloading LFS objects. Git reported:
  `Encountered 3 files that should have been pointers, but weren't`.
- `Run go generate ./...` finished cleanly in `gitmap/`.
- `Detect drift` executed `git diff --exit-code` from `gitmap/`. Because no path was provided, Git checked the entire repository, surfacing diffs on the 3 raw zip files.
- `Installer Smoke` ran `check-changelog-version-sync.sh`, extracted `6.143.0` from `constants.go`, and found no matching `## v6.143.0` in `changelog.md`.
- `Lint Baseline Guard` ran `gocritic` on Go 1.25 against the Go 1.24 cached baseline. Upon discovering differences, it piped keys to `jq --argjson l "$LINE"`, crashing on unvalidated input.

### 3. Root Cause

1. Unscoped `git diff` in `generate-check` combined with Git LFS treating committed raw zips as un-smudged files.
2. Version drift across `constants.go` (`6.143.0`) vs `version.json` / `readme.md` / `changelog.md` (`6.144.0`).
3. Unsanitized `LINE` variable passed to `jq --argjson` in diff scripts, plus stale Go 1.24 cache key in `ci.yml`.

### 4. Code Fix

1. **LFS & Drift Scope Fix**:
   - Added `data/**/*.zip -filter -diff -merge` and `settings/**/*.zip -filter -diff -merge` to `.gitattributes`.
   - Updated `ci.yml` `generate-check` to scope `git diff --exit-code .` and `git diff --name-only .` to `gitmap/`.
2. **Version Synchronization**:
   - Updated `var Version = "6.144.0"` in `gitmap/constants/constants.go`.
   - Updated `changelog.md` heading to `## [v6.144.0] - 2026-08-29`.
   - Updated `check-changelog-version-sync.sh` to accept optional `v` in headings: `^##[[:space:]]+\[?v?${VERSION//./\\.}\]?`.
3. **Linter Diff Script Hardening & Cache Bump**:
   - Updated `check-single-linter-diff.sh` and `check-misspell-diff.sh` to sanitize `LINE` (`if ! [[ "$LINE" =~ ^[0-9]+$ ]]; then LINE=1; fi`) and normalize newlines in linter messages (`gsub("[\r\n]+"; " ")`).
   - Bumped `lint_baseline_cache_version` default and fallback to `"v2"` in `ci.yml`.
4. **Local Runner Upgrade**:
   - Added `run_changelog_sync_check` and `run_generate_drift_check` to `.lovable/ai-fix-scripts/03-cicd-local-runner.py`.
