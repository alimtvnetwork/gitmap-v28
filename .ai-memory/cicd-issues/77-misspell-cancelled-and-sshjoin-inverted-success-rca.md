# CI/CD RCA: Misspell `cancelled` in Pipeline Tests & Inverted Success Check in `sshjoin_common.go`

- Job: Lint / Full Suite Guard (`golangci-lint`) & Boolean Guidelines Linter (`check-boolean-guidelines.py`)
- Type: FAIL
- Run ID: 35889824634 (Commit: `cfe28b1`)
- Status: ✅ Resolved

---

## 1. Symptom

On remote GitHub Actions CI run `#35889824634`, three jobs failed:

1. **Lint (`golangci-lint` strict, fail on any error):**
   ```text
   cmdpipeline/pipeline_commit_groups_test.go:63:40: `cancelled` is a misspelling of `canceled` (misspell)
   cmdpipeline/pipeline_commit_groups_test.go:65:40: `cancelled` is a misspelling of `canceled` (misspell)
   cmdpipeline/pipeline_commit_groups_test.go:66:45: `cancelled` is a misspelling of `canceled` (misspell)
   cmdpipeline/pipeline_history_test.go:430:22: `cancelled` is a misspelling of `canceled` (misspell)
   cmdpipeline/pipeline_history_test.go:452:29: `cancelled` is a misspelling of `canceled` (misspell)
   cmdpipeline/pipeline_history_test.go:471:31: `Cancelled` is a misspelling of `Canceled` (misspell)
   ```

2. **Boolean Guidelines Linter (`linter-scripts/check-boolean-guidelines.py`):**
   ```text
   ❌ FAIL: Found 1 boolean guideline violation(s) across 1 file(s):
     cli/cmdssh/sshjoin_common.go:96: Inverted success check (!isSuccess): if !t.IsSuccess {
   ```

3. **Full Suite Guard (`golangci-lint` strict, full suite):**
   Failed on the same misspell violations in `pipeline_commit_groups_test.go` and `pipeline_history_test.go`.

---

## 2. Root Cause

1. **British English Spelling in Tests (`misspell`):**
   `golangci-lint` enforces American English spelling via the `misspell` linter. In `pipeline_commit_groups_test.go` and `pipeline_history_test.go`, the British spelling `cancelled` (with two L's) was written in test case strings and function names (`TestIsCommitGroupFailure_Cancelled`, `TestGroupRunsByCommit_CancelledWorkflows`). Because the underlying pipeline functions (`isFailingConclusion` and `formatWorkflowShortStatus`) check `strings.HasPrefix(lower, "cancel")`, either spelling is supported by runtime logic, but tests must use standard US English `canceled` to satisfy static analysis without lint bypasses.

2. **Inverted Success Check (`check-boolean-guidelines.py`):**
   In `cli/cmdssh/sshjoin_common.go:96`, the batch join result rendering logic had:
   ```go
   if !t.IsSuccess {
       statusStr = constants.ColorRed + "FAILED" + constants.ColorReset
       details = t.ErrorMsg
   }
   ```
   Repository coding guidelines strictly ban inverted checks on success (`!isSuccess` / `!t.IsSuccess`), requiring positive framing. In addition, `renderCommonJoinTable` exceeded the 15-line canonical function size limit.

---

## 3. Resolution

1. **American English Standardization:**
   - In `cli/cmdpipeline/pipeline_commit_groups_test.go`, renamed `TestGroupRunsByCommit_CancelledWorkflows` to `TestGroupRunsByCommit_CanceledWorkflows`, updated `makeTestRun` arguments from `"cancelled"` to `"canceled"`, and renamed helpers to `verifyCanceledFirstGroup` and `verifyCanceledSecondGroup`.
   - In `cli/cmdpipeline/pipeline_history_test.go`, updated test cases and function names from `Cancelled` to `Canceled` and test inputs from `"cancelled"` to `"canceled"`.

2. **Affirmative Boolean & Modular Extraction:**
   - In `cli/cmdssh/sshjoin_common.go`, extracted `renderCommonJoinTargetRow(w io.Writer, t SSHCommonJoinTarget)` which tests `if t.IsSuccess` positively and sets `SUCCESS` / `enrolled`, defaulting to `FAILED` / `t.ErrorMsg`.
   - Reduced `renderCommonJoinTable` to 9 lines, satisfying both the boolean linter and function length requirements.

---

## 4. Prevention & Learnings

- **US English Everywhere in Go:** Always use US English `canceled` (one L) across Go files, comments, and tests.
- **Positive Boolean Framing:** Never write `if !isSuccess` or `if !t.IsSuccess`. Default to the failure branch or test `if t.IsSuccess` explicitly.
- **Local Runner Pre-Check:** Run `python 03-ai-scripts/06-cicd-local-runner.py --filter "golangci-lint (strict)"` and `python linter-scripts/check-boolean-guidelines.py` prior to cutting a release.
