# Issue: Legacy Refs Whitelist, Bare Stderr Prints in Cmd, and Lint Script KeyError

- date: 2026-08-29
- status: resolved
- scope: ci/linter/policy

## 1. Why it happened
Three parallel CI checks failed:
1. `Legacy Refs Check`: The policy script `check-legacy-refs.sh` scans for old version patterns (`gitmap-v[567]\b`). Previous RCA files documented resolving `gitmap-v6` in `spec/` files without an exemption tag. <!-- gitmap-legacy-ref-allow -->
2. `Bare Stderr Check`: Older subcommands in `gitmap/cmd/` printed errors using `fmt.Fprintln(os.Stderr, err)` rather than using the centralized `cliexit.Fail` standard.
3. `Lint Baseline Diff`: Python analysis scripts `lint-suggest.py` and `lint-diff.py` queried `res["is_failure"]`, but `query_wrapper` returns `is_fail`.

## 2. How it happened
- In `policy-check`, `check-legacy-refs.sh .` detected 4 lines containing `gitmap-v6` in `.lovable/cicd-issues/20-cicd-four-headed-hydra.md` and `.lovable/cicd-issues/index.md`. <!-- gitmap-legacy-ref-allow -->
- In `lint-baseline-guard`, `check-bare-stderr-err.sh` ran `rg` across `gitmap/cmd` and found 8 bare error prints across `dedupe.go`, `orphans.go`, `size.go`, `stale.go`, `visibilityhistory.go`, and `visibilityundo.go`.
- In `lint-baseline-diff`, `lint-suggest.py` threw `KeyError: 'is_failure'` during `load_findings()`.

## 3. Root Cause
- Missing whitelist annotation (`<!-- gitmap-legacy-ref-allow -->`) on legitimate historical references in CI log markdown files.
- Lack of uniform adoption of `cliexit.Fail` across older CLI subcommands in `gitmap/cmd`.
- Key name mismatch between `query_wrapper.py` (`is_fail`) and consumers (`res["is_failure"]`).

## 4. Code Fix
- Appended `<!-- gitmap-legacy-ref-allow -->` to historical lines in `.lovable/cicd-issues/`.
- Converted all 8 bare error prints in `gitmap/cmd/` to `cliexit.Fail`.
- Modified `lint-suggest.py` and `lint-diff.py` to evaluate `res.get("is_fail") or res.get("is_failure")`.
- Created `.lovable/ai-fix-scripts/03-cicd-local-runner.py` automating native local execution of all checks.
- Verified all checks pass with exit code 0.
