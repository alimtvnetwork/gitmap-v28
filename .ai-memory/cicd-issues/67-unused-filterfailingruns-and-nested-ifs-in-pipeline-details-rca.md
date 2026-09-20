# RCA: Unused `filterFailingRunsByTargetSha` and Nested-If Violations in Pipeline Details

- **Date:** 2026-09-21
- **Affected Run:** [CI #35521722367](https://github.com/alimtvnetwork/gitmap-v28/actions/runs/35521722367)
- **Target Commit:** `3b4d2fe255379093d31c434761b370ae6f2ad3ec`
- **Scope:** CI/CD Quality Gates (`golangci-lint`, `check-nested-ifs.py`, `check-enum-and-boolean.py`)

---

## 1. Symptom

Pipeline run [#35521722367](https://github.com/alimtvnetwork/gitmap-v28/actions/runs/35521722367) failed across 4 jobs:
1. **Lint** & **Full Suite Guard** (`golangci-lint`):
   ```text
   cmdpipeline/pipeline_logs.go:517:6: func `filterFailingRunsByTargetSha` is unused (unused)
   func filterFailingRunsByTargetSha(runs []ghRunItem) []ghRunItem {
        ^
   Process completed with exit code 1.
   ```
2. **Boolean & Enum Linter** & **Nested If Linter** (`check-nested-ifs.py`):
   ```text
   ❌ FAIL: Found 3 nested-if / anti-compression violation(s) across 1 file(s):
     cli/cmdpipeline/pipeline_details.go:168: Nested if statement found (depth 2): if offset < len(runs) {
     cli/cmdpipeline/pipeline_details.go:249: Nested if statement found (depth 2): if failStep := findFailingStepName(j.Steps); len(failStep) > 0 {
     cli/cmdpipeline/pipeline_details.go:254: Nested if statement found (depth 2): if activeStep := findActiveStepName(j.Steps); len(activeStep) > 0 {
   Process completed with exit code 1.
   ```

---

## 2. Root Cause

1. **Unused Function**: In `cli/cmdpipeline/pipeline_logs.go`, during the refactoring to scope pipeline errors strictly to target commit SHAs, `filterFailingRunsByTargetSha` became orphaned when callers switched to `findPrimaryTargetRun` and `collectRunsMatchingSha` directly.
2. **Nested If Violations**: In `cli/cmdpipeline/pipeline_details.go`:
   - `selectRunByIndex` enclosed `if offset < len(runs)` directly inside `if flags.HasIndex && flags.Index < 0`.
   - `resolveJobActiveOrFailingStep` enclosed `if failStep := findFailingStepName(j.Steps); len(failStep) > 0` and `if activeStep := findActiveStepName(j.Steps); len(activeStep) > 0` directly inside outer status checks.
   Both violated the repository guideline strictly banning nested `if` blocks (max depth 1).

---

## 3. Resolution

1. **Removed Dead Code**: Deleted `filterFailingRunsByTargetSha` from `cli/cmdpipeline/pipeline_logs.go`.
2. **Decomposed Nested Conditionals**:
   - Refactored `selectRunByIndex` to delegate negative offset handling to a dedicated helper `resolveNegativeOffsetRun`.
   - Refactored `resolveJobActiveOrFailingStep` to delegate to `resolveFailedJobStep` and `resolveActiveJobStep`.
3. **Corrected Struct Field Mapping**: Fixed `initBaseDetailsPayload` to properly set `Failures: make([]SectionFailure, 0)` and `DurationSeconds: calculateRunDuration(...)`.
4. **Verified Gates**: Ran `python 03-ai-scripts/06-cicd-local-runner.py --skip-tests` — all 38 gates passed 100% green in 39.52s.

---

## 4. Prevention & Learnings

- **Mandatory Pre-Push Local Gate**: Always execute `python 03-ai-scripts/06-cicd-local-runner.py --skip-tests` to catch both AST policy checks (`check-nested-ifs.py`) and Go linters (`golangci-lint`) before pushing.
- **Immediate Dead Code Pruning**: When refactoring helper methods, grep the codebase to ensure no orphaned declarations remain.
