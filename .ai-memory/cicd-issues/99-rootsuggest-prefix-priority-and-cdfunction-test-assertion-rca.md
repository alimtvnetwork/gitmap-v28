# RCA-099: Root Suggest Prefix Match Priority and PowerShell CD Shim Test Assertion

**Date:** 2026-10-04
**Status:** ✅ Resolved
**Severity:** Moderate (Unit Test Failure in CI macOS/Ubuntu test suites)
**Affected Workflows:** `CI (#37227985564)`
**Run URLs:**
- `https://github.com/alimtvnetwork/gitmap-v28/actions/runs/37227985564`

---

## 1. Symptom

During the unit test execution on macOS and Ubuntu runners for commit `6885eeff`, two test cases failed:

1. **`TestSuggestTopLevelCommands_Typos` in `cli/cmd/rootsuggest_test.go:52`:**
   ```text
   rootsuggest_test.go:52: suggestTopLevelCommands("ss") = [os sc as], expected "ssh"
   ```

2. **`TestRenderPowerShellCommandShimPinsInstalledExe` in `cli/completion/cdfunction_test.go:223`:**
   ```text
   cdfunction_test.go:223: shim missing "Set-Location -LiteralPath ([string]$dest)"
   ```

---

## 2. Root Cause

1. **Command Suggestion Ranking Disregarded Prefix Matches:**
   - In `cli/cmd/rootsuggest_calc.go:rankCandidateCommands()`, candidate commands with identical edit distance (such as `d = 1` for `"ss"` to `"os"`, `"sc"`, `"as"`, and `"ssh"`) were sorted strictly by length difference `candidateLenDiff`.
   - Because 2-letter typo commands (`"os"`, `"sc"`, `"as"`) had `lenDiff = 0` while `"ssh"` had `lenDiff = 1`, the 3-element cutoff in `selectBestSuggestions` truncated `"ssh"`, despite `"ssh"` being an exact prefix match for `"ss"`.

2. **PowerShell Command Shim Test Expected Outdated Variable Name:**
   - In `cli/completion/cdfunction_test.go:216`, the test asserted `"Set-Location -LiteralPath ([string]$dest)"`.
   - However, canonical PowerShell shim templates in `cli/constants/constants_cd_shim.go` and `cli/constants/constants_cd.go` consistently generate `$target` (`Set-Location -LiteralPath ([string]$target)`), causing an assertion mismatch against the actual generated script.

---

## 3. Resolution

1. **Prioritized Prefix Matches in Suggestion Sorting:**
   - In `cli/cmd/rootsuggest_calc.go:rankCandidateCommands()`, updated `sort.SliceStable` to compare `strings.HasPrefix(cmd, input)`. If one candidate is a prefix match and the other is not, the prefix match is always sorted ahead.
2. **Corrected Test Assertion in `cdfunction_test.go`:**
   - Updated `wants` slice in `TestRenderPowerShellCommandShimPinsInstalledExe` to check for `"Set-Location -LiteralPath ([string]$target)"`.

---

## 4. Verification

- `go test -v ./cmd -run TestSuggestTopLevelCommands`: All tests passed cleanly.
- `go test -v ./completion`: All tests passed cleanly, including `TestRenderPowerShellCommandShimPinsInstalledExe`.
- `python linter-scripts/check-relative-paths.py`: 8,656 files checked, 0 violations.
- `python 03-ai-scripts/49-verify-privacy-and-relative-paths.py`: All 31 files passed with 0 violations.
