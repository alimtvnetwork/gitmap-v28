# CI Issue 84: Unused Variable `isAgyFixExplicitOnly` in `cmdpipeline` RCA

## 1. Reproduction & Symptoms
- **CI Run ID:** [36281677651](https://github.com/alimtvnetwork/gitmap-v28/actions/runs/36281677651)
- **Failing Workflows:**
  - `CI / Lint / golangci-lint (strict, fail on any error)`
  - `CI / Lint Baseline Guard / Unused diff (baseline-diff, full-path only)`
  - `CI / Full Suite Guard / golangci-lint (strict, full suite)`
  - `CI / Lint Baseline Diff / Diff vs baseline (fail only on NEW findings)`
- **Symptoms:**
  `golangci-lint` reported unused variable violation in `cli/cmdpipeline/pipeline_cmd.go`:
  ```text
  ##[error]cmdpipeline/pipeline_cmd.go:9:2: var `isAgyFixExplicitOnly` is unused (unused)
  	isAgyFixExplicitOnly bool
  	^
  ##[error]Process completed with exit code 1.
  ```

## 2. Root Cause Analysis
1. **Unused Package-Level Variable**: In `cli/cmdpipeline/pipeline_cmd.go`, `var isAgyFixExplicitOnly bool` was declared and assigned `isAgyFixExplicitOnly = true` inside `setupExplicitOnlyRunner()`, but was never referenced or checked anywhere in the codebase.
2. **Strict Linter Enforcement**: GitHub Actions CI enforces `unused` in `golangci-lint` and blocks new unused variables against the baseline via `.github/scripts/check-single-linter-diff.py`.

## 3. Code Fix
- Removed unused `isAgyFixExplicitOnly` declaration from `var (...)` block and removed assignment in `setupExplicitOnlyRunner()` in `cli/cmdpipeline/pipeline_cmd.go`.
- Verified clean `golangci-lint run --issues-exit-code=1 ./...` exit code 0 locally.

## 4. Prevention
- Always run `golangci-lint run --issues-exit-code=1 ./...` across modified Go packages before committing and tagging release versions.
