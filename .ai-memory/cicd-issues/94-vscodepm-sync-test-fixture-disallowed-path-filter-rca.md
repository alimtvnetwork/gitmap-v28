# RCA-094: VS Code PM Sync Test Fixture Disallowed Path Filter Exclusion

**Date:** 2026-10-01
**Status:** Resolved
**Severity:** Critical (Full Suite Guard and Cross-Platform Build CI failure on Linux & macOS)
**Affected Packages:** `cli/cmdvscode`
**Run IDs:** `36750229027`, `36750228882`

---

## 1. Why It Happened (Architectural Reason)

In RCA-55 (commit `3e1162d059`), path guards (`IsDisallowedProjectPath` and `filterProductionPairs`) were integrated into `vscodepm.SyncMode` (`cli/vscodepm/sync.go`) to prevent transient test mock directories from being persisted to developers' production VS Code `projects.json` and Antigravity workspace databases. While `vscodepm.BypassDisallowedPathFilterForTesting` was introduced to allow unit tests using `t.TempDir()` to test sync logic, test fixture helper `swapHomeEnv` in `cli/cmdvscode/vscodepmsync_testhelper_test.go` omitted this bypass flag.

---

## 2. How It Happened (Execution Flow)

1. On Linux (`ubuntu-latest`) and macOS (`macos-latest`) in CI, `cmdvscode` test cases (`TestVSCodePMSyncDedupesAcrossSources`, `TestVSCodePMSyncModeUnionDefaultPreservesUserTags`, `TestVSCodePMSyncModeReplaceDropsUserTags`, `TestVSCodePMSyncModeIntersectionDropsExclusiveTags`, `TestVSCodePMSyncModeIntersectionPinsBrandWhenAbsent`, and `TestVSCodePMSyncEndToEndPreservesUserTags`) initialized mock project structures using `setupVSCodePMSyncFixtureWithTags(t, ...)`.
2. `setupVSCodePMSyncFixtureWithTags` called `t.TempDir()`, creating directories prefixed with `/tmp/` on Linux and `/var/folders/.../T/` on macOS, and initialized `swapHomeEnv(t, tmp)`.
3. Tests called `runVSCodePMSync(nil)` (or with `--mode replace/intersection`), which dispatched to `vscodepm.SyncMode(pairs, opts.Mode)`.
4. `vscodepm.SyncMode` invoked `filterProductionPairs(pairs)`, where `IsDisallowedProjectPath(p.RootPath)` evaluated to `true` because the paths started with `/tmp` or `os.TempDir()`.
5. All pairs were dropped (`len(filtered) == 0`), causing `SyncMode` to early-return `SyncSummary{}, nil` without modifying `projects.json`.
6. Test assertions checking tag modifications and brand-tag additions read the unmutated seed files, failing with missing/duplicated tags (`want 1, got 2`, `brand tag missing`).

---

## 3. Root Cause

`cli/cmdvscode/vscodepmsync_testhelper_test.go` (`swapHomeEnv`):
The environment sandbox helper isolated `HOME`, `XDG_CONFIG_HOME`, and `APPDATA`, but did not set `vscodepm.BypassDisallowedPathFilterForTesting = true`.

---

## 4. Code Fix

In `cli/cmdvscode/vscodepmsync_testhelper_test.go`:

```go
// swapHomeEnv points HOME / XDG_CONFIG_HOME / APPDATA at tmp using t.Setenv.
func swapHomeEnv(t *testing.T, tmp string) func() {
	t.Helper()
	t.Setenv(constants.VSCodeEnvHome, tmp)
	t.Setenv(constants.VSCodeEnvXDGConfigHome, filepath.Join(tmp, ".config"))
	t.Setenv(constants.VSCodeEnvAppData, tmp)

	vscodepm.BypassDisallowedPathFilterForTesting = true
	t.Cleanup(func() {
		vscodepm.BypassDisallowedPathFilterForTesting = false
	})

	return func() {
		vscodepm.BypassDisallowedPathFilterForTesting = false
	}
}
```

---

## 5. Verification

1. Verified targeted package tests: `go test -v ./cmdvscode` runs clean with all fixture tests passing.
2. Verified newline styling and boolean guidelines linters (`check-newline-styling.py`, `05-guideline-autofixer.py`).
3. Recorded modified files under lock via `03-ai-scripts/33-test-inventory-generator.py --record`.
