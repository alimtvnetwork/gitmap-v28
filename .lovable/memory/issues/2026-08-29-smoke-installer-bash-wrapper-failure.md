# RCA — smoke-installer wrapper script failure in CI

**Date:** 2026-08-29
**Workflow:** .github/workflows/release.yml and .github/workflows/ci.yml
**Symptom:**

`	ext
Run bash .github/scripts/smoke-installer.sh release
  bash .github/scripts/smoke-installer.sh release
  shell: /usr/bin/bash --noprofile --norc -e -o pipefail {0}
`

The CI job fails when invoking the bash wrapper script smoke-installer.sh.

## Root cause

The scripts .github/scripts/smoke-installer.sh and .github/scripts/smoke-installer.ps1 were legacy wrappers that simply forwarded arguments to the cross-platform smoke-installer.py. Using a ash wrapper in GitHub Actions (especially on Windows runners using MSYS2 bash) introduces pathing, environment, and execution translation issues that cause the wrapper to fail before it can even invoke the Python script.

## Fix

1. **Removed Wrappers:** Deleted .github/scripts/smoke-installer.sh and .github/scripts/smoke-installer.ps1 entirely.
2. **Direct Python Invocation:** Updated .github/workflows/ci.yml and .github/workflows/release.yml to execute the Python script natively.
   - Old:
un: bash .github/scripts/smoke-installer.sh release
   - New:
un: python .github/scripts/smoke-installer.py release

## Permanent guarantees

- **No Shell Wrappers for Python Scripts:** Cross-platform Python tools should be invoked directly via python script.py in GitHub Actions. Avoid wrapping Python scripts in .sh or .ps1 proxy files, as they introduce unnecessary points of failure across different OS runners.
