# Root Cause Analysis: Installer Smoke Release Mode & Strict Relative Paths Enforcement

## 1. Context and Problem Statement
During GitHub Actions release pipeline runs, the `Installer Smoke (release)` job failed when attempting to probe the published release assets immediately after tag publication due to propagation latency on release asset download URLs. Additionally, multiple documentation, spec, and plan files contained absolute drive paths (`D:\...`, `C:\...`) and `file:///` URIs.

## 2. Root Cause
- `smoke-installer.sh` was written in bash with a single probe attempt without backoff, causing transient failure if asset CDN propagation had not finished.
- Absolute machine-specific file URIs (`file:///d:/work/gitmap/...`) were recorded during automated audit planning runs, breaking cross-platform portability.

## 3. Corrective and Preventive Actions
- Converted `smoke-installer.sh` and `smoke-installer.ps1` to `.github/scripts/smoke-installer.py`, adding robust multi-attempt retry logic with backoff for `release` mode.
- Created `.lovable/ai-fix-scripts/04-relative-path-fixer.py` and converted all absolute paths repository-wide to relative Git root paths.
- Built `linter-scripts/check-relative-paths.py` and registered `Relative Path Check` in both `.github/workflows/ci.yml` and `.lovable/ai-fix-scripts/03-cicd-local-runner.py`.

## 4. Verification
- Ran `python .github/scripts/smoke-installer.py source` and `python .github/scripts/smoke-installer.py release` (both passed).
- Ran `python linter-scripts/check-relative-paths.py` (passed with 0 violations across 6,478 tracked files).
- Ran `python .lovable/ai-fix-scripts/03-cicd-local-runner.py` (all 21 quality gates passed with exit 0).
