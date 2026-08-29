# Memory Issue: 2026-08-29 LFS Zip Drift, Changelog Sync, and JQ Diff ArgJSON

## 1. Why it happened
CI failed simultaneously on 3 distinct validation jobs:
1. `Generate Drift Check` failed because Git LFS tracked `*.zip`, but three historical settings zip files in `data/` and `settings/` were committed directly as binary blobs. When `git diff` ran in the job without directory scoping, Git reported these LFS-filter differences as code generation drift.
2. `Installer Smoke` failed because `constants.Version` was pinned to `6.143.0`, whereas `version.json`, `readme.md`, and `changelog.md` had advanced to `6.144.0`.
3. `Lint Baseline Guard` crashed with `jq: invalid JSON text passed to --argjson` when parsing linter messages containing multi-line content or non-numeric line markers.

---

## 2. How it happened
- On GitHub Actions (`ubuntu-latest`), `actions/checkout@v6` detected 3 files that should have been LFS pointers but were raw blobs.
- `Detect drift` ran `git diff --exit-code`, surfacing binary differences in `data/npp-settings/npp-settings.zip` and `settings/` as drift.
- `check-changelog-version-sync.sh` checked `constants.go` for `var Version = "6.143.0"` against `changelog.md`, failing because `changelog.md` documented `6.144.0`.
- In `check-single-linter-diff.sh`, `gocritic` ran against a stale Go 1.24 baseline cache, found diffs, and passed non-numeric tokens into `jq --argjson l "$LINE"`, causing a fatal script crash.

---

## 3. Root Cause
- Unscoped `git diff` in `ci.yml` evaluating non-module root zip files affected by Git LFS.
- Incomplete version propagation across `constants.go` and `changelog.md`.
- Unchecked variable interpolation in bash script calling `jq --argjson`.

---

## 4. Code Fix
- **`.gitattributes`**: Explicitly exempt committed archive files (`data/**/*.zip -filter -diff -merge`, `settings/**/*.zip -filter -diff -merge`).
- **`.github/workflows/ci.yml`**: Scope `git diff --exit-code .` and `git diff --name-only .` to `gitmap/`. Bump baseline cache version to `v2`.
- **`gitmap/constants/constants.go`**: Update `var Version = "6.144.0"`.
- **`changelog.md`**: Update release heading to `## [v6.144.0] - 2026-08-29`.
- **`.github/scripts/check-changelog-version-sync.sh`**: Accept optional `v` in heading regex.
- **`.github/scripts/check-single-linter-diff.sh` & `check-misspell-diff.sh`**: Sanitize `LINE` to ensure integer values before `jq --argjson` and strip newlines from issue text.
- **`.lovable/ai-fix-scripts/03-cicd-local-runner.py`**: Added `run_changelog_sync_check` and `run_generate_drift_check`.
