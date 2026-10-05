# Subtask 04: Parser Unit Tests, Verification & Release Ceremony

> **Parent Plan:** `220-pipeline-pe-unit-test-traceback-and-heatmap`  
> **Status:** DONE  
> **Target:** `cli/cmdpipeline/pipeline_error_extract_test.go`, `version.json`, `package.json`, `changelog.md`  

---

## Objectives
1. Add regression tests in `cli/cmdpipeline/pipeline_error_extract_test.go`:
   - Verifying Python traceback extraction from GitHub Actions raw log.
   - Verifying test failure summary elevation (`FAIL: test_...`).
   - Verifying noise line filtering (`Ran 18 tests in 75.725s`).
2. Run targeted tests to verify 100% pass rate.
3. Perform minor version bump (`v6.475.0` -> `v6.476.0`) across `version.json`, `package.json`, `changelog.md`.
4. Create `.ai-memory/release/release-notes-v6.476.0.md`.
5. Atomic commit and push via GitMap hyphen convention.
