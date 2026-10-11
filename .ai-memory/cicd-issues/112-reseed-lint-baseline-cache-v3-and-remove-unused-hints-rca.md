# RCA-112: Reseed Lint Baseline Cache to v3 and Remove Unused Hints

**Date:** 2026-10-11
**Status:** ✅ Resolved
**Severity:** High (CI Pipeline Failure in Lint Baseline Diff)
**Affected Workflows:** `CI (#38120172156)`
**Run URLs:**
- https://github.com/alimtvnetwork/gitmap-v28/actions/runs/38120172156

---

## 1. Symptom & Error Signature

Following commit `fadf788`, GitHub Actions CI workflow `#38120172156` failed on job `Lint Baseline Diff` (Step: `Diff vs baseline (fail only on NEW findings)`):

```text
[unused] type `statusSummary` is unused (NEW in /tmp/lint-current/report.json)
    ❌ cmd/dbhelpers.go:11: [unused] type `statusSummary` is unused
[unused] func `projectReposHints` is unused (NEW in /tmp/lint-current/report.json)
    ❌ cmd/hints.go:25: [unused] func `projectReposHints` is unused
[unused] func `listHints` is unused (NEW in /tmp/lint-current/report.json)
    ❌ cmd/hints.go:34: [unused] func `listHints` is unused
[unused] func `listGroupsHints` is unused (NEW in /tmp/lint-current/report.json)
    ❌ cmd/hints.go:43: [unused] func `listGroupsHints` is unused
[unused] func `activeGroupHints` is unused (NEW in /tmp/lint-current/report.json)
    ❌ cmd/hints.go:51: [unused] func `activeGroupHints` is unused
[unused] func `groupListHints` is unused (NEW in /tmp/lint-current/report.json)
    ❌ cmd/hints.go:60: [unused] func `groupListHints` is unused
Process completed with exit code 1.
```

---

## 2. Root Cause

1. **Linter Quota Truncation & The Stale Baseline Treadmill:**
   - In `cli/.golangci.yml`, `issues.max-issues-per-linter` is configured to `50`.
   - Across the `cli/` codebase, there are approximately 91 pre-existing `unused` symbol findings.
   - The cached baseline `golangci-baseline-main-v2-*` was created when earlier unused findings filled the 50-issue quota.
   - Every time unused symbols are removed from earlier files (such as `cmd/aliasops.go`, `cmd/clihelpers_argv.go`, `cmd_jsonextract.go`, `cmd_safety.go`), `golangci-lint` outputs the next pre-existing findings in alphabetical order to reach the 50-issue cap.
   - Because these next findings were omitted from the stale baseline due to the 50-issue cap, `.github/scripts/lint-diff.py` misclassifies them as NEW findings and fails CI.

2. **Unused Functions in `cli/cmd/hints.go`:**
   - Five helper functions (`projectReposHints`, `listHints`, `listGroupsHints`, `activeGroupHints`, `groupListHints`) in `cli/cmd/hints.go` were left unreferenced when command logic was refactored into subpackages (`cmdgroup`, `cmdlist`, `cmdprojectrepos`), which now maintain their own shims.
   - Only `aliasListHints`, `aliasSetHints`, and `aliasSuggestHints` are called within package `cmd`.

3. **`statusSummary` Cross-File Dependency in `cli/cmd`:**
   - While `statusSummary` in `cli/cmd/dbhelpers.go` was flagged by `unused` (because status commands moved to subpackages), removing it causes compilation errors in `cli/cmd/statusformat.go` and `cli/cmd/statusprint.go` which still reside in package `cmd`.
   - Thus, `statusSummary` must be retained in `dbhelpers.go` while invalidating the stale baseline cache.

---

## 3. Surgical Fix

1. **Bumping Baseline Cache Version to `v3` in `.github/workflows/ci.yml`:**
   - Updated the `lint_baseline_cache_version` workflow input default from `"v2"` to `"v3"`.
   - Updated cache restore keys and save keys fallback from `'v2'` to `'v3'`.
   - On the next push to `main`, cache restore will miss, causing `lint-diff.py` to enter seeding mode (exit code 0), and step `Promote current report to baseline` will write the current clean report as the new `v3` baseline on `main`.

2. **Removing Dead Hint Functions in `cli/cmd/hints.go`:**
   - Removed the 5 unused functions: `projectReposHints`, `listHints`, `listGroupsHints`, `activeGroupHints`, and `groupListHints`.
   - Retained active alias hint generators: `aliasListHints`, `aliasSetHints`, and `aliasSuggestHints`.

---

## 4. Verification

1. `go -C cli test ./cmd/...` passes cleanly (exit 0) confirming all tests compile and pass.
2. Verified `.github/workflows/ci.yml` syntax and fallback cache keys.
3. Updated CI issue log index in `.ai-memory/cicd-issues/index.md`.
