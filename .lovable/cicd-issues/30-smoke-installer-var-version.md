# 30-smoke-installer-var-version

## Error Summary

The CI job `Installer Smoke (source-build)` failed during step `Run installer smoke (source mode)`:
```text
Run bash .github/scripts/smoke-installer.sh source
Error: Could not determine expected version
Error: Process completed with exit code 2.
```

---

## 4-Part Root Cause Analysis

### 1. Why it happened

In `.github/scripts/smoke-installer.sh`, the version fallback logic extracted the version from `gitmap/constants/constants.go` using:
```bash
EXPECTED="${EXPECTED:-$(awk -F'"' '/^const Version/ {print $2}' "$REPO_ROOT/gitmap/constants/constants.go")}"
```
Because `constants.Version` is declared as a package variable (`var Version = "6.144.0"`) so it can be overridden at link-time via `-ldflags`, the regex `/^const Version/` found no matching lines and returned an empty string. This caused the subsequent guard `if [ -z "$EXPECTED" ]; then echo "::error::Could not determine expected version" >&2; exit 2; fi` to fail with exit code 2.

### 2. How it happened

- In `ci.yml`, job `Installer Smoke (source-build)` ran `bash .github/scripts/smoke-installer.sh source`.
- `smoke-installer.sh` evaluated line 20 without `$EXPECTED` pre-set in the environment.
- `awk` scanned `gitmap/constants/constants.go` for `^const Version`, yielding empty output.
- Line 83 detected `$EXPECTED` was empty and exited 2.

### 3. Root Cause

Rigid regex pattern `/^const Version/` in `.github/scripts/smoke-installer.sh` assuming `Version` was declared with `const` rather than `var`.

### 4. Code Fix

- In `.github/scripts/smoke-installer.sh`, updated the extraction pattern to match both `var` and `const`:
  ```bash
  EXPECTED="${EXPECTED:-$(grep -E '^(var|const) Version\b' "$REPO_ROOT/gitmap/constants/constants.go" | head -1 | sed -E 's/.*"([^"]+)".*/\1/')}"
  ```
- Verified source-mode smoke test passes: `✅ Installer smoke test passed: gitmap v6.144.0`.
- Added `Installer Smoke` to `.lovable/ai-fix-scripts/03-cicd-local-runner.py`.
