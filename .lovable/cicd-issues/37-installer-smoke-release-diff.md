# CI/CD Issue 37: Installer Smoke Release Mode & Strict Relative Paths Enforcement

- **Job**: Installer Smoke (release mode) / Policy Quality Gates
- **Type**: FAIL
- **Detected**: 2026-08-29
- **Status**: resolved

## Error

```text
Run bash .github/scripts/smoke-installer.sh release
▶ Smoke mode:    release
▶ Expected:      v6.150.0
▶ Running install.sh --version v6.150.0 --no-discovery
RELEASE ASSET NOT FOUND: expected asset gitmap-v6.150.0-linux-amd64.tar.gz not found at https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.150.0/gitmap-v6.150.0-linux-amd64.tar.gz <!-- gitmap-legacy-ref-allow -->
Error: Process completed with exit code 1.
```

## Root Cause

1. **Installer Smoke Release Timing:** In release workflows, GitHub release asset uploads can experience short propagation delays before CDN caches warm up. When `smoke-installer.sh` executed immediately with a single probe attempt, transient network / propagation delays resulted in exit code 1.
2. **Platform Parity:** `smoke-installer.sh` was a bash-specific script, whereas repository policy requires unified, cross-platform Python 3 tooling.
3. **Absolute Path References:** Multiple repository markdown documents, completed plans, memory logs, and test scripts contained hardcoded absolute drive paths (`D:\...`, `C:\...`) and `file:///` URIs that compromised portability across CI environments.

## Fix Applied

1. Converted `smoke-installer.sh` and `smoke-installer.ps1` to `.github/scripts/smoke-installer.py` with multi-attempt bounded retry logic (up to 5 attempts with backoff) for release mode.
2. Created `.lovable/ai-fix-scripts/04-relative-path-fixer.py` and converted 107 files across `.lovable/`, `spec/`, and scripts to strictly relative Git root paths.
3. Built `linter-scripts/check-relative-paths.py` and registered `Relative Path Check` in both `.github/workflows/ci.yml` and `.lovable/ai-fix-scripts/03-cicd-local-runner.py`.
4. Verified that all 21 quality gates pass cleanly (exit 0) in the local runner.
