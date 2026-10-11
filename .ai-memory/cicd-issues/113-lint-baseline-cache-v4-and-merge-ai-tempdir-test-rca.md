# RCA-113: Reseed Lint Baseline Cache to v4 and Hermetic Merge-AI TempDir Test Fix

**Date:** 2026-10-11
**Status:** ✅ Resolved
**Severity:** High (CI Pipeline Failures in Lint Baseline Guard, Full Suite Guard, and Cross-Platform Build)
**Affected Workflows:**
- `CI (#38136605865, #38137083987)`
- `Cross-Platform Build (#38136605746, #38137083855)`
**Run URLs:**
- https://github.com/alimtvnetwork/gitmap-v28/actions/runs/38136605865
- https://github.com/alimtvnetwork/gitmap-v28/actions/runs/38136605746

---

## 1. Symptom & Error Signature

Following commits `aa487f0` and subsequent docs pushes on `main`:
1. In `CI (#38136605865)`:
   - **Job:** `Lint Baseline Guard (unused, gosec G115, misspell, gocritic, exhaustive)`
   - **Step:** `Unused diff (baseline-diff, full-path only)` (step #7)
   - **Error Signature:**
     ```text
     FAIL: 86 new unused finding(s) detected!
     Script:   .github/scripts/check-single-linter-diff.py
     Target:   cli
     To reproduce locally:
       golangci-lint run --no-config --disable-all --enable=unused ./cli/...
     ```
   - **Job:** `Full Suite Guard`
   - **Step:** `Run full test suite`
   - **Error Signature:**
     ```text
     --- FAIL: TestParseMergeAIArgs_MultiWordDestination (0.13s)
         source_parser_test.go:228: unexpected error: source file not found: merge-ai-icon-studio.txt
     FAIL	github.com/alimtvnetwork/gitmap-v28/cli/cmdmergeai	0.765s
     ```

2. In `Cross-Platform Build (#38136605746, #38137083855)`:
   - **Jobs:** `ubuntu-latest / go build + test`, `macos-latest / go build + test`
   - **Step:** `go test ./... (cached)`
   - **Error Signature:**
     ```text
     --- FAIL: TestParseMergeAIArgs_MultiWordDestination (0.07s)
         source_parser_test.go:228: unexpected error: source file not found: merge-ai-icon-studio.txt
     FAIL	github.com/alimtvnetwork/gitmap-v28/cli/cmdmergeai	0.232s
     ```

---

## 2. Root Cause

1. **Ambient Non-Hermetic Test File Dependency in `TestParseMergeAIArgs_MultiWordDestination`:**
   - In `cli/cmdmergeai/source_parser_test.go`, `TestParseMergeAIArgs_MultiWordDestination` invoked `parseMergeAIArgs([]string{"Icon", "Studio", "v1", "merge-ai-icon-studio.txt"})`.
   - `merge-ai-icon-studio.txt` was located locally in `d:\work\repo-cache\merge-ai-icon-studio.txt` on the developer workstation, so the test passed locally via `probeWorkBaseCandidate()`.
   - In remote GitHub Actions CI environments (Ubuntu, macOS), no such ambient file existed in workspace or repo-cache, causing `ResolveSourceFilePath` to return `source file not found: merge-ai-icon-studio.txt`.

2. **Unversioned Stale Baseline Cache in `lint-baseline-guard`:**
   - In `.github/workflows/ci.yml`, the `lint-baseline-guard` job restored `unused`, `gosec-g115`, `misspell`, `gocritic`, and `exhaustive` baselines using unversioned cache keys (e.g. `golangci-unused-baseline-main-`).
   - While `lint-baseline-diff` was previously versioned with `v3`, `lint-baseline-guard` retained unversioned restore keys that matched an old cached report containing only 1 issue.
   - When compared against the full 87 pre-existing findings across `cli`, `check-single-linter-diff.py` reported 86 "NEW" findings and failed the step, preventing the cache from ever being refreshed on `main`.

---

## 3. Surgical Fix

1. **Hermetic TempDir in `cli/cmdmergeai/source_parser_test.go`:**
   - Updated `TestParseMergeAIArgs_MultiWordDestination` to write a temporary manifest file inside `t.TempDir()` with test source URLs.
   - Tested both positional arguments and multi-word destinations against the isolated temporary file path.

2. **Bumping Baseline Cache Version to `v4` and Versioning All Baseline Guards:**
   - Updated `lint_baseline_cache_version` default input in `.github/workflows/ci.yml` from `"v3"` to `"v4"`.
   - Updated all 5 single-linter baseline cache keys and restore keys in `lint-baseline-guard` (`unused`, `gosec-g115`, `misspell`, `gocritic`, `exhaustive`) to use `${{ inputs.lint_baseline_cache_version || 'v4' }}`.
   - Updated `lint-baseline-diff` restore and save keys to use `'v4'`.
   - This invalidates the stale unversioned cache, allowing `check-single-linter-diff.py` to enter seeding mode on the first run on `main`, saving the complete baseline.

---

## 4. Verification

1. `go -C cli test -v -run TestParseMergeAIArgs_MultiWordDestination ./cmdmergeai/...` passes cleanly.
2. `go -C cli test ./cmdmergeai/...` passes cleanly (exit 0).
3. `python .github/scripts/tests/test_ci_scripts.py` passes all 18 test cases cleanly.
4. `.github/workflows/ci.yml` lint baseline cache keys verified across all jobs.
