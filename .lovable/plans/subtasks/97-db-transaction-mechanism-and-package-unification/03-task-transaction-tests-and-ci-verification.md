# Subtask 03: Unit Testing & CI Verification

Parent Plan: `plans/pending/97-db-transaction-mechanism-and-package-unification.md`  
Status: In Progress  
Ownership: Execution Subagent 1 & 2

---

## 1. Objectives

1. Verify and update unit tests:
   - `gitmap/cmd/sequence_cmd_test.go`
   - `gitmap/cmd/sshjoin_cmd_test.go`
   - `gitmap/store/ssh_repo_test.go`
   - `gitmap/store/transaction_atomicity_test.go`

2. Add test coverage for:
   - Rollback on error in `gitmap/store/makeallvisibility.go`.
   - Rollback on error in `gitmap/store/pendingtask.go`.
   - Rollback on error in `gitmap/store/owner_repo_name_index.go`.
   - Atomic multi-profile export in `gitmap/cmd/chromeprofile_export_all.go`.

3. Run AST linters and CI local runner:
   - Run `python 03-ai-scripts/06-cicd-local-runner.py --force --no-tests`.
   - Ensure all 29 quality gates pass with exit code 0.

---

## 2. Acceptance Criteria

- All unit tests pass cleanly.
- `python 03-ai-scripts/06-cicd-local-runner.py --force --no-tests` passes with 0 failures.
- No new linters or formatting regressions introduced.
