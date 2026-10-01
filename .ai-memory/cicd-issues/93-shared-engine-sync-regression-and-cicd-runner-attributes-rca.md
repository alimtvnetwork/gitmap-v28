# RCA-093: Shared Engine Sync Regression and CI/CD Runner Attribute Deletion

**Date:** 2026-09-30
**Status:** Resolved
**Severity:** Critical (CI/CD Pipeline Failure across Boolean Linter, Relative Path Check, Nested If Linter, and Script Unit Tests)
**Affected Workflows:** `CI` (`Boolean & Enum Linter`, `Relative Path Check`, `Nested If Linter`, `Lint Script Unit Tests`)
**Run ID:** `36730811906`

---

## 1. Root Cause Analysis
During an external sync commit (`e92ebb5f`), an older copy of `03-ai-scripts/02-shared-engine.py` and `03-ai-scripts/06-cicd-local-runner.py` from `coding-guidelines` was synced into the repository, overwriting GitMap-specific functions and classes:
1. **`AttributeError: module '02-shared-engine' has no attribute 'chunk_items'`:**
   `chunk_items()` and `WorkerHeartbeatMonitor` were removed from `02-shared-engine.py`. This caused immediate crashes across three policy check linters (`check-enum-and-boolean.py`, `check-relative-paths.py`, and `check-nested-ifs.py`) which import `chunk_items` from `02-shared-engine`.
2. **`AttributeError: module '06-cicd-local-runner' has no attribute 'CICD_DIR'` and `normalize_repo_rel`:**
   `REPO_ROOT`, `CICD_DIR`, `normalize_repo_rel`, `format_banner_metadata`, and `format_log_locations_section` were removed from `06-cicd-local-runner.py`.
3. **`TypeError: JobResult.__init__() takes 6 positional arguments but 7 were given`:**
   The flexible `*args, **kwargs` constructor for `JobResult` in `06-cicd-local-runner.py` was reverted to a rigid 5-field dataclass, breaking `.github/scripts/tests/test_ci_scripts.py`.

## 2. Why Not Caught Earlier?
The sync commit was pushed to `origin/main` without running local CI script tests (`python .github/scripts/tests/test_ci_scripts.py`). Because `chunk_items` is dynamically imported by linter scripts, static compilation did not catch the missing attribute.

## 3. Remediation
1. **Restored `chunk_items` and `WorkerHeartbeatMonitor`:**
   Re-added `chunk_items(items, chunk_size)` and `WorkerHeartbeatMonitor` to `03-ai-scripts/02-shared-engine.py`.
2. **Restored Path Utilities and Flexible `JobResult` Constructor:**
   Re-added `REPO_ROOT`, `CICD_DIR`, `normalize_repo_rel`, `format_banner_metadata`, `format_log_locations_section`, and the flexible `JobResult.__init__` to `03-ai-scripts/06-cicd-local-runner.py`.
3. **Verification:**
   - `python .github/scripts/tests/test_ci_scripts.py`: Ran 18 tests in 62.3s -> **OK**.
   - `python linter-scripts/check-enum-and-boolean.py`: Scanned 3087 files -> **PASS**.
   - `python linter-scripts/check-relative-paths.py`: Scanned 8387 files -> **PASS**.
   - `python linter-scripts/check-nested-ifs.py`: **PASS**.

## 4. Prevention
- Protect GitMap-specific extensions in `03-ai-scripts/` when running multi-repo sync scripts.
- Ensure `test_ci_scripts.py` is included in local pre-commit or pre-push verification.
